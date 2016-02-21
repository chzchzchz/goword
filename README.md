# goword
Spell checking and godoc checking for golang comments.

## Mechanism
`goword` uses both a natural language dictionary and a source token dictionary to spell check comments.
The source token dictionary avoids flagging acceptable go-style comments (e.g., vim's `:set spell`
will complain about function names).

`goword` understands godoc formatting and can detect incorrect labels in godoc comments. This is
useful for catching when a function, field, or type name drifts from the godoc documentation.

Unlike misspell checkers, `goword` has few false negatives. On the other hand, `goword` is more likely
to give false positives. Most false positives may be elided by passing `goword` an ignore list.

## Example Output

A misspelled local function:
```
/usr/lib/go/src/net/addrselect.go.38: // srcsAddrs tries to UDP-connect to each address to see if it has a (godoc-local: srcsAddrs -> srcAddrs?)
```

A misspelled comment:
```
/usr/lib/go/src/os/str.go.5: // Simple converions to avoid depending on strconv. (spell: converions -> conversions?)
```

An exported function drifted to a local function:
```
/usr/lib/go/src/go/types/conversions.go.11: // Conversion type-checks the conversion T(x). (godoc-local: Conversion -> conversion?)
```

Should begin with the function identifier:
```
/usr/lib/go/src/cmd/compile/internal/gc/builtin/unsafe.go.15: // return types here are ignored; see unsafe.go (godoc-export: return -> Offsetof?)
```

A local function drifted into an export function:
```
/usr/lib/go/src/internal/trace/parser.go.583: // symbolizeTrace attaches func/file/line info to stack traces. (godoc-export: symbolizeTrace -> Symbolize?)
```

## Requirements

`goword` links against:
* [GNU Aspell](http://aspell.net/)
* [Hunspell](http://hunspell.github.io/)

and hence requires the corresponding development headers to build.


## Running

From your golang project root,
```bash
goword file.go [file2.go, file3.go, ...]
```
