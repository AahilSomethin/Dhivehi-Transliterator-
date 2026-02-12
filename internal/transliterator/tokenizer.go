package transliterator

// IsConsonant checks if a rune is a Thaana consonant.
func IsConsonant(r rune) bool {
	_, ok := ConsonantMap[r]
	return ok
}

// IsVowel checks if a rune is a Thaana vowel sign (fili).
func IsVowel(r rune) bool {
	_, ok := VowelMap[r]
	return ok
}

// IsSukun checks if a rune is the sukūn mark.
func IsSukun(r rune) bool {
	const Sukun rune = 'ް'
	return r == Sukun
}

// IsWhitespace handles word boundaries.
func IsWhitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t'
}

// IsLabial returns true for labial consonants (ބ, ޕ).
// Used for nasalization feature.
func IsLabial(r rune) bool {
	return r == 'ބ' || r == 'ޕ'
}

// IsDiphthongPair checks if two vowels form a diphthong (no glottal stop needed).
// Common Dhivehi diphthongs: ai, au, oi, ei
func IsDiphthongPair(prev, curr rune) bool {
	prevLatin := VowelMap[prev]
	currLatin := VowelMap[curr]

	if len(prevLatin) == 0 || len(currLatin) == 0 {
		return false
	}

	lastChar := prevLatin[len(prevLatin)-1]
	firstChar := currLatin[0]

	// Common diphthongs: a+i, a+u, o+i, e+i
	if (lastChar == 'a' || lastChar == 'o' || lastChar == 'e') &&
		(firstChar == 'i' || firstChar == 'u') {
		return true
	}

	return false
}
