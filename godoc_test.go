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

// this will trigger a warning that TestGoDocNoFuncNameReject isn't the function name
func TestGoDocNoFuncNameReject(t *testing.T) {
	goDocReject(t, "TestGoDocNoFuncNameReject")
}

// TestGoDocNoFuncNameReject has the wrong function name
func TestGoDocWrongFuncNameReject(t *testing.T) {
	goDocReject(t, "TestGoDocWrongFuncNameReject")
}

// this will trigger TestGoDocNoTypeNameReject
type godocStructReject struct {
	// this will trigger TestGoDocNoFieldNameReject
	oopsie int
}

// TestGoDocNoTypeNameReject rejects a type missing name in comments
func TestGoDocNoTypeNameReject(t *testing.T) {
	goDocReject(t, "TestGoDocNoFuncNameReject")
}

// TestGoDocNoFieldNameReject rejects a field missing name in comments
func TestGoDocNoFieldNameReject(t *testing.T) {
	goDocReject(t, "TestGoDocNoFieldNameReject")
}

// TestGoDocFuncPass will not trigger a warning
func TestGoDocFuncPass(t *testing.T) {
	for _, ct := range godocChk.get(t, "godoc_test.go").ct {
		if strings.Contains(ct.lit, "TestGoDocFuncPass") {
			t.Fatalf("unexpected error %v", ct.lit)
		}
	}
}

func goDocReject(t *testing.T, f string) {
	for _, ct := range godocChk.get(t, "godoc_test.go").ct {
		if strings.Contains(ct.lit, f) {
			return
		}
	}
	t.Fatal("did not flag bad godoc")
}
