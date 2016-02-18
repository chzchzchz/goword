package main

import (
	"fmt"
	"regexp"
	"strings"
)

type Spellcheck struct {
	speller Speller
	toks    map[string]struct{}
	check   CheckFunc
	strips  []*regexp.Regexp
}

func NewSpellcheck(ignoreFile string) (sc *Spellcheck, err error) {
	igns, err := WithPassIgnores(ignoreFile)
	if err != nil {
		return nil, err
	}

	strips := []*regexp.Regexp{
		regexp.MustCompile("http(s|):[^ ]*"),
		regexp.MustCompile("TODO[ ]*\\([a-z]*"),
	}

	hsp := NewHunSpeller()
	if hsp == nil {
		return nil, fmt.Errorf("bad hunspell")
	}
	asp, err := NewASpeller()
	if err != nil {
		return nil, err
	}

	sc = &Spellcheck{
		strips:  strips,
		speller: NewMultiSpeller(hsp, asp),
	}

	checks := []CheckFunc{
		sc.WithPassTokens(),
		igns,
		WithPassNumbers(),
		sc.WithSpeller(),
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

func (sc *Spellcheck) Close() { sc.speller.Close() }

func (sc *Spellcheck) WithPassTokens() CheckFunc {
	return func(word string) bool {
		_, ok := sc.toks[word]
		return ok
	}
}

func (sc *Spellcheck) WithSpeller() CheckFunc {
	return func(word string) bool { return sc.speller.Check(word) }
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
		lc, lerr := LexemeChan(p)
		if lerr != nil {
			go func() {
				errc <- lerr
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
	if sugg := sc.speller.Suggest(word); len(sugg) > 0 {
		return sugg[0]
	}
	return ""
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
