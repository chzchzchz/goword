package main

import (
	"strings"
	"testing"
)

var commChk docCheck

// TestCommentMisspell finds a misspeling
func TestCommentMisspell(t *testing.T) {
	for _, ct := range commChk.get(t, "comment_test.go").ct {
		if strings.Contains(ct.lit, "TestCommentMisspell") {
			return
		}
	}
	t.Fatal("did not flag misspelling")
}

// TestCommentFuncName has a comment with a function name TestCommentFuncName
func TestCommentFuncName(t *testing.T) {
	for _, ct := range commChk.get(t, "comment_test.go").ct {
		if strings.Contains(ct.lit, "TestCommentFuncName") {
			t.Errorf("unexpected error %v", ct.lit)
		}
	}
}
