// ackage should trigger TestGoDocPackageReject
package main

import (
	"testing"
)

type docCheck struct {
	ct []*CheckedLexeme
}

func (dc *docCheck) get(t *testing.T, path string) *docCheck {
	if len(dc.ct) == 0 {
		cts, err := CheckAll([]string{path})
		if err != nil {
			t.Fatal(err)
		}
		dc.ct = cts
	}
	return dc
}

// this will trigger a warning that TestGoDocNoFuncNameReject isn't the function name
func TestGoDocNoFuncNameReject(t *testing.T) {
	reject(t, "godoc_test.go", "TestGoDocNoFuncNameReject")
}

// TestGoDocNoFuncNameReject has the wrong name; want TestGoDocWrongFuncNameReject.
func TestGoDocWrongFuncNameReject(t *testing.T) {
	reject(t, "godoc_test.go", "TestGoDocWrongFuncNameReject")
}

// this will trigger TestGoDocNoTypeNameReject
type godocStructReject struct {
	// this will trigger TestGoDocNoFieldNameReject
	oopsie int
	// okthough will pass TestGoDocFieldNamePass
	okthough int
}

// TestGoDocNoTypeNameReject rejects a type missing name in comments
func TestGoDocNoTypeNameReject(t *testing.T) {
	reject(t, "godoc_test.go", "TestGoDocNoFuncNameReject")
}

// TestGoDocNoFieldNameReject rejects a field missing name in comments
func TestGoDocNoFieldNameReject(t *testing.T) {
	reject(t, "godoc_test.go", "TestGoDocNoFieldNameReject")
}

// TestGoDocFieldNamePass accepts a godoc of a field
func TestGoDocFieldNamePass(t *testing.T) {
	reject(t, "godoc_test.go", "TestGoDocFieldNamePass")
}

// TestGoDocFuncPass will not trigger a warning
func TestGoDocFuncPass(t *testing.T) {
	accept(t, "godoc_test.go", "TestGoDocFuncPass")
}

// TestGoDocMultiLinePass should pass
// even though this is a multiple line comment
func TestGoDocMultiLinePass(t *testing.T) {
	accept(t, "godoc_test.go", "TestGoDocMultiLinePass")
}

// this comment not part of the function documentation

// TestGoDocCommentBreakPass will not trigger a warning
func TestGoDocCommentBreakPass(t *testing.T) {
	accept(t, "godoc_test.go", "TestGoDocCommentBreak")
}

type godocSideComment struct {
	sideComment int // a side comment isn't a godoc
	unrelated   int
}

func TestGoDocSideCommentPass(t *testing.T) {
	accept(t, "godoc_test.go", "TestGoDocSideCommentPass")
}

// TestGoDocSkipLinesPass will not trigger a warning
//
// even though there are multiple lines
// that could
//
// break the checker's
// ability to parse
func TestGoDocSkipLinesPass(t *testing.T) {
	accept(t, "godoc_test.go", "TestGoDocSkipLinesPass")
}

// it's OK not to match TestGoDocSkipAndMissingPass since blank below

func TestGoDocSkipAndMissingPass(t *testing.T) {
	accept(t, "godoc_test.go", "TestGoDocSkipAndMissingPass")
}

type hello interface {
	// this should trigger TestInterfaceFuncReject
	f() int
}

func TestInterfaceFuncReject(t *testing.T) {
	reject(t, "godoc_test.go", "TestInterfaceFuncReject")
}

type morestruct struct {
	// this should trigger TestStructFieldReject
	hello int
	// this should trigger TestStructFieldPtrReject
	hello2 *int
	// this should trigger TestStructFieldChanReject
	hello3 chan int
	// this should trigger TestStructFieldRChanReject
	hello4 <-chan int
}

func TestStructFieldReject(t *testing.T) {
	reject(t, "godoc_test.go", "TestStructFieldReject")
}
func TestStructFieldPtrReject(t *testing.T) {
	reject(t, "godoc_test.go", "TestStructFieldPtrReject")
}
func TestStructFieldChanReject(t *testing.T) {
	reject(t, "godoc_test.go", "TestStructFieldChanReject")
}
func TestStructFieldRChanReject(t *testing.T) {
	reject(t, "godoc_test.go", "TestStructFieldRChanReject")
}

func TestGoDocPackageReject(t *testing.T) {
	reject(t, "godoc_test.go", "TestGoDocPackageReject")
}

func TestGoDocPackageNameReject(t *testing.T) {
	reject(t, "util_test.go", "TestGoDocPackageNameReject")
}
