package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	ign := flag.String("ignore-file", "", "additional words to ignore")
	flag.Parse()

	ignfile := ""
	if ign != nil {
		ignfile = *ign
	}
	sp, serr := NewSpellcheck(ignfile)
	if serr != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", serr)
		os.Exit(1)
	}
	defer sp.Close()

	ct, cerr := sp.Check(flag.Args())
	if cerr != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", cerr)
		os.Exit(1)
	}
	for _, c := range ct {
		sugg := sp.Suggest(c.lit)
		fmt.Printf("%s.%d: %s (%s?)\n", c.pos.Filename, c.pos.Line, c.lit, sugg)
	}
}
