# goword
Spell checking for golang comments.

## Mechanism
`goword` uses both a natural language dictionary and a source token dictionary to spell check comments.
The source token dictionary avoids flagging acceptable go-style comments (e.g., vim's `:set spell`
will complain about function names).

Unlike misspell checkers, `goword` has few false negatives. On the other hand, `goword` is more likely
to give false positives.

## Running

From your golang project root,
```bash
goword file.go [file2.go, file3.go, ...]
```
