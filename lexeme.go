package main

import (
	"go/scanner"
	"go/token"
	"io/ioutil"
	"os"
)

type Lexeme struct {
	pos  token.Position
	p    token.Pos
	tok  token.Token
	lit  string
	prev *Lexeme
}

// LexemeChan streams a source file
func LexemeChan(srcpath string) (chan *Lexeme, error) {
	fs := token.NewFileSet()
	st, err := os.Stat(srcpath)
	if err != nil {
		return nil, err
	}
	tf := fs.AddFile(srcpath, fs.Base(), int(st.Size()))
	src, err := ioutil.ReadFile(tf.Name())
	if err != nil {
		return nil, err
	}
	s := &scanner.Scanner{}
	s.Init(tf, src, nil, scanner.ScanComments)

	lexc := make(chan *Lexeme)
	go func() {
		var prev *Lexeme
		defer close(lexc)
		for {
			p, t, lit := s.Scan()
			if t == token.EOF {
				return
			}
			lexeme := &Lexeme{tf.Position(p), p, t, lit, prev}
			lexc <- lexeme
			prev = lexeme
		}
	}()

	return lexc, nil
}

func Filter(lc chan *Lexeme, f func(*Lexeme) bool) chan *Lexeme {
	retc := make(chan *Lexeme)
	go func() {
		for l := range lc {
			if f(l) {
				retc <- l
			} else {
				l.prev = nil
			}
		}
		close(retc)
	}()
	return retc
}

// CommentChan streams the comment tokens from a source file.
func CommentChan(srcpath string) (chan *Lexeme, error) {
	ch, err := LexemeChan(srcpath)
	if err != nil {
		return nil, err
	}
	return Filter(ch,
		func(l *Lexeme) bool {
			return l.tok == token.COMMENT
		}), nil
}
