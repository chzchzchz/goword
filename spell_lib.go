// +build spell

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

// isCodeRE matches assignment, function calls, imports, and braces.
var isCodeRE *regexp.Regexp = regexp.MustCompile(
	"(:=|[a-zA-Z1-9]\\(|{|}|[a-z]/[a-z]|\\[)")

func NewSpellcheck(ts TokenSet, ignoreFile string) (sc *Spellcheck, err error) {
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
		toks:    ts,
		speller: NewMultiSpeller(asp, hsp),
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

func (sc *Spellcheck) Check() CheckPipe {
	return func(lc <-chan *Lexeme, outc chan<- *CheckedLexeme) {
		sc.checkComments(lc, outc)
	}
}

func (sc *Spellcheck) checkComments(lc <-chan *Lexeme, outc chan<- *CheckedLexeme) {
	ch := Filter(lc, CommentFilter)
	for comm := range ch {
		if isCodeComment(comm.lit) {
			continue
		}
		if ct := sc.checkLexeme(comm); ct != nil {
			outc <- ct
		}
	}
}

// isCodeComment finds comments that contain example code as indented blocks:
// 	abc, def := indented.with(spaceTab)
//	a, b := tabOnly()
//      if xyz != nil { panic(indentWithSpace) }
// Since it's an example, sometimes the identifiers aren't in the
// code, which will cause a spell check false positive.
func isCodeComment(s string) bool {
	hasPfx := false
	for _, pfx := range []string{"//  ", "// \t", "//\t"} {
		if hasPfx = strings.HasPrefix(s, pfx); hasPfx {
			break
		}
	}
	// Check if line seems to be code to avoid false negatives on
	// indented comments like lists.
	return hasPfx && isCodeRE.MatchString(s)
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
