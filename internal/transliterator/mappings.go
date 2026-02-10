package transliterator

// ConsonantMap maps Thaana consonants to their base Latin forms.
// NOTE: Special behavior (doubling, omission, context rules)
// is handled later in rules.go.
var ConsonantMap = map[rune]string{
	'ހ': "h",
	'ށ': "sh",
	'ނ': "n",
	'ރ': "r",
	'ބ': "b",
	'ޅ': "ḷ",
	'ކ': "k",
	'އ': "",
	'ވ': "v",
	'މ': "m",
	'ފ': "f",
	'ދ': "dh",
	'ތ': "t",
	'ލ': "l",
	'ގ': "g",
	'ސ': "s",
	'ޑ': "ḍ",
	'ޖ': "j",
	'ޒ': "z",
	'ޓ': "ṭ",
	'ޕ': "p",
	'ޔ': "y",
	'ޗ': "c",

	// Arabic-derived letters
	'ޘ': "th",
	'ޙ': "ḥ",
	'ޚ': "kh",
	'ޛ': "dh",
	'ޜ': "ṣ",
	'ޝ': "sh",
	'ޞ': "ṣ",
	'ޟ': "ḍ",
	'ޠ': "ṭ",
	'ޡ': "ẓ",
	'ޢ': "ʻ",
	'ޣ': "gh",
	'ޤ': "q",
}


// VowelMap maps Thaana vowel signs (fili) to Latin vowels.
// These are applied to the *previous consonant*.
var VowelMap = map[rune]string{
	'ަ': "a",
	'ާ': "ā",
	'ި': "i",
	'ީ': "ī",
	'ު': "u",
	'ޫ': "ū",
	'ެ': "e",
	'ޭ': "ē",
	'ޮ': "o",
	'ޯ': "ō",
}
