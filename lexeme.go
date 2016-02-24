package main

import (
	"go/scanner"
	"go/token"
	"io/ioutil"
	"os"
)

type Lexeme struct {
	pos token.Position
	p   token.Pos
	tok token.Token
	lit string
}

// LexemeChan streams a source file
func LexemeChan(srcpath string) (<-chan *Lexeme, error) {
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
		defer close(lexc)
		for {
			p, t, lit := s.Scan()
			if t == token.EOF {
				return
			}
			lexeme := &Lexeme{tf.Position(p), p, t, lit}
			lexc <- lexeme
		}
	}()

	return lexc, nil
}

func Filter(lc <-chan *Lexeme, f func([]*Lexeme) State) <-chan *Lexeme {
	retc := make(chan *Lexeme)
	go func() {
		ls := []*Lexeme{}
		for l := range lc {
			ls = append(ls, l)
			st := f(ls)
			if st == Reject {
				ls = nil
				continue
			}
			if st != Accept {
				continue
			}
			for _, v := range ls {
				retc <- v
			}
			if len(ls) > 1 {
				retc <- &Lexeme{tok: token.ILLEGAL}
			}
			ls = nil
		}
		close(retc)
	}()
	return retc
}

func CommentFilter(l []*Lexeme) State {
	return dfa([]xfer{
		{token.COMMENT: Accept},
	}, l)
}

// DeclRootCommentFilter gives a comment header for types and functions.
func DeclRootCommentFilter(l []*Lexeme) State {
	return dfa([]xfer{
		{token.COMMENT: 1},
		{token.COMMENT: 1, token.TYPE: 2, token.FUNC: 3, token.IDENT: 5},
		{token.IDENT: Accept},
		{token.IDENT: Accept, token.LPAREN: 4},
		{token.RPAREN: 5, token.ILLEGAL: 4},
		{token.IDENT: Accept},
	}, l)
}

// DeclTypeFilter captures the contents of types.
func DeclTypeFilter(l []*Lexeme) State {
	return dfa([]xfer{
		{token.TYPE: 1},
		{token.IDENT: 2},
		{token.INTERFACE: 3, token.STRUCT: 3},
		{token.LBRACE: 4},
		{token.RBRACE: Accept, token.ILLEGAL: 4},
	}, l)
}

// DeclIdentCommentFilter captures comments preceding an identifier.
func DeclIdentCommentFilter(l []*Lexeme) State {
	return dfa([]xfer{
		{token.COMMENT: 1},
		{token.COMMENT: 1, token.IDENT: Accept},
	}, l)
}

func LexemeMux(lc <-chan *Lexeme, n int) []chan *Lexeme {
	ret := []chan *Lexeme{}
	for i := 0; i < n; i++ {
		ret = append(ret, make(chan *Lexeme, 32))
	}
	go func() {
		for l := range lc {
			for i := 0; i < n; i++ {
				ret[i] <- l
			}
		}
		for i := 0; i < n; i++ {
			close(ret[i])
		}
	}()

	return ret
}
