package transliterator

// applyConsonant returns the base Latin for a consonant
func applyConsonant(r rune) string {
	return ConsonantMap[r]
}

// applyVowel attaches a vowel to the previous consonant
func applyVowel(prev string, vowel rune) string {
	return prev + VowelMap[vowel]
}

func romanizeSha(isDoubling bool) string {
	if isDoubling {
		return "ḫ"
	}
	return "ś"
}

// ނ → n (baseline; nasalization rules can be refined later)
func romanizeNu(_ int, _ bool) string {
	return "n"
}

func romanizeAlif(position int, hasVowel bool, finalWithSukun bool) string {
	// (a) Initial with vowel → omit
	if position == 0 && hasVowel {
		return ""
	}

	// (b) Medial with vowel → apostrophe
	if position > 0 && hasVowel {
		return "’"
	}

	// (d) Final with sukūn → h
	if finalWithSukun {
		return "h"
	}

	return ""
}


//(Note: The above functions are examples of how to implement the special rules for certain consonants and contexts. The actual logic for determining when to apply these rules would be implemented in the main transliteration loop in transliterator.go, using the state to track the necessary context.)