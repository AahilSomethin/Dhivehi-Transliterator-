package transliterator

import "strings"

// Options configures transliteration features.
type Options struct {
	Gemination          bool // consonant + sukun + same consonant → doubled output
	SuppressGlottalStop bool // suppress apostrophe between adjacent vowels
}

// Fast lookup tables indexed by (r - thaanaBase), replacing map access.
const (
	thaanaBase rune = 0x0780
	thaanaSize      = 0x07B1 - 0x0780 // 49 slots: U+0780 through U+07B0
)

var (
	cLat [thaanaSize]string
	cOk  [thaanaSize]bool
	vLat [thaanaSize]string
	vOk  [thaanaSize]bool
)

func init() {
	for r, s := range ConsonantMap {
		if i := int(r - thaanaBase); i >= 0 && i < thaanaSize {
			cLat[i] = s
			cOk[i] = true
		}
	}
	for r, s := range VowelMap {
		if r == 'ް' {
			continue // sukun handled separately
		}
		if i := int(r - thaanaBase); i >= 0 && i < thaanaSize {
			vLat[i] = s
			vOk[i] = true
		}
	}
}

func consonant(r rune) (string, bool) {
	if i := int(r - thaanaBase); i >= 0 && i < thaanaSize {
		return cLat[i], cOk[i]
	}
	return "", false
}

func vowel(r rune) (string, bool) {
	if i := int(r - thaanaBase); i >= 0 && i < thaanaSize {
		return vLat[i], vOk[i]
	}
	return "", false
}

// Transliterate converts Dhivehi (Thaana) text to Latin.
func Transliterate(input string) string {
	return TransliterateWithOptions(input, Options{})
}

func TransliterateWithOptions(input string, opts Options) string {
	runes := []rune(input)
	n := len(runes)

	var b strings.Builder
	b.Grow(len(input))

	var (
		lastRune     rune
		lastLatin    string
		pending      bool
		posInWord    int
		geminateNext bool
		lastVowel    rune
	)

	for i := 0; i < n; i++ {
		r := runes[i]

		var next rune
		if i+1 < n {
			next = runes[i+1]
		}

		// Whitespace — word boundary reset
		if r <= ' ' && (r == ' ' || r == '\n' || r == '\t') {
			if pending && lastLatin != "" {
				b.WriteString(lastLatin)
			}
			b.WriteRune(r)
			lastRune = 0
			lastLatin = ""
			pending = false
			posInWord = 0
			geminateNext = false
			lastVowel = 0
			continue
		}

		// Sukun
		if r == 'ް' {
			if pending {
				if opts.Gemination && i+1 < n && runes[i+1] == lastRune {
					b.WriteString(lastLatin)
				}

				switch {
				case lastRune == 'ށ':
					b.WriteByte('h')
				case lastRune == 'ނ' && (next == 'ބ' || next == 'ޕ'):
					b.WriteByte('m')
				case lastRune == 'އ':
					if _, ok := consonant(next); ok {
						geminateNext = true
					} else {
						b.WriteByte('h')
					}
				default:
					b.WriteString(lastLatin)
				}
				pending = false
				lastLatin = ""
			}
			continue
		}

		// Vowels (fili)
		if vl, ok := vowel(r); ok {
			if pending {
				if lastRune == 'އ' {
					switch {
					case opts.SuppressGlottalStop:
					case lastVowel != 0 && isDiphthong(lastVowel, r):
					case posInWord == 1:
					default:
						b.WriteByte('\'')
					}
				} else {
					b.WriteString(lastLatin)
				}
				pending = false
			}
			b.WriteString(vl)
			lastVowel = r
			continue
		}

		// Consonants
		if cl, ok := consonant(r); ok {
			if pending && lastLatin != "" {
				b.WriteString(lastLatin)
			}
			pending = false

			lat := cl
			if r == 'ނ' && (next == 'ބ' || next == 'ޕ') {
				lat = "m"
			}

			if geminateNext {
				b.WriteString(lat)
				geminateNext = false
			}

			lastRune = r
			lastLatin = lat
			pending = true
			posInWord++
			continue
		}

		// Fallback: pass through
		if pending && lastLatin != "" {
			b.WriteString(lastLatin)
		}
		pending = false
		b.WriteRune(r)
	}

	// Final flush
	if lastRune == 'އ' && pending {
		b.WriteByte('h')
	} else if pending && lastLatin != "" {
		b.WriteString(lastLatin)
	}

	return b.String()
}

func isDiphthong(prev, curr rune) bool {
	pi := int(prev - thaanaBase)
	ci := int(curr - thaanaBase)
	if pi < 0 || pi >= thaanaSize || ci < 0 || ci >= thaanaSize {
		return false
	}
	p := vLat[pi]
	c := vLat[ci]
	if len(p) == 0 || len(c) == 0 {
		return false
	}
	return (p[len(p)-1] == 'a' || p[len(p)-1] == 'o' || p[len(p)-1] == 'e') &&
		(c[0] == 'i' || c[0] == 'u')
}
