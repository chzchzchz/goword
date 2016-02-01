package main

import (
	"strings"
	"testing"
)

var commChk docCheck

// TestCommentMisspellReject finds a misspeling
func TestCommentMisspellReject(t *testing.T) {
	for _, ct := range commChk.get(t, "comment_test.go").ct {
		if strings.Contains(ct.lit, "TestCommentMisspell") {
			return
		}
	}
	t.Fatal("did not flag misspelling")
}

// TestCommentFuncName has a comment with a function name TestCommentFuncName
func TestCommentFuncNamePass(t *testing.T) {
	for _, ct := range commChk.get(t, "comment_test.go").ct {
		if strings.Contains(ct.lit, "TestCommentFuncName") {
			t.Errorf("unexpected error %v", ct.lit)
		}
	}
}

// TestCommentSomeUserReject should reject direct someuser callouts so TODO works
func TestCommentSomeUserReject(t *testing.T) {
	for _, ct := range commChk.get(t, "comment_test.go").ct {
		if strings.Contains(ct.lit, "TestCommentSomeUserReject") {
			return
		}
	}
	t.Fatal("did not flag floating user")
}

// TestTODOPass should accept TODO(someuser) comments
func TestTODOPass(t *testing.T) {
	for _, ct := range commChk.get(t, "comment_test.go").ct {
		if strings.Contains(ct.lit, "TestTODOPass") {
			t.Errorf("unexpected error %v", ct.lit)
		}
	}
}

// TestTODOSpacePass should accept TODO (someuser) comments
func TestTODOSpacePass(t *testing.T) {
	for _, ct := range commChk.get(t, "comment_test.go").ct {
		if strings.Contains(ct.lit, "TestTODOSpacePass") {
			t.Errorf("unexpected error %v", ct.lit)
		}
	}
}
