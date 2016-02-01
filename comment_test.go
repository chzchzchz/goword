package main

import (
	"testing"
)

var commChk docCheck

// TestCommentMisspellReject finds a misspeling
func TestCommentMisspellReject(t *testing.T) {
	reject(t, "comment_test.go", "TestCommentMisspell")
}

// TestCommentFuncName has a comment with a function name TestCommentFuncName
func TestCommentFuncNamePass(t *testing.T) {
	accept(t, "comment_test.go", "TestCommentFuncName")
}

// TestCommentSomeUserReject should reject direct someuser callouts so testing TODO works
func TestCommentSomeUserReject(t *testing.T) {
	reject(t, "comment_test.go", "TestCommentSomeUserReject")
}

// TestTODOPass should accept TODO(someuser) comments
func TestTODOPass(t *testing.T) {
	accept(t, "comment_test.go", "TestTODOPass")
}

// TestTODOSpacePass should accept TODO (someuser) comments
func TestTODOSpacePass(t *testing.T) {
	accept(t, "comment_test.go", "TestTODOSpacePass")
}
