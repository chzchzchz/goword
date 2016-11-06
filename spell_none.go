// +build !spell

package main

type Spellcheck struct {
	check CheckFunc
}

func NewSpellcheck(ts TokenSet, ignoreFile string) (*Spellcheck, error) {
	return &Spellcheck{func(w string) bool { return true }}, nil
}

func (sc *Spellcheck) Close() {}

func (sc *Spellcheck) WithPassTokens() CheckFunc {
	return func(string) bool { return true }
}

func (sc *Spellcheck) WithSpeller() CheckFunc {
	return func(string) bool { return true }
}

func (sc *Spellcheck) Check() CheckPipe {
	return func(lc <-chan *Lexeme, outc chan<- *CheckedLexeme) {
		for range lc {
		}
	}
}
