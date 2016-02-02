package main

import (
	"go/token"
)

// CommentChan streams the comment tokens from a source file.
func CommentChan(srcpath string) (chan *Lexeme, error) {
	ch, err := LexemeChan(srcpath)
	if err != nil {
		return nil, err
	}

	retc := make(chan *Lexeme)
	go func() {
		defer close(retc)
		for l := range ch {
			if l.tok != token.COMMENT {
				l.prev = nil
				continue
			}
			retc <- l
		}
	}()

	return retc, nil
}
