# dhivehi-translit

A command-line tool and Go library for transliterating Dhivehi (Thaana script) to Latin script, following the Male Latin romanization standard.

## Features

- Transliterates Thaana consonants and vowels to their Latin equivalents
- Handles sukun (ް), nasalization, gemination, and glottal stops
- Supports Arabic-derived dotted letters (ޝ, ޤ, ޢ, etc.)
- Works as a CLI tool (arguments, file input, or stdin) or as a Go library
- Zero external dependencies — standard library only

## Examples

| Thaana     | Latin    |
| ---------- | -------- |
| ދިވެހި     | dhivehi  |
| ބަސް       | bas      |
| މާލެ       | maale    |
| އަދު       | adhu     |
| ހަމްދު     | hamdhu   |
| އަންބަރަ   | ambara   |
| ޝަރުޠު    | sharuthu |
| ޤައުމު     | qaumu    |

## Installation

Requires **Go 1.22** or later.

```bash
# Install globally
go install ./cmd/dhivehi-translit/

# Or build a binary
go build -o dhivehi-translit ./cmd/dhivehi-translit/
```

## Usage

### CLI

Pass Thaana text as arguments:

```bash
dhivehi-translit ދިވެހި
# Output: dhivehi

dhivehi-translit ދިވެހި ބަސް
# Output: dhivehi bas
```

Transliterate from a text file:

```bash
dhivehi-translit input.txt
# Output: (transliterated contents of input.txt)
```

If the first argument is a path to an existing file, its contents are read and transliterated. Otherwise, arguments are treated as direct text input.

Or pipe text via stdin:

```bash
echo ދިވެހި | dhivehi-translit
# Output: dhivehi
```

Running without arguments enters interactive mode — type a line and press Enter to see the transliteration.

### Library

```go
import "dhivehi-translit/internal/transliterator"

// Basic transliteration
result := transliterator.Transliterate("ދިވެހި") // "dhivehi"

// With options (gemination, suppress glottal stop)
opts := transliterator.Options{
    Gemination:         true,
    SuppressGlottalStop: true,
}
result := transliterator.TransliterateWithOptions("ބައްބަ", opts) // "babba"
```

#### Options

| Option                | Default | Description                                                        |
| --------------------- | ------- | ------------------------------------------------------------------ |
| `Gemination`          | `false` | Double a consonant when sukun is followed by the same consonant    |
| `SuppressGlottalStop` | `false` | Omit the apostrophe (`'`) between adjacent vowels across syllables |

## Running Tests

```bash
go test ./...
```

or

go run ./cmd/dhivehi-translit/ "(what you want to write in Dhivehi)"

Benchmarks:

```bash
go test -bench Benchmark -benchmem ./internal/transliterator/
```

## Project Structure

```
dhivehi-translit/
├── cmd/
│   └── dhivehi-translit/
│       └── main.go                # CLI entry point
├── docs/                          # Reference PDFs
├── internal/
│   └── transliterator/
│       ├── mappings.go            # Consonant & vowel maps
│       ├── rules.go               # Special-case rules
│       ├── state.go               # Transliteration state
│       ├── tokenizer.go           # Rune classification helpers
│       ├── transliterator.go      # Core transliteration logic
│       └── transliterator_test.go # Tests & benchmarks
├── samples/                       # Sample input/output files
├── go.mod
└── README.md
```

## License

See the repository for license information.
