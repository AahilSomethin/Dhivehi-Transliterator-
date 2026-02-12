package transliterator

type State struct {
	LastRune            rune
	LastLatin           string
	HasPendingConsonant bool
	PositionInWord      int  // 0 = start
	GeminateNext        bool // Alifu + sukun before consonant causes gemination
	LastVowel           rune // Track last vowel for diphthong detection
}
