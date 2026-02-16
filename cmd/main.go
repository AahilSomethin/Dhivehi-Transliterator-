package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	translit1 "dhivehi-translit/internal/translit1"
	translit2 "dhivehi-translit/internal/translit2"
	translit3 "dhivehi-translit/internal/translit3"
)

func main() {
	// Engine selection
	v1 := flag.Bool("v1", false, "use v1 engine")
	v2 := flag.Bool("v2", false, "use v2 engine")
	v3 := flag.Bool("v3", false, "use v3 engine (default)")

	// Option flags
	gemination := flag.Bool("gemination", false, "consonant + sukun + same consonant → doubled (v1, v3)")
	suppressGlottal := flag.Bool("suppress-glottal-stop", false, "suppress glottal-stop apostrophe (v1)")
	normalizeArabic := flag.Bool("normalize-arabic", false, "collapse Arabic-derived letters to Latin (v3)")

	// Timer flag
	timer := flag.Bool("timer", false, "print transliteration runtime to stderr")
	shortTimer := flag.Bool("t", false, "shorthand for -timer")

	// Short flags
	shortGem := flag.Bool("g", false, "shorthand for -gemination")
	shortSuppress := flag.Bool("s", false, "shorthand for -suppress-glottal-stop")
	shortNorm := flag.Bool("n", false, "shorthand for -normalize-arabic")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: dhivehi-translit [flags] [file]\n\n")
		fmt.Fprintf(os.Stderr, "Transliterate Dhivehi (Thaana) text to Latin script.\n\n")
		fmt.Fprintf(os.Stderr, "If a file path is given, its contents are transliterated to stdout.\n")
		fmt.Fprintf(os.Stderr, "Otherwise reads line-by-line from stdin (interactive or piped).\n\n")
		fmt.Fprintf(os.Stderr, "Engine:\n")
		fmt.Fprintf(os.Stderr, "  -v1                        use v1 engine\n")
		fmt.Fprintf(os.Stderr, "  -v2                        use v2 engine\n")
		fmt.Fprintf(os.Stderr, "  -v3                        use v3 engine (default)\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -g, -gemination            consonant + sukun + same consonant → doubled (v1, v3)\n")
		fmt.Fprintf(os.Stderr, "  -s, -suppress-glottal-stop suppress glottal-stop apostrophe (v1)\n")
		fmt.Fprintf(os.Stderr, "  -n, -normalize-arabic      collapse Arabic-derived letters to Latin (v3)\n")
		fmt.Fprintf(os.Stderr, "  -t, -timer                 print transliteration runtime to stderr\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  dhivehi-translit input.txt\n")
		fmt.Fprintf(os.Stderr, "  echo \"ދިވެހި\" | dhivehi-translit\n")
		fmt.Fprintf(os.Stderr, "  dhivehi-translit -v1 -g input.txt\n")
		fmt.Fprintf(os.Stderr, "  dhivehi-translit -n -gemination input.txt\n")
	}

	flag.Parse()

	// Merge short and long flags
	gem := *gemination || *shortGem
	suppress := *suppressGlottal || *shortSuppress
	normalize := *normalizeArabic || *shortNorm
	showTimer := *timer || *shortTimer

	// Ensure at most one engine is selected
	vCount := 0
	if *v1 {
		vCount++
	}
	if *v2 {
		vCount++
	}
	if *v3 {
		vCount++
	}
	if vCount > 1 {
		fmt.Fprintln(os.Stderr, "error: specify only one of -v1, -v2, or -v3")
		os.Exit(1)
	}

	// Validate option compatibility
	if *v2 && (gem || suppress || normalize) {
		fmt.Fprintln(os.Stderr, "error: v2 engine does not support any option flags")
		os.Exit(1)
	}
	if *v1 && normalize {
		fmt.Fprintln(os.Stderr, "error: -normalize-arabic is only supported by v3")
		os.Exit(1)
	}

	// Build the transliterate function (default: v3)
	var transliterate func(string) string

	switch {
	case *v1:
		opts := translit1.Options{
			Gemination:          gem,
			SuppressGlottalStop: suppress,
		}
		transliterate = func(s string) string {
			return translit1.TransliterateWithOptions(s, opts)
		}

	case *v2:
		transliterate = translit2.Transliterate

	default: // v3
		opts := translit3.Options{
			Gemination:          gem,
			SuppressGlottalStop: suppress,
			NormalizeArabic:     normalize,
		}
		transliterate = func(s string) string {
			return translit3.TransliterateWithOptions(s, opts)
		}
	}

	// Determine engine name for timer output
	engineName := "v3"
	if *v1 {
		engineName = "v1"
	} else if *v2 {
		engineName = "v2"
	}

	// Process input
	args := flag.Args()
	if len(args) > 0 {
		input, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		start := time.Now()
		result := transliterate(string(input))
		elapsed := time.Since(start)

		fmt.Print(result)
		if showTimer {
			fmt.Fprintf(os.Stderr, "[%s] %v (%.3f ms)\n", engineName, elapsed, float64(elapsed.Nanoseconds())/1e6)
		}
	} else {
		fi, _ := os.Stdin.Stat()
		if (fi.Mode() & os.ModeCharDevice) != 0 {
			eofHint := "Ctrl+D"
			if runtime.GOOS == "windows" {
				eofHint = "Ctrl+Z then Enter"
			}
			fmt.Fprintf(os.Stderr, "Type Thaana text (%s to exit):\n", eofHint)
		}

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			start := time.Now()
			result := transliterate(scanner.Text())
			elapsed := time.Since(start)

			fmt.Println(result)
			if showTimer {
				fmt.Fprintf(os.Stderr, "[%s] %v (%.3f ms)\n", engineName, elapsed, float64(elapsed.Nanoseconds())/1e6)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", err)
			os.Exit(1)
		}
	}
}
