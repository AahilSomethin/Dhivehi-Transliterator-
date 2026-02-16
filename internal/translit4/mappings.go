package transliterator

const thaanaBase = 0x0780
const thaanaLen = 0x07B1 - 0x0780 // 49: covers U+0780 to U+07B0

const (
	alifuIdx     = '\u0787' - thaanaBase
	shaviyaniIdx = '\u0781' - thaanaBase
	noonuIdx     = '\u0782' - thaanaBase
	raaIdx       = '\u0783' - thaanaBase
	ainuIdx      = '\u07A2' - thaanaBase
	sukunIdx     = '\u07B0' - thaanaBase
	meymuIdx     = '\u0789' - thaanaBase
	baaIdx       = '\u0784' - thaanaBase
	paviyaniIdx  = '\u0795' - thaanaBase
)

// Bitmask constants: bit N is set if index N is a valid member.
// uint64 is 64 bits, thaanaLen is 49, so all indices fit.
// For out-of-range idx (uint), shift >= 64 yields 0 — no bounds guard needed.
const (
	akuruMask       uint64 = ((1 << 28) - 1) | (((1 << 9) - 1) << 29) // bits 0-27, 29-37
	filiMask        uint64 = ((1 << 10) - 1) << 38                     // bits 38-47
	sukunOvrdMask   uint64 = (1 << 12) | (1 << 15) | (1 << 34)
	akuruNameMask   uint64 = (1 << 38) - 1 // bits 0-37 (includes 0x079C)
	filiOrAkuruMask uint64 = filiMask | akuruMask
)

var akuruValues = [thaanaLen]string{
	'\u0780' - thaanaBase: "h",
	'\u0781' - thaanaBase: "sh",
	'\u0782' - thaanaBase: "n",
	'\u0783' - thaanaBase: "r",
	'\u0784' - thaanaBase: "b",
	'\u0785' - thaanaBase: "lh",
	'\u0786' - thaanaBase: "k",
	'\u0787' - thaanaBase: "",
	'\u0788' - thaanaBase: "v",
	'\u0789' - thaanaBase: "m",
	'\u078A' - thaanaBase: "f",
	'\u078B' - thaanaBase: "dh",
	'\u078C' - thaanaBase: "th",
	'\u078D' - thaanaBase: "l",
	'\u078E' - thaanaBase: "g",
	'\u078F' - thaanaBase: "gn",
	'\u0790' - thaanaBase: "s",
	'\u0791' - thaanaBase: "d",
	'\u0792' - thaanaBase: "z",
	'\u0793' - thaanaBase: "t",
	'\u0794' - thaanaBase: "y",
	'\u0795' - thaanaBase: "p",
	'\u0796' - thaanaBase: "j",
	'\u0797' - thaanaBase: "ch",
	'\u0798' - thaanaBase: "th'",
	'\u0799' - thaanaBase: "h'",
	'\u079A' - thaanaBase: "kh",
	'\u079B' - thaanaBase: "dh'",
	// '\u079C': not in the official ruleset
	'\u079D' - thaanaBase: "sh'",
	'\u079E' - thaanaBase: "s'",
	'\u079F' - thaanaBase: "l'",
	'\u07A0' - thaanaBase: "t'",
	'\u07A1' - thaanaBase: "z'",
	'\u07A2' - thaanaBase: "'",
	'\u07A3' - thaanaBase: "gh",
	'\u07A4' - thaanaBase: "q",
	'\u07A5' - thaanaBase: "w",
}

var filiValues = [thaanaLen]string{
	'\u07A6' - thaanaBase: "a",
	'\u07A7' - thaanaBase: "aa",
	'\u07A8' - thaanaBase: "i",
	'\u07A9' - thaanaBase: "ee",
	'\u07AA' - thaanaBase: "u",
	'\u07AB' - thaanaBase: "oo",
	'\u07AC' - thaanaBase: "e",
	'\u07AD' - thaanaBase: "ey",
	'\u07AE' - thaanaBase: "o",
	'\u07AF' - thaanaBase: "oa",
}

var sukunOvrdValues = [thaanaLen]string{
	'\u078C' - thaanaBase: "iy",
	'\u078F' - thaanaBase: "", // nyaviyani sukun doesn't exist in dhivehi
	'\u07A2' - thaanaBase: "u",
}

var akuruNameValues = [thaanaLen]string{
	'\u0780' - thaanaBase: "haa",
	'\u0781' - thaanaBase: "shaviyani",
	'\u0782' - thaanaBase: "noonu",
	'\u0783' - thaanaBase: "raa",
	'\u0784' - thaanaBase: "baa",
	'\u0785' - thaanaBase: "lhaviyani",
	'\u0786' - thaanaBase: "kaafu",
	'\u0787' - thaanaBase: "alifu",
	'\u0788' - thaanaBase: "vaavu",
	'\u0789' - thaanaBase: "meemu",
	'\u078A' - thaanaBase: "faafu",
	'\u078B' - thaanaBase: "dhaalu",
	'\u078C' - thaanaBase: "thaalu",
	'\u078D' - thaanaBase: "laamu",
	'\u078E' - thaanaBase: "gaafu",
	'\u078F' - thaanaBase: "gnaviyani",
	'\u0790' - thaanaBase: "seenu",
	'\u0791' - thaanaBase: "daviyani",
	'\u0792' - thaanaBase: "zaviyani",
	'\u0793' - thaanaBase: "taviyani",
	'\u0794' - thaanaBase: "yaa",
	'\u0795' - thaanaBase: "paviyani",
	'\u0796' - thaanaBase: "javiyani",
	'\u0797' - thaanaBase: "chaviyani",
	'\u0798' - thaanaBase: "tsaa",
	'\u0799' - thaanaBase: "haa",
	'\u079A' - thaanaBase: "khaa",
	'\u079B' - thaanaBase: "zhaalu",
	'\u079C' - thaanaBase: "zaa",
	'\u079D' - thaanaBase: "sheenu",
	'\u079E' - thaanaBase: "soadhu",
	'\u079F' - thaanaBase: "dzoadhu",
	'\u07A0' - thaanaBase: "thoa",
	'\u07A1' - thaanaBase: "zoa",
	'\u07A2' - thaanaBase: "ainu",
	'\u07A3' - thaanaBase: "ghainu",
	'\u07A4' - thaanaBase: "gaafu",
	'\u07A5' - thaanaBase: "vaavu",
}
