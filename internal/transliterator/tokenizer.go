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
