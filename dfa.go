package main

import (
	"go/token"
)

type State int

const (
	Accept = iota - 2
	Reject
	More
)

// xfer represents a state transition from a current state to the next
type xfer map[token.Token]State

// dfa runs a dfa over a lexeme string
func dfa(x []xfer, ll []*Lexeme) State {
	st := State(0)
	for _, l := range ll {
		if int(st) >= len(x) {
			return Reject
		}
		xm := &x[st]
		nextst, ok := (*xm)[l.tok]
		if !ok {
			// check for otherwise case
			nextst, ok = (*xm)[token.ILLEGAL]
			if !ok {
				nextst = Reject
			}
		}
		st = nextst
	}
	return st
}
