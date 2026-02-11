package transliterator

import "strings"

func Transliterate(input string) string {
	runes := []rune(input)
	var result strings.Builder

	i := 0
	for i < len(runes) {
		if i+1 < len(runes) {
			pair := string(runes[i : i+2])
			if latin, ok := Mappings[pair]; ok {
				result.WriteString(latin)
				i += 2
				continue
			}
		}

		result.WriteRune(runes[i])
		i++
	}

	return result.String()
}
