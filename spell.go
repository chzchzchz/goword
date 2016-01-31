package main

import (
	"io/ioutil"
	"regexp"
	"strings"
	"sync"

	"github.com/trustmaster/go-aspell"
)

type Spellcheck struct {
	speller aspell.Speller
	toks    map[string]struct{}
	ignores map[string]struct{}
	mu      sync.Mutex
	reURL   *regexp.Regexp
	reTODO  *regexp.Regexp
}

func NewSpellcheck(ignoreFile string) (sc *Spellcheck, err error) {
	ignmap := make(map[string]struct{})
	if ignoreFile != "" {
		igns, rerr := ioutil.ReadFile(ignoreFile)
		if rerr != nil {
			return nil, rerr
		}
		for _, word := range strings.Fields(string(igns)) {
			ignmap[word] = struct{}{}
		}
	}

	opts := map[string]string{
		"lang":           "en",
		"filter":         "url",
		"mode":           "url",
		"encoding":       "ascii",
		"guess":          "true",
		"ignore":         "0",
		"ignore-case":    "false",
		"ignore-accents": "false",
	}
	sc = &Spellcheck{
		reURL:   regexp.MustCompile("http(s|):[^ ]*"),
		reTODO:  regexp.MustCompile("TODO[ ]*\\([a-z]*"),
		ignores: ignmap,
	}
	sc.speller, err = aspell.NewSpeller(opts)
	if err != nil {
		return nil, err
	}
	return sc, nil
}

func (sc *Spellcheck) Close() {
	defer sc.speller.Delete()
}

func (sc *Spellcheck) Check(srcpaths []string) ([]*CommentToken, error) {
	toks, err := GoTokens(srcpaths)
	if err != nil {
		return nil, err
	}

	sc.toks = make(map[string]struct{})
	for k, _ := range toks {
		sc.toks[strings.ToLower(k)] = struct{}{}
	}

	errc := make(chan error)
	badcommc := make(chan *CommentToken)
	badcomms := &[]*CommentToken{}
	go func() {
		for comm := range badcommc {
			*badcomms = append(*badcomms, comm)
		}
		errc <- nil
	}()

	// process all comments
	for _, p := range srcpaths {
		go func(path string) {
			commc, cerr := GoCommentChan(path)
			if cerr != nil {
				errc <- cerr
				return
			}
			for comm := range commc {
				if sc.badComment(comm) {
					badcommc <- comm
				}
			}
			errc <- nil
		}(p)
	}

	// wait for completion of readers
	for range srcpaths {
		if curErr := <-errc; curErr != nil {
			err = curErr
		}
	}

	// wait to collect all bad comments
	close(badcommc)
	<-errc

	return *badcomms, err
}

func (sc *Spellcheck) tokenize(s string) []string {
	s = string(sc.reURL.ReplaceAll([]byte(s), []byte("")))
	s = string(sc.reTODO.ReplaceAll([]byte(s), []byte("")))
	x := []string{
		".", "`", "\"", ",", "!", "?",
		";", ")", "(", "/", ":", "=",
		"*", "-", ">", "]", "[", "_",
		"|", "{", "}", "+", "\t", "' ",
		" '", "&", "<", "'s "}
	for _, v := range x {
		s = strings.Replace(s, v, " ", -1)
	}
	return strings.Fields(s)
}

func (sc *Spellcheck) isGoodWord(word string) bool {
	if _, ok := sc.ignores[word]; ok {
		return true
	}
	lower := strings.ToLower(word)
	if _, ok := sc.toks[lower]; ok {
		return true
	}
	for i := 0; i < len(word); i++ {
		if word[i] >= '0' && word[i] <= '9' {
			return true
		}
	}
	// aspell's check isn't thread-safe-- why!?
	sc.mu.Lock()
	aspell_ok := sc.speller.Check(word)
	sc.mu.Unlock()
	return aspell_ok
}

func (sc *Spellcheck) badComment(ct *CommentToken) bool {
	for _, word := range sc.tokenize(ct.lit) {
		if sc.isGoodWord(word) {
			continue
		}
		return true
	}
	return false
}

func (sc *Spellcheck) Suggest(s string) string {
	for _, word := range sc.tokenize(s) {
		if sc.isGoodWord(word) {
			continue
		}
		// aspell's check isn't thread-safe-- why!?
		sc.mu.Lock()
		aspell_ok := sc.speller.Check(word)
		sc.mu.Unlock()
		if !aspell_ok {
			sc.mu.Lock()
			suggest := sc.speller.Suggest(word)
			sc.mu.Unlock()
			if len(suggest) > 0 {
				return word + " -> " + suggest[0]
			}
		}
	}
	return ""
}
