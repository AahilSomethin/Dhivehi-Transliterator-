package transliterator

import (
	"strings"
)

func Transliterate(input string) string {
	runes := []rune(input)
	var result strings.Builder

	i := 0
	for i < len(runes) {
		r := runes[i]

		if akuru, ok := Akuru[r]; ok {
			if i+1 < len(runes) {
				next := runes[i+1]

				if fili, ok := Fili[next]; ok {
					result.WriteString(akuru)
					result.WriteString(fili)
					i += 2
					continue
				}

				if next == Sukun {
					if override, ok := SukunOverrides[r]; ok {
						result.WriteString(override)
					} else {
						result.WriteString(akuru)
					}
					i += 2
					continue
				}
			}

			if r == Noonu && i != 0 && i < len(runes)-1 {
				_, ok1 := Fili[runes[i-1]]
				_, ok2 := Akuru[runes[i+1]]
				if ok1 && ok2 {
					result.WriteString("n'")
					i++
					continue
				}
			}
			if name, ok := AkuruNames[r]; ok {
				result.WriteString(name)
			}

			i++
			continue
		}

		if nishaan, ok := Nishaan[r]; ok {
			result.WriteRune(nishaan)
			i++
			continue
		}

		result.WriteRune(r)
		i++
	}

	return result.String()
}
