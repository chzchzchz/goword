package main

import (
	"io/ioutil"
	"strings"
)

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
