package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"dhivehi-translit/internal/transliterator"
)

// CLI Usage Examples:
//
//   # Transliterate from command-line arguments:
//   dhivehi-translit.exe ދިވެހި
//   Output: dhivehi
//
//   # Transliterate multiple words:
//   dhivehi-translit.exe ދިވެހި ބަސް
//   Output: dhivehi bas
//
//   # Transliterate from stdin (pipe or interactive):
//   echo ދިވެހި | dhivehi-translit.exe
//   Output: dhivehi
//
//   # Interactive mode (type input, press Enter):
//   dhivehi-translit.exe
//   > ދިވެހި
//   dhivehi

func main() {
	var input string

	if len(os.Args) > 1 {
		// Read from command-line arguments
		input = strings.Join(os.Args[1:], " ")
	} else {
		// Read from stdin
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			input = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
	}

	if input == "" {
		return
	}

	// Transliterate using v1 default behavior
	result := transliterator.Transliterate(input)
	fmt.Println(result)
}



//With command-line arguments:
// go run ./cmd/dhivehi-translit/ ދިވެހި

//or

//go run ./cmd/dhivehi-translit/ "ދިވެހި ބަހުގެ ދުވަހުން ދުވަހަށް ބޭނުންތަކަށް"

//With stdin (pipe):
//echo ދިވެހި | go run ./cmd/dhivehi-translit/

//To run tests:
//go test ./...

