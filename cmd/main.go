package main

import (
	"bufio"
	"fmt"
	"os"

	transliterator "dhivehi-translit/internal"
)

func main() {
	if len(os.Args) > 1 {
		input, err := os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(transliterator.Transliterate(string(input)))
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Println(transliterator.Transliterate(scanner.Text()))
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
}
