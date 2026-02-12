package transliterator

type Options struct {
	Gemination          bool // consonant + sukun + same consonant → doubled output
	SuppressGlottalStop bool // suppress apostrophe between adjacent vowels
}

// DefaultOptions returns v1-compatible options (all features disabled).
func DefaultOptions() Options {
	return Options{}
}

// Transliterate converts Dhivehi (Thaana) text to Latin using v1 default behavior.
func Transliterate(input string) string {
	return TransliterateWithOptions(input, DefaultOptions())
}

// TransliterateWithOptions converts Dhivehi (Thaana) text to Latin with configurable v2 features.
func TransliterateWithOptions(input string, opts Options) string {
	runes := []rune(input)
	var output []rune
	state := State{}

	// flush writes the buffered consonant to output if it hasn't been paired with a vowel
	flush := func() {
		if state.HasPendingConsonant {
			if state.LastLatin != "" {
				output = append(output, []rune(state.LastLatin)...)
			}
			state.HasPendingConsonant = false
		}
	}

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		var nextRune rune
		if i+1 < len(runes) {
			nextRune = runes[i+1]
		}

		// Whitespace
		if IsWhitespace(r) {
			flush()
			state = State{}
			output = append(output, r)
			continue
		}


		// Sukūn 
		if IsSukun(r) {
			if state.HasPendingConsonant {
				// v2: Gemination - if consonant + sukun + same consonant, double the output
				if opts.Gemination && i+1 < len(runes) && runes[i+1] == state.LastRune {
					output = append(output, []rune(state.LastLatin)...)
					// The next consonant iteration will add another instance
				}

				if state.LastRune == 'ށ' {
					// ށ + sukun → h
					output = append(output, 'h')
				} else if state.LastRune == 'ނ' && IsLabial(nextRune) {
					// Nasalization - ނ + sukun before labial (ބ/ޕ) becomes "m"
					output = append(output, 'm')
				} else if state.LastRune == 'އ' {
					// Alifu + sukun handling
					if IsConsonant(nextRune) {
						// Alifu + sukun before consonant → geminate next consonant
						state.GeminateNext = true
					} else {
						// Alifu + sukun at word end or before non-consonant → soft 'h'
						output = append(output, 'h')
					}
				} else {
					// Normal dead consonant
					output = append(output, []rune(state.LastLatin)...)
				}
				state.HasPendingConsonant = false
				state.LastLatin = ""
			}
			continue
		}


		// Vowels (fili)

		if IsVowel(r) {
			if state.HasPendingConsonant {
				var latin string
				if state.LastRune == 'އ' {
					// v2: Glottal stop suppression - suppress apostrophe between adjacent vowels
					if opts.SuppressGlottalStop {
						latin = "" // Suppress the glottal stop apostrophe
					} else if state.LastVowel != 0 && IsDiphthongPair(state.LastVowel, r) {
						// No glottal stop for diphthong pairs (e.g., o+i, a+u)
						latin = ""
					} else {
						latin = romanizeAlif(state.PositionInWord-1, true, false)
					}
				} else {
					latin = state.LastLatin
				}

				combined := latin + VowelMap[r]
				output = append(output, []rune(combined)...)
				state.HasPendingConsonant = false
			} else {
				// Leading vowel
				output = append(output, []rune(VowelMap[r])...)
			}
			state.LastVowel = r // Track for diphthong detection
			continue
		}


		// Consonants

		if IsConsonant(r) {
			flush()

			var latin string
			switch r {
			case 'ށ':
				latin = romanizeSha(false)
			case 'ނ':
				latin = romanizeNu(nextRune)
			case 'އ':
				latin = "" // Alifu carrier
			case 'ޏ':
				latin = "gn"
			default:
				latin = ConsonantMap[r]
			}

			// Gemination from alifu + sukun before this consonant
			if state.GeminateNext {
				latin = latin + latin // Double the consonant
				state.GeminateNext = false
			}

			state.LastRune = r
			state.LastLatin = latin
			state.HasPendingConsonant = true
			state.PositionInWord++
			continue
		}

		// Fallback
		flush()
		output = append(output, r)
	}

	// Final flush
	if state.LastRune == 'އ' && state.HasPendingConsonant {
		output = append(output, 'h')
	} else {
		flush()
	}

	return string(output)
}
