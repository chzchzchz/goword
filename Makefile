.PHONY: all
all: goword

SRCS=$(filter-out %_test.go, $(wildcard *.go */*.go))
TESTSRCS=$(wildcard *_test.go */*_test.go) 

goword: $(SRCS)
	go build -v

.PHONY: test
test: test.out
	cat test.out

test.out: goword $(TESTSRCS)
	go test -v >$@ 2>&1 || cat $@
