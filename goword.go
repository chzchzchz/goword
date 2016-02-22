package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func newSpellCheck(srcpaths []string, ignfile string) (*Spellcheck, error) {
	toks, err := GoTokens(srcpaths)
	if err != nil {
		return nil, err
	}
	splitToks := make(map[string]struct{})
	for k, _ := range toks {
		for _, field := range strings.Fields(k) {
			splitToks[field] = struct{}{}
		}
	}
	sp, serr := NewSpellcheck(splitToks, ignfile)
	if serr != nil {
		return nil, serr
	}
	return sp, nil
}

func main() {
	ign := flag.String("ignore-file", "", "additional words to ignore")
	useSpell := flag.Bool("use-spell", true, "check spelling")
	useGoDoc := flag.Bool("use-godoc", true, "check godocs")
	flag.Parse()

	var cps []CheckPipe

	if *useSpell {
		ignfile := ""
		if ign != nil {
			ignfile = *ign
		}
		sp, serr := newSpellCheck(flag.Args(), ignfile)
		if serr != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", serr)
			os.Exit(1)
		}
		defer sp.Close()
		cps = append(cps, sp.Check())
	}

	if *useGoDoc {
		cps = append(cps, CheckGoDocs)
	}

	// find all errors
	ct, cerr := Check(flag.Args(), cps)
	if cerr != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", cerr)
		os.Exit(1)
	}
	for _, c := range ct {
		fmt.Printf("%s.%d: %s (%s: %s -> %s?)\n",
			c.ctok.pos.Filename,
			c.ctok.pos.Line,
			c.ctok.lit,
			c.rule,
			c.words[0].word,
			c.words[0].suggest)
	}
}
