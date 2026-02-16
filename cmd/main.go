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
	translit4 "dhivehi-translit/internal/translit4"
)

func main() {
	v1 := flag.Bool("v1", false, "use v1 engine")
	v2 := flag.Bool("v2", false, "use v2 engine")
	v3 := flag.Bool("v3", false, "use v3 engine")
	v4 := flag.Bool("v4", false, "use v4 engine (default)")
	timer := flag.Bool("timer", false, "print transliteration runtime to stderr")
	shortTimer := flag.Bool("t", false, "shorthand for -timer")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: dhivehi-translit [flags] [file]\n\n")
		fmt.Fprintf(os.Stderr, "Transliterate Dhivehi (Thaana) text to Latin script.\n\n")
		fmt.Fprintf(os.Stderr, "If a file path is given, its contents are transliterated to stdout.\n")
		fmt.Fprintf(os.Stderr, "Otherwise reads line-by-line from stdin (interactive or piped).\n\n")
		fmt.Fprintf(os.Stderr, "Engine:\n")
		fmt.Fprintf(os.Stderr, "  -v1    use v1 engine\n")
		fmt.Fprintf(os.Stderr, "  -v2    use v2 engine\n")
		fmt.Fprintf(os.Stderr, "  -v3    use v3 engine\n")
		fmt.Fprintf(os.Stderr, "  -v4    use v4 engine (default)\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -t, -timer    print transliteration runtime to stderr\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  dhivehi-translit input.txt\n")
		fmt.Fprintf(os.Stderr, "  echo \"ދިވެހި\" | dhivehi-translit\n")
		fmt.Fprintf(os.Stderr, "  dhivehi-translit -v2 input.txt\n")
	}

	flag.Parse()

	showTimer := *timer || *shortTimer

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
	if *v4 {
		vCount++
	}
	if vCount > 1 {
		fmt.Fprintln(os.Stderr, "error: specify only one of -v1, -v2, -v3, or -v4")
		os.Exit(1)
	}

	var transliterate func(string) string
	engineName := "v4"

	switch {
	case *v1:
		transliterate = translit1.Transliterate
		engineName = "v1"
	case *v2:
		transliterate = translit2.Transliterate
		engineName = "v2"
	case *v3:
		transliterate = translit3.Transliterate
		engineName = "v3"
	default:
		transliterate = translit4.Transliterate
	}

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
