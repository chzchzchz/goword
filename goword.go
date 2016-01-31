package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	ct, err := Spellcheck(flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	for _, c := range ct {
		fmt.Println(c)
	}
}
