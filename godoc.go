package main

import (
	"go/token"
	"strings"
)

func checkGoDocs(lc <-chan *Lexeme, outc chan<- *CheckedLexeme) {
	tch := Filter(lc, DeclCommentFilter)
	for {
		comm, ok := <-tch
		if !ok {
			return
		}
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
		fields := strings.Fields(comm.lit)
		if len(fields) < 2 {
			continue
		}

		cmplex := ll[len(ll)-1]
		if len(ll) >= 2 && ll[len(ll)-2].tok == token.IDENT {
			cmplex = ll[len(ll)-2]
		}
		if fields[1] == cmplex.lit {
			continue
		}
		cw := []CheckedWord{{fields[1], cmplex.lit}}
		cl := &CheckedLexeme{comm, "godoc", cw}
		outc <- cl
	}
}
