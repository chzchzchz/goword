package main

import (
	"testing"
)

var commChk docCheck

// TestCommentMisspellReject finds a misspeling
func TestCommentMisspellReject(t *testing.T) {
	reject(t, "comment_test.go", "TestCommentMisspell")
}

// TestCommentFuncNamePass has a comment with a function name TestCommentFuncNamePass
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

// TestPackage_underscoreTest checks the spell checker will not
// tokenize identifiers with underscores; the individual components
// of the function name will not pass spell check, but taken together
// they match a function name pulled from the source.
func TestPackage_underscoreTest(t *testing.T) {
	accept(t, "comment_test.go", "TestPackage_underscoreTest")
}
