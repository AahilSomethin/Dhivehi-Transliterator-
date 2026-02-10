package transliterator

type State struct {
	LastRune             rune
	LastLatin            string
	HasPendingConsonant  bool
	PositionInWord       int // 0 = start
}
