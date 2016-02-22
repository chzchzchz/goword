package main

import (
	"go/token"
	"strings"
)

func CheckGoDocs(lc <-chan *Lexeme, outc chan<- *CheckedLexeme) {
	tch := Filter(lc, DeclCommentFilter)
	for {
		ll := []*Lexeme{}
		for {
			l, ok := <-tch
			if !ok {
				return
			}
			if l.tok == token.ILLEGAL {
				break
			}
			ll = append(ll, l)
		}

		godoc := godocBlock(ll)

		// does the comment line up with the next line?
		after := afterGoDoc(ll)
		if after.pos.Column != godoc[0].pos.Column {
			continue
		}
		// is the comment on the line immediately before the code?
		if after.pos.Line != godoc[len(godoc)-1].pos.Line+1 {
			continue
		}

		// does the comment have a token for documentation?
		fields := strings.Fields(godoc[0].lit)
		if len(fields) < 2 {
			continue
		}

		// what token should the documentation match?
		cmplex := ll[len(ll)-1]
		if len(ll) >= 2 && ll[len(ll)-2].tok == token.IDENT {
			cmplex = ll[len(ll)-2]
		}
		if fields[1] == cmplex.lit {
			continue
		}

		// bad godoc
		label := "godoc-local"
		if strings.ToUpper(cmplex.lit)[0] == cmplex.lit[0] {
			label = "godoc-export"
		}
		cw := []CheckedWord{{fields[1], cmplex.lit}}
		cl := &CheckedLexeme{godoc[0], label, cw}
		outc <- cl
	}
}

// godocBlock gets the godoc comment block from a comment prefixed token string
func godocBlock(ll []*Lexeme) (comm []*Lexeme) {
	wantLine := 0
	for _, l := range ll {
		if l.tok != token.COMMENT {
			break
		}
		if l.pos.Line != wantLine {
			comm = []*Lexeme{}
		}
		wantLine = l.pos.Line + 1
		comm = append(comm, l)
	}
	return comm
}

// afterGoDoc gets the first token following the comments
func afterGoDoc(ll []*Lexeme) *Lexeme {
	for _, l := range ll {
		if l.tok != token.COMMENT {
			return l
		}
	}
	return nil
}
