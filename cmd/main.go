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

		fmt.Println(transliterator.Transliterate(string(input)))

		// i := 0
		// // fatherTime := time.Duration(0)
		// manOfTheHour := time.Duration(0)

		// var start time.Time
		// for i < 1 {
		// 	// start = time.Now()
		// 	// transliterator.TransliterateV2(string(input))
		// 	// fatherTime += time.Since(start)

		// 	start = time.Now()
		// 	transliterator.Transliterate(string(input))
		// 	manOfTheHour += time.Since(start)

		// 	i++
		// }
		// // fmt.Println("Transliteration without growing took an average of: ", fatherTime/1)
		// fmt.Println("Transliteration with growing took an average of: ", manOfTheHour/1)

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
