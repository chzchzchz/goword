package main

import (
	"go/scanner"
	"go/token"
	"io/ioutil"
	"os"
)

// GoTokens gets the tokens from set of source files.
func GoTokens(srcpaths []string) (toks []string, err error) {
	tokc := make(chan []string)
	errc := make(chan error)

	fs := token.NewFileSet()
	files := 0
	for _, path := range srcpaths {
		st, serr := os.Stat(path)
		if serr != nil {
			break
		}
		f := fs.AddFile(path, 0, int(st.Size()))
		go func(tf *token.File) {
			s, e := fileTokens(tf)
			tokc <- s
			errc <- e
		}(f)
		files++
	}

	tokmap := make(map[string]struct{})
	for i := 0; i < files; i++ {
		if curToks := <-tokc; curToks != nil {
			for _, tok := range curToks {
				tokmap[tok] = struct{}{}
			}
		}
		if curErr := <-errc; curErr != nil {
			err = curErr
		}
	}

	for tok, _ := range tokmap {
		toks = append(toks, tok)
	}
	return toks, err
}

func fileTokens(tf *token.File) (toks []string, err error) {
	src, err := ioutil.ReadFile(tf.Name())
	if err != nil {
		return nil, err
	}
	s := &scanner.Scanner{}
	s.Init(tf, src, nil, 0)
	tokmap := make(map[string]struct{})
	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		tokmap[lit] = struct{}{}
	}

	for k, _ := range tokmap {
		toks = append(toks, k)
	}
	return toks, nil
}
