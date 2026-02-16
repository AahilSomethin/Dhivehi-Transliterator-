# dhivehi-translit

A command-line tool and Go library for transliterating Dhivehi (Thaana script) to Latin script, following the Male' Latin romanization standard.

## Features

- Two transliteration engines (`v1` and `v2`) with different approaches
- Handles sukun (ް), nasalization, gemination, and glottal stops
- Supports Arabic-derived dotted letters (ޝ, ޤ, ޢ, etc.)
- Works as a CLI tool (file input or stdin) or as a Go library
- Zero external dependencies — standard library only
- Fast performance with array-based lookups and pre-allocated buffers

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
go install ./cmd/...

# Or build a binary
go build -o dhivehi-translit ./cmd/
```

## Usage

### CLI

A version flag (`-v1` or `-v2`) is required. The `v1` engine is the original rule-based implementation; `v2` is an optimized rewrite with additional features like configurable gemination and glottal stop suppression.

**Transliterate a file:**

```bash
dhivehi-translit -v2 input.txt
```

**Pipe text via stdin:**

```bash
echo ދިވެހި | dhivehi-translit -v2
# Output: dhivehi
```

**Interactive mode** — run without a file argument to enter line-by-line mode:

```bash
dhivehi-translit -v1
```

### Library

**v1 — simple transliteration:**

```go
import translit1 "dhivehi-translit/internal/translit1"

result := translit1.Transliterate("ދިވެހި") // "dhivehi"
```

**v2 — with options:**

```go
import translit2 "dhivehi-translit/internal/translit2"

// Basic
result := translit2.Transliterate("ދިވެހި") // "dhivehi"

// With options
opts := translit2.Options{
    Gemination:          true,
    SuppressGlottalStop: true,
}
result := translit2.TransliterateWithOptions("ބައްބަ", opts) // "babba"
```

#### Options (v1 only)

| Option                | Default | Description                                                        |
| --------------------- | ------- | ------------------------------------------------------------------ |
| `Gemination`          | `false` | Double a consonant when sukun is followed by the same consonant    |
| `SuppressGlottalStop` | `false` | Omit the apostrophe (`'`) between adjacent vowels across syllables |

## Running Tests

```bash
go test ./...
```

Benchmarks:

```bash
go test -bench Benchmark -benchmem ./internal/translit1/
go test -bench Benchmark -benchmem ./internal/translit2/
```

## Project Structure

```
dhivehi-translit/
├── cmd/
│   └── main.go                    # CLI entry point (flag parsing, I/O)
├── docs/                          # Reference PDFs
├── internal/
│   ├── translit1/
│   │   ├── engine.go              # v1 transliteration logic
│   │   ├── mappings.go            # v1 consonant & vowel maps
│   │   └── transliterator_test.go # v1 tests & benchmarks
│   └── translit2/
│       ├── transliterator.go      # v2 transliteration logic
│       ├── mappings.go            # v2 character maps & overrides
│       └── transliterator_test.go # v2 tests & benchmarks
├── go.mod
└── README.md
```

## License

See the repository for license information.
