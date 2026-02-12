package transliterator

const Sukun = '\u07B0'
const Noonu = '\u0782'
const Ainu = '\u07A2'

var Akuru = map[rune]string{
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

	'\u0799': "h'",
	'\u079A': "kh",
	'\u079B': "dh'",
	// '\u079C': "z", not there in the official ruleset
	'\u079D': "sh'",
	'\u079E': "s'",
	'\u079F': "l'",
	'\u07A0': "t'",
	'\u07A1': "z'",
	'\u0798': "th'",
	'\u07A2': "'",
	'\u07A3': "gh",
	'\u07A4': "q",
	'\u07A5': "w",
}

var Fili = map[rune]string{
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

var SukunOverrides = map[rune]string{
	'\u078C': "iy",
	'\u078F': "", // nyaviyani sukun doesn't exist in dhivehi
	'\u07A2': "u",
}

var AkuruNames = map[rune]string{
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

	'\u0799': "haa",
	'\u079A': "khaa",
	'\u079B': "zhaalu",
	'\u079C': "zaa",
	'\u079D': "sheenu",
	'\u079E': "soadhu",
	'\u079F': "dzoadhu",
	'\u07A0': "thoa",
	'\u07A1': "zoa",
	'\u0798': "tsaa",
	'\u07A2': "ainu",
	'\u07A3': "ghainu",
	'\u07A4': "gaafu",
	'\u07A5': "vaavu",
}

var Nishaan = map[rune]rune{
	'\u061F': '?',
	'\u060C': ',',
	'\u061B': ';',
}
