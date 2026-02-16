package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	translit1 "dhivehi-translit/internal/translit1"
	translit2 "dhivehi-translit/internal/translit2"
)

func main() {
	v1 := flag.Bool("v1", false, "use v1 transliterator")
	v2 := flag.Bool("v2", false, "use v2 transliterator")
	flag.Parse()

	if !*v1 && !*v2 {
		fmt.Fprintln(os.Stderr, "error: specify -v1 or -v2")
		flag.Usage()
		os.Exit(1)
	}
	if *v1 && *v2 {
		fmt.Fprintln(os.Stderr, "error: specify only one of -v1 or -v2")
		os.Exit(1)
	}

	transliterate := translit2.Transliterate
	if *v1 {
		transliterate = translit1.Transliterate
	}

	args := flag.Args()
	if len(args) > 0 {
		input, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(transliterate(string(input)))
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Start typing in thaana...")
		for scanner.Scan() {
			fmt.Println(transliterate(scanner.Text()))
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
}
