package main

import (
	"io/ioutil"
	"strings"
)

type CheckPipe func(lc <-chan *Lexeme, outc chan<- *CheckedLexeme)

type CheckFunc func(string) bool

type CheckedWord struct {
	word    string
	suggest string
}

type CheckedLexeme struct {
	ctok  *Lexeme
	rule  string
	words []CheckedWord
}

type CheckedLexemes []*CheckedLexeme

// sort.Interface implementation:

func (s CheckedLexemes) Len() int { return len(s) }
func (s CheckedLexemes) Less(i, j int) bool {
	return s[i].ctok.pos.Filename < s[j].ctok.pos.Filename ||
		(s[i].ctok.pos.Filename == s[j].ctok.pos.Filename &&
			s[i].ctok.pos.Line < s[j].ctok.pos.Line) ||
		(s[i].ctok.pos.Filename == s[j].ctok.pos.Filename &&
			s[i].ctok.pos.Line == s[j].ctok.pos.Line &&
			s[i].ctok.pos.Column < s[j].ctok.pos.Column)
}
func (s CheckedLexemes) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func WithPassIgnores(ignoreFile string) (CheckFunc, error) {
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
	return func(word string) bool {
		_, ok := ignmap[word]
		return ok
	}, nil
}

func WithPassNumbers() CheckFunc {
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

func Check(srcs []string, cps []CheckPipe) ([]*CheckedLexeme, error) {
	var err error
	errc := make(chan error)
	badcommc := make(chan *CheckedLexeme)
	badcomms := []*CheckedLexeme{}
	go func() {
		for comm := range badcommc {
			badcomms = append(badcomms, comm)
		}
		errc <- nil
	}()

	// process all files under all checkers
	for _, p := range srcs {
		lc, lerr := LexemeChan(p)
		if lerr != nil {
			go func() {
				errc <- lerr
				errc <- nil
			}()
			continue
		}
		mux := LexemeMux(lc, len(cps))
		for i := range cps {
			go func(k int) {
				cps[k](mux[k], badcommc)
				errc <- nil
			}(i)
		}
	}

	// wait for completion of readers
	for i := 0; i < len(srcs)*len(cps); i++ {
		if curErr := <-errc; curErr != nil {
			err = curErr
		}
	}

	// wait to collect all bad comments
	close(badcommc)
	<-errc

	return badcomms, err
}

func CheckAll(paths []string) ([]*CheckedLexeme, error) {
	sp, err := newSpellCheck(paths, "")
	if err != nil {
		return nil, err
	}
	defer sp.Close()
	cps := []CheckPipe{CheckGoDocs, sp.Check()}
	cts, cerr := Check(paths, cps)
	if cerr != nil {
		return nil, cerr
	}
	return cts, nil
}
