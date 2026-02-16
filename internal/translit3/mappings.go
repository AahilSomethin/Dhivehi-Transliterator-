package transliterator

// Thaana Unicode range for array-indexed lookups.
const (
	thaanaBase rune = 0x0780
	thaanaSize      = 0x07B1 - 0x0780 // 49 slots: U+0780 through U+07B0
)

// Rune constants for rule-based handling.
const (
	Sukun     rune = '\u07B0'
	Noonu     rune = '\u0782'
	Ainu      rune = '\u07A2'
	Alifu     rune = '\u0787'
	Shaviyani rune = '\u0781'
	Raa       rune = '\u0783'
	Meemu     rune = '\u0789'
	Baa       rune = '\u0784'
	Paviyani  rune = '\u0795'
)

// Fast lookup arrays indexed by (r - thaanaBase).
var (
	cLat     [thaanaSize]string // consonant → Latin (V2-style, preserving Arabic distinctions)
	cLatNorm [thaanaSize]string // consonant → Latin (V1-style, normalized Arabic)
	cOk      [thaanaSize]bool   // is consonant?

	vLat [thaanaSize]string // vowel → Latin
	vOk  [thaanaSize]bool   // is vowel?

	skOver   [thaanaSize]string // sukun override values
	skOverOk [thaanaSize]bool   // has sukun override?

	akNames   [thaanaSize]string // standalone letter names
	akNamesOk [thaanaSize]bool   // has letter name?

	nishaanLat [256]rune // punctuation lookup (small range, indexed directly)
	nishaanOk  [256]bool
)

// V2-style consonant mappings (preserving Arabic distinctions with apostrophe).
var consonantData = map[rune]string{
	'\u0780': "h",
	'\u0781': "sh",
	'\u0782': "n",
	'\u0783': "r",
	'\u0784': "b",
	'\u0785': "lh",
	'\u0786': "k",
	'\u0787': "",
	'\u0788': "v",
	'\u0789': "m",
	'\u078A': "f",
	'\u078B': "dh",
	'\u078C': "th",
	'\u078D': "l",
	'\u078E': "g",
	'\u078F': "gn",
	'\u0790': "s",
	'\u0791': "d",
	'\u0792': "z",
	'\u0793': "t",
	'\u0794': "y",
	'\u0795': "p",
	'\u0796': "j",
	'\u0797': "ch",

	'\u0798': "th'",
	'\u0799': "h'",
	'\u079A': "kh",
	'\u079B': "dh'",
	'\u079D': "sh'",
	'\u079E': "s'",
	'\u079F': "l'",
	'\u07A0': "t'",
	'\u07A1': "z'",
	'\u07A2': "'",
	'\u07A3': "gh",
	'\u07A4': "q",
	'\u07A5': "w",
}

// V1-style normalized Arabic mappings (collapses Arabic-derived letters).
var consonantNormData = map[rune]string{
	'\u0798': "th",
	'\u0799': "h",
	'\u079A': "kh",
	'\u079B': "dh",
	'\u079D': "sh",
	'\u079E': "s",
	'\u079F': "d",
	'\u07A0': "th",
	'\u07A1': "z",
	'\u07A2': "'",
	'\u07A3': "gh",
	'\u07A4': "q",
	'\u07A5': "w",
}

var vowelData = map[rune]string{
	'\u07A6': "a",
	'\u07A7': "aa",
	'\u07A8': "i",
	'\u07A9': "ee",
	'\u07AA': "u",
	'\u07AB': "oo",
	'\u07AC': "e",
	'\u07AD': "ey",
	'\u07AE': "o",
	'\u07AF': "oa",
}

var sukunOverrideData = map[rune]string{
	'\u078C': "iy",
	'\u078F': "",
	'\u07A2': "u",
}

var akuruNameData = map[rune]string{
	'\u0780': "haa",
	'\u0781': "shaviyani",
	'\u0782': "noonu",
	'\u0783': "raa",
	'\u0784': "baa",
	'\u0785': "lhaviyani",
	'\u0786': "kaafu",
	'\u0787': "alifu",
	'\u0788': "vaavu",
	'\u0789': "meymu",
	'\u078A': "faafu",
	'\u078B': "dhaalu",
	'\u078C': "thaalu",
	'\u078D': "laamu",
	'\u078E': "gaafu",
	'\u078F': "nyaviyani",
	'\u0790': "seynu",
	'\u0791': "daviyani",
	'\u0792': "zaviyani",
	'\u0793': "taviyani",
	'\u0794': "yaa",
	'\u0795': "paviyani",
	'\u0796': "javiyani",
	'\u0797': "chaviyani",

	'\u0798': "tsaa",
	'\u0799': "haa",
	'\u079A': "khaa",
	'\u079B': "zhaalu",
	'\u079C': "zaa",
	'\u079D': "sheenu",
	'\u079E': "soadhu",
	'\u079F': "dzoadhu",
	'\u07A0': "thoa",
	'\u07A1': "zoa",
	'\u07A2': "ainu",
	'\u07A3': "ghainu",
	'\u07A4': "gaafu",
	'\u07A5': "vaavu",
}

var nishaanData = map[rune]rune{
	'\u061F': '?',
	'\u060C': ',',
	'\u061B': ';',
}

func init() {
	for r, s := range consonantData {
		if i := int(r - thaanaBase); i >= 0 && i < thaanaSize {
			cLat[i] = s
			cOk[i] = true
			cLatNorm[i] = s // default: same as V2
		}
	}
	for r, s := range consonantNormData {
		if i := int(r - thaanaBase); i >= 0 && i < thaanaSize {
			cLatNorm[i] = s // override Arabic-derived letters with normalized form
		}
	}
	for r, s := range vowelData {
		if i := int(r - thaanaBase); i >= 0 && i < thaanaSize {
			vLat[i] = s
			vOk[i] = true
		}
	}
	for r, s := range sukunOverrideData {
		if i := int(r - thaanaBase); i >= 0 && i < thaanaSize {
			skOver[i] = s
			skOverOk[i] = true
		}
	}
	for r, s := range akuruNameData {
		if i := int(r - thaanaBase); i >= 0 && i < thaanaSize {
			akNames[i] = s
			akNamesOk[i] = true
		}
	}
	for r, lat := range nishaanData {
		if int(r) < len(nishaanLat) {
			nishaanLat[r] = lat
			nishaanOk[r] = true
		}
	}
}
