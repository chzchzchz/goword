package main

import (
	"strings"
	"testing"
)

var godocChk docCheck

type docCheck struct {
	ct []*CommentToken
}

func (dc *docCheck) get(t *testing.T, path string) *docCheck {
	if len(dc.ct) == 0 {
		sp, err := NewSpellcheck("")
		if err != nil {
			t.Fatal(err)
		}
		defer sp.Close()
		ct, err := sp.Check([]string{path})
		if err != nil {
			t.Fatal(err)
		}
		dc.ct = ct
	}
	return dc
}

// this will trigger a warning that the first word isn't the function name
func TestGoDocBadFirstWord(t *testing.T) {
	for _, ct := range godocChk.get(t, "godoc_test.go").ct {
		if strings.Contains(ct.lit, "this will trigger") {
			return
		}
	}
	t.Fatal("did not flag bad godoc")
}

// TestGoDocGood will not trigger a warning
func TestGoDocGood(t *testing.T) {
	for _, ct := range godocChk.get(t, "godoc_test.go").ct {
		if strings.Contains(ct.lit, "TestGoodGoDoc will") {
			t.Fatalf("unexpected error %v", ct.lit)
		}
	}
}
