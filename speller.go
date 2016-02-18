package main

import (
	"sync"

	"github.com/akhenakh/hunspellgo"
	"github.com/trustmaster/go-aspell"
)

type Speller interface {
	Check(w string) bool
	Suggest(w string) []string
	Close()
}

type aspeller struct {
	sp aspell.Speller
	mu sync.Mutex
}

func NewASpeller() (Speller, error) {
	opts := map[string]string{
		"lang":           "en",
		"filter":         "url",
		"mode":           "url",
		"encoding":       "ascii",
		"guess":          "true",
		"ignore":         "0",
		"ignore-case":    "false",
		"ignore-accents": "false",
	}

	sp, err := aspell.NewSpeller(opts)
	if err != nil {
		return nil, err
	}
	return &aspeller{sp: sp}, nil
}

func (s *aspeller) Check(w string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sp.Check(w)
}

func (s *aspeller) Suggest(w string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sp.Suggest(w)
}

func (s *aspeller) Close() { s.sp.Delete() }

type hunspeller struct {
	sp *hunspellgo.Hunhandle
}

func NewHunSpeller() Speller {
	sp := hunspellgo.Hunspell("/usr/share/hunspell", "/usr/share/hunspell")
	if sp == nil {
		return nil
	}
	return &hunspeller{sp}
}

func (s *hunspeller) Check(w string) bool {
	return s.sp.Spell(w)
}

func (s *hunspeller) Suggest(w string) []string {
	return s.sp.Suggest(w)
}

func (s *hunspeller) Close() { s.sp = nil }

type multispeller struct {
	sp []Speller
}

func NewMultiSpeller(sp ...Speller) Speller {
	m := &multispeller{}
	for _, s := range sp {
		m.sp = append(m.sp, s)
	}
	return m
}

func (s *multispeller) Check(w string) bool {
	for _, sp := range s.sp {
		if sp.Check(w) {
			return true
		}
	}
	return false
}

func (s *multispeller) Suggest(w string) (ret []string) {
	for _, sp := range s.sp {
		ret = append(ret, sp.Suggest(w)...)
	}
	return ret
}

func (s *multispeller) Close() {
	for _, sp := range s.sp {
		sp.Close()
	}
	s.sp = nil
}
