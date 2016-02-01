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
	check   CheckFunc
	mu      sync.Mutex
	reURL   *regexp.Regexp
	reTODO  *regexp.Regexp
}

type CheckFunc func(string) bool

type CheckedWord struct {
	word    string
	suggest string
}

type CheckedToken struct {
	ctok  *CommentToken
	words []CheckedWord
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

	speller, err := aspell.NewSpeller(opts)
	if err != nil {
		return nil, err
	}

	sc = &Spellcheck{
		reURL:   regexp.MustCompile("http(s|):[^ ]*"),
		reTODO:  regexp.MustCompile("TODO[ ]*\\([a-z]*"),
		ignores: ignmap,
		speller: speller,
	}

	checks := []CheckFunc{
		sc.WithPassTokens(),
		sc.WithPassIgnores(),
		sc.WithPassNumbers(),
		sc.WithASpell(),
	}
	sc.check = func(w string) bool {
		for _, chk := range checks {
			if chk(w) {
				return true
			}
		}
		return false
	}
	return sc, nil
}

func (sc *Spellcheck) Close() {
	defer sc.speller.Delete()
}

func (sc *Spellcheck) WithPassTokens() CheckFunc {
	return func(word string) bool {
		_, ok := sc.toks[strings.ToLower(word)]
		return ok
	}
}

func (sc *Spellcheck) WithPassIgnores() CheckFunc {
	return func(word string) bool {
		_, ok := sc.ignores[word]
		return ok
	}
}

func (sc *Spellcheck) WithPassNumbers() CheckFunc {
	return func(word string) bool {
		// contains a number?
		for i := 0; i < len(word); i++ {
			if word[i] >= '0' && word[i] <= '9' {
				return true
			}
		}
		return false
	}
}

func (sc *Spellcheck) WithASpell() CheckFunc {
	return func(word string) bool {
		// aspell's check isn't thread-safe-- why!?
		sc.mu.Lock()
		defer sc.mu.Unlock()
		return sc.speller.Check(word)
	}
}

func (sc *Spellcheck) Check(srcpaths []string) ([]*CheckedToken, error) {
	toks, err := GoTokens(srcpaths)
	if err != nil {
		return nil, err
	}

	sc.toks = make(map[string]struct{})
	for k, _ := range toks {
		sc.toks[strings.ToLower(k)] = struct{}{}
	}

	errc := make(chan error)
	badcommc := make(chan *CheckedToken)
	badcomms := &[]*CheckedToken{}
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
				if ct := sc.checkComment(comm); ct != nil {
					badcommc <- ct
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

func (sc *Spellcheck) suggest(word string) string {
	if sc.check(word) {
		return ""
	}
	// aspell's check isn't thread-safe-- why!?
	sc.mu.Lock()
	suggest := sc.speller.Suggest(word)
	sc.mu.Unlock()
	if len(suggest) == 0 {
		return ""
	}
	return suggest[0]
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

func (sc *Spellcheck) checkComment(ct *CommentToken) (ret *CheckedToken) {
	for _, word := range sc.tokenize(ct.lit) {
		if sc.check(word) {
			continue
		}
		if ret == nil {
			ret = &CheckedToken{ct, nil}
		}
		ret.words = append(ret.words, CheckedWord{word, sc.suggest(word)})
	}
	return ret
}
