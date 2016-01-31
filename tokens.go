package main

import (
	"go/scanner"
	"go/token"
	"io/ioutil"
	"os"
)

// GoTokens gets the tokens from set of source files.
func GoTokens(srcpaths []string) (toks map[string]struct{}, err error) {
	tokc := make(chan []string)
	errc := make(chan error)

	fs := token.NewFileSet()
	files := 0
	for _, path := range srcpaths {
		st, serr := os.Stat(path)
		if serr != nil {
			break
		}
		f := fs.AddFile(path, fs.Base(), int(st.Size()))
		go func(tf *token.File) {
			s, e := fileTokens(tf)
			tokc <- s
			errc <- e
		}(f)
		files++
	}

	toks = make(map[string]struct{})
	for i := 0; i < files; i++ {
		if curToks := <-tokc; curToks != nil {
			for _, tok := range curToks {
				toks[tok] = struct{}{}
			}
		}
		if curErr := <-errc; curErr != nil {
			err = curErr
		}
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
		if tok == token.STRING {
			// XXX: what if strings are misspelled?
			lit = lit[1 : len(lit)-1]
		}
		tokmap[lit] = struct{}{}
	}

	for k, _ := range tokmap {
		toks = append(toks, k)
	}
	return toks, nil
}
