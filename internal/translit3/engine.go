package transliterator

import "strings"

// Options configures transliteration features.
type Options struct {
	Gemination          bool // consonant + sukun + same consonant → doubled output
	SuppressGlottalStop bool // suppress apostrophe between adjacent vowels (no effect in default mode)
	NormalizeArabic     bool // collapse Arabic-derived letters to standard Latin (V1 style)
}

// Array accessor helpers — inlined by the compiler.

func consonant(r rune, norm bool) (string, bool) {
	if i := int(r - thaanaBase); i >= 0 && i < thaanaSize {
		if norm {
			return cLatNorm[i], cOk[i]
		}
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

func sukunOverride(r rune) (string, bool) {
	if i := int(r - thaanaBase); i >= 0 && i < thaanaSize {
		return skOver[i], skOverOk[i]
	}
	return "", false
}

func nishaan(r rune) (rune, bool) {
	if int(r) < len(nishaanLat) {
		return nishaanLat[r], nishaanOk[r]
	}
	return 0, false
}

func isConsonant(r rune) bool {
	if i := int(r - thaanaBase); i >= 0 && i < thaanaSize {
		return cOk[i]
	}
	return false
}

func isVowel(r rune) bool {
	if i := int(r - thaanaBase); i >= 0 && i < thaanaSize {
		return vOk[i]
	}
	return false
}

// Transliterate converts Dhivehi (Thaana) text to Latin with default options.
func Transliterate(input string) string {
	return TransliterateWithOptions(input, Options{})
}

// TransliterateWithOptions converts Dhivehi (Thaana) text to Latin with the given options.
func TransliterateWithOptions(input string, opts Options) string {
	runes := []rune(input)
	n := len(runes)

	var b strings.Builder
	b.Grow(len(input))

	norm := opts.NormalizeArabic

	var (
		lastRune     rune
		lastLatin    string
		pending      bool
		posInWord    int
		geminateNext bool
	)

	for i := 0; i < n; i++ {
		r := runes[i]

		var next rune
		if i+1 < n {
			next = runes[i+1]
		}

		// --- Whitespace: word boundary reset ---
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
			continue
		}

		// --- Punctuation (Nishaan) ---
		if lat, ok := nishaan(r); ok {
			if pending && lastLatin != "" {
				b.WriteString(lastLatin)
			}
			pending = false
			b.WriteRune(lat)
			continue
		}

		// --- Sukun ---
		if r == Sukun {
			if pending {
				// Gemination: consonant + sukun + same consonant → doubled
				if opts.Gemination && i+1 < n && runes[i+1] == lastRune {
					b.WriteString(lastLatin)
				}

				// SukunOverride check (thaalu→"iy", ainu→"u", nyaviyani→"")
				if override, ok := sukunOverride(lastRune); ok {
					b.WriteString(override)
					pending = false
					lastLatin = ""
					continue
				}

				switch {
				case lastRune == Shaviyani:
					// Shaviyani + sukun: if next consonant exists, output its first byte; else "h"
					if i+1 < n {
						if cl, ok := consonant(next, norm); ok && len(cl) > 0 {
							b.WriteByte(cl[0])
						} else {
							b.WriteByte('h')
						}
					} else {
						b.WriteByte('h')
					}

				case lastRune == Alifu:
					// Alifu + sukun: if next is a consonant, geminate; else "h"
					if cl, ok := consonant(next, norm); ok && len(cl) > 0 {
						geminateNext = true
					} else {
						b.WriteByte('h')
					}

				case lastRune == Noonu:
					// Noonu + sukun: nasalization before meemu/baa/paviyani
					if i+1 < n && (next == Meemu || next == Baa || next == Paviyani) {
						if cl, ok := consonant(next, norm); ok {
							b.WriteByte(cl[0])
						} else {
							b.WriteString(lastLatin)
						}
					} else {
						b.WriteString(lastLatin)
					}

				default:
					b.WriteString(lastLatin)
				}
				pending = false
				lastLatin = ""
			}
			continue
		}

		// --- Vowels (Fili) ---
		if vl, ok := vowel(r); ok {
			if pending {
				if lastRune == Alifu {
					// Alifu is a silent carrier — no glottal stop (V2 accuracy)
				} else {
					b.WriteString(lastLatin)
				}
				pending = false
			}
			b.WriteString(vl)
			continue
		}

		// --- Consonants ---
		if cl, ok := consonant(r, norm); ok {
			// Flush pending consonant
			if pending && lastLatin != "" {
				b.WriteString(lastLatin)
			}
			pending = false

			lat := cl

			// Ainu + fili: output first char of fili, then apostrophe, then rest (V2 rule)
			if r == Ainu && i+1 < n {
				if vl, vOk := vowel(next); vOk {
					filiRunes := []rune(vl)
					b.WriteRune(filiRunes[0])
					b.WriteByte('\'')
					if len(filiRunes) > 1 {
						b.WriteString(string(filiRunes[1:]))
					}
					lastRune = r
					lastLatin = ""
					posInWord++
					i++ // skip the vowel; loop's own i++ advances past Ainu
					continue
				}
			}

			// Noonu between fili and next consonant → "n'" (V2 syllable boundary)
			if r == Noonu && i > 0 && i < n-1 {
				if isVowel(runes[i-1]) && isConsonant(runes[i+1]) {
					b.WriteString("n'")
					lastRune = r
					lastLatin = ""
					posInWord++
					continue // loop's i++ advances past noonu; next consonant processed normally
				}
			}

			// Noonu before baa/paviyani → nasalize to "m"
			if r == Noonu && (next == Baa || next == Paviyani) {
				lat = "m"
			}

			// Apply gemination from alifu+sukun
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

		// --- Fallback: pass through ---
		if pending && lastLatin != "" {
			b.WriteString(lastLatin)
		}
		pending = false
		b.WriteRune(r)
	}

	// Final flush
	if pending {
		if lastRune == Alifu {
			b.WriteByte('h')
		} else if lastLatin != "" {
			b.WriteString(lastLatin)
		}
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
