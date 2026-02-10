package transliterator

func Transliterate(input string) string {
	var output []rune
	state := State{}

	for _, r := range []rune(input) {

		// Whitespace resets state
		if IsWhitespace(r) {
			state = State{}
			output = append(output, r)
			continue
		}

		// ✅ Leading vowel (no pending consonant)
		if IsVowel(r) && !state.HasPendingConsonant {
			output = append(output, []rune(VowelMap[r])...)
			continue
		}

		// Vowel sign (fili) after consonant
		if IsVowel(r) && state.HasPendingConsonant {

			// Special rule: ހ + short vowel → ha
			if state.LastRune == 'ހ' {
				output = output[:len(output)-len([]rune(state.LastLatin))]
				output = append(output, []rune("ha")...)
			} else if state.LastRune == 'އ' {
				latin := romanizeAlif(state.PositionInWord-1, true, false)
				output = append(output, []rune(latin+VowelMap[r])...)
			} else {
				output = output[:len(output)-len([]rune(state.LastLatin))]
				combined := state.LastLatin + VowelMap[r]
				output = append(output, []rune(combined)...)
			}

			state.HasPendingConsonant = false
			continue
		}

		// Sukūn: suppress inherent vowel
		if IsSukun(r) {
			state.HasPendingConsonant = false
			continue
		}

		// Consonant handling
		if IsConsonant(r) {

			var latin string

			switch r {
			case 'ށ':
				latin = romanizeSha(false)
			case 'ނ':
				latin = romanizeNu(state.PositionInWord, false)
			case 'އ':
				latin = "" // defer Alif
			default:
				latin = ConsonantMap[r]
			}

			state.LastRune = r
			state.LastLatin = latin
			state.HasPendingConsonant = true
			state.PositionInWord++

			if latin != "" {
				output = append(output, []rune(latin)...)
			}
			continue
		}

		// Fallback
		output = append(output, r)
	}

	// ✅ Final Alif → h
	if state.LastRune == 'އ' {
		output = append(output, 'h')
	}

	return string(output)
}
