// Package maij should trigger TestGoDocPackageNameReject
package main

import (
	"strings"
	"testing"
)

func reject(t *testing.T, filename, funcname string) {
	chk := &docCheck{}
	for _, ct := range chk.get(t, filename).ct {
		if strings.Contains(ct.ctok.lit, funcname) {
			return
		}
	}
	t.Fatal("did not flag bad test " + funcname)
}

func accept(t *testing.T, filename, funcname string) {
	chk := &docCheck{}
	for _, ct := range chk.get(t, filename).ct {
		if strings.Contains(ct.ctok.lit, funcname) {
			t.Errorf("unexpected error %v (%s -> %s?)",
				ct.ctok.lit,
				ct.words[0].word,
				ct.words[0].suggest)
		}
	}
}
