package main

import (
	"testing"
)

var stringChk docCheck

var stringInComment = "httpMethod"

// TestStringInCommentPass tests "httpMethod" is accepted
func TestStringInCommentPass(t *testing.T) {
	accept(t, "string_test.go", "TestStringInCommentPass")
}

// TestStringBadCapsReject tests "httpmethod" is rejected for going lower case
func TestStringBadCapsReject(t *testing.T) {
	reject(t, "string_test.go", "TestStringBadCapsReject")
}

var multiWordString = "foo bar baz"

// TestStringTokenizePass tests multiWordString causes acceptance for foo bar baz
func TestStringTokenizePass(t *testing.T) {
	accept(t, "string_test.go", "TestStringTokenizePass")
}
