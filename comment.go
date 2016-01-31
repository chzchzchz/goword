package main

import (
	"go/scanner"
	"go/token"
	"io/ioutil"
	"os"
)

type CommentToken struct {
	pos token.Position
	p   token.Pos
	tok token.Token
	lit string
}

// GoCommentChan streams the comment tokens from a source file.
func GoCommentChan(srcpath string) (chan *CommentToken, error) {
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

	commc := make(chan *CommentToken)
	go func() {
		for {
			p, t, l := s.Scan()
			if t == token.EOF {
				break
			}
			if t != token.COMMENT {
				continue
			}
			// TODO know comment's context
			// * before function
			// * before type
			// * before global
			// * before struct field
			// * inside function
			commc <- &CommentToken{tf.Position(p), p, t, l}
		}
		close(commc)
	}()

	return commc, nil
}
