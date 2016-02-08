package main

import (
	"regexp"
	"strings"
	"sync"

	"github.com/trustmaster/go-aspell"
)

type Spellcheck struct {
	speller aspell.Speller
	toks    map[string]struct{}
	check   CheckFunc
	mu      sync.Mutex
	strips  []*regexp.Regexp
}

func NewSpellcheck(ignoreFile string) (sc *Spellcheck, err error) {
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

	igns, err := WithPassIgnores(ignoreFile)
	if err != nil {
		return nil, err
	}

	speller, err := aspell.NewSpeller(opts)
	if err != nil {
		return nil, err
	}

	strips := []*regexp.Regexp{
		regexp.MustCompile("http(s|):[^ ]*"),
		regexp.MustCompile("TODO[ ]*\\([a-z]*"),
	}

	sc = &Spellcheck{
		strips:  strips,
		speller: speller,
	}

	checks := []CheckFunc{
		sc.WithPassTokens(),
		igns,
		WithPassNumbers(),
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
		_, ok := sc.toks[word]
		return ok
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

func (sc *Spellcheck) Check(srcpaths []string) ([]*CheckedLexeme, error) {
	toks, err := GoTokens(srcpaths)
	if err != nil {
		return nil, err
	}

	sc.toks = make(map[string]struct{})
	for k, _ := range toks {
		for _, field := range strings.Fields(k) {
			sc.toks[field] = struct{}{}
		}
	}

	errc := make(chan error)
	badcommc := make(chan *CheckedLexeme)
	badcomms := &[]*CheckedLexeme{}
	go func() {
		for comm := range badcommc {
			*badcomms = append(*badcomms, comm)
		}
		errc <- nil
	}()

	// process all comments
	for _, p := range srcpaths {
		lc, err := LexemeChan(p)
		if err != nil {
			go func() {
				errc <- err
				errc <- nil
			}()
			continue
		}
		mux := LexemeMux(lc, 2)
		go func() {
			sc.checkComments(mux[0], badcommc)
			errc <- nil
		}()
		go func() {
			checkGoDocs(mux[1], badcommc)
			errc <- nil
		}()
	}

	// wait for completion of readers
	for i := 0; i < len(srcpaths)*2; i++ {
		if curErr := <-errc; curErr != nil {
			err = curErr
		}
	}

	// wait to collect all bad comments
	close(badcommc)
	<-errc

	return *badcomms, err
}

func (sc *Spellcheck) checkComments(lc <-chan *Lexeme, outc chan<- *CheckedLexeme) {
	ch := Filter(lc, CommentFilter)
	for comm := range ch {
		if ct := sc.checkLexeme(comm); ct != nil {
			outc <- ct
		}
	}
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
	for _, re := range sc.strips {
		s = string(re.ReplaceAll([]byte(s), []byte("")))
	}
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

func (sc *Spellcheck) checkLexeme(ct *Lexeme) (ret *CheckedLexeme) {
	for _, word := range sc.tokenize(ct.lit) {
		if sc.check(word) {
			continue
		}
		if ret == nil {
			ret = &CheckedLexeme{ct, "spell", nil}
		}
		ret.words = append(ret.words, CheckedWord{word, sc.suggest(word)})
	}
	return ret
}
