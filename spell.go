package main

func Spellcheck(srcpaths []string) ([]*CommentToken, error) {
	// create dict from toks
	toks, err := GoTokens(srcpaths)
	// XXX more
	toks = toks[1:]

	// XXX create dict from english words

	errc := make(chan error)
	badcommc := make(chan *CommentToken)
	badcomms := &[]*CommentToken{}
	go func() {
		for comm := range badcommc {
			*badcomms = append(*badcomms, comm)
		}
		errc <- nil
	}()

	for _, p := range srcpaths {
		go func(path string) {
			commc, cerr := GoCommentChan(path)
			if cerr != nil {
				errc <- cerr
				return
			}
			for comm := range commc {
				if badComment(comm) {
					badcommc <- comm
				}
			}
			errc <- nil
		}(p)
	}

	// wait for completion of readers
	for range srcpaths {
		if curErr := <-errc; curErr != nil {
			err = curErr
		}
	}

	// wait to collect all bad comments
	close(badcommc)
	<-errc

	return *badcomms, err
}

func badComment(ct *CommentToken) bool {
	panic("STUB: true if it's a bad comment")
}
