package transliterator

func Transliterate(input string) string {
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
				if state.LastRune == 'ށ' {
					// ށ + sukun → h
					output = append(output, 'h')
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
					latin = romanizeAlif(state.PositionInWord-1, true, false)
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
