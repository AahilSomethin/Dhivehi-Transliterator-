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
					if r == Ainu {

						filiRunes := []rune(fili)
						result.WriteRune(filiRunes[0])
						result.WriteRune('\'')
						if len(filiRunes) > 1 {
							result.WriteString(string(filiRunes[1:]))
						}
					} else {
						result.WriteString(akuru)
						result.WriteString(fili)
					}
					i += 2
					continue
				}

				if next == Sukun {
					if r == Alifu || r == Shaviyani {

						if i+2 < len(runes) {
							if nextAkuru, isAkuru := Akuru[runes[i+2]]; isAkuru && len(nextAkuru) > 0 {
								result.WriteByte(nextAkuru[0])
								i += 2
								continue
							}
						}

						result.WriteString("h")
						i += 2
						continue
					}

					if r == Noonu {

						if i+2 < len(runes) {
							nextR := runes[i+2]
							if nextR == '\u0789' || nextR == '\u0784' || nextR == '\u0795' {
								if nextAkuru, isAkuru := Akuru[nextR]; isAkuru {
									result.WriteByte(nextAkuru[0])
									i += 2
									continue
								}
							}
						}

						result.WriteString(akuru)
						i += 2
						continue
					}

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

			if r == Raa {
				afterFiliOrAkuru := false
				beforeAkuru := false
				if i > 0 {
					_, pf := Fili[runes[i-1]]
					_, pa := Akuru[runes[i-1]]
					afterFiliOrAkuru = pf || pa
				}
				if i+1 < len(runes) {
					_, na := Akuru[runes[i+1]]
					beforeAkuru = na
				}
				if afterFiliOrAkuru || beforeAkuru {
					result.WriteString("r")
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
