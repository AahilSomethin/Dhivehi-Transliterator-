package transliterator

// ConsonantMap provides BASE Latin values only.
// Contextual behavior (See Notes in DBA spec) MUST be handled in rules.go.
var ConsonantMap = map[rune]string{
    // Core Thaana letters
    'ހ': "h",   // Haa
    'ށ': "rh",  // Shaviyani (Linguistic 'rh' to distinguish from Sheenu)
    'ނ': "n",   // Noonu
    'ރ': "r",   // Raa
    'ބ': "b",   // Baa
    'ޅ': "lh",  // Lhaviyani
    'ކ': "k",   // Kaafu
    'އ': "a",   // Alifu (Vowel carrier)
    'ވ': "v",   // Vaavu
    'މ': "m",   // Meemu
    'ފ': "f",   // Faafu
    'ދ': "dh",  // Dhaalu
    'ތ': "th",  // Thaa
    'ލ': "l",   // Laamu
    'ގ': "g",   // Gaafu
    'ޏ': "gn",  // Nyaviyani (Matches Malé Latin chart)
    'ސ': "s",   // Seenu
    'ޑ': "d",   // Daviyani
    'ޖ': "j",   // Javiyani
    'ޒ': "z",   // Zaviyani
    'ޓ': "t",   // Taviyani
    'ޕ': "p",   // Paviyani
    'ޔ': "y",   // Yaa
    'ޗ': "ch",  // Chaviyani

    // Arabic-derived letters (Standard Latin compatibility)
    'ޘ': "th",
    'ޙ': "h",
    'ޚ': "kh",
    'ޛ': "dh",
    'ޜ': "z",
    'ޝ': "sh",  // Sheenu (Standard "sh")
    'ޞ': "s",   // Saad
    'ޟ': "d",   // Daad
    'ޠ': "t",   // Taa
    'ޡ': "z",   // Zaa
    'ޢ': "'",   // Ainu (Apostrophe for glottal stop)
    'ޣ': "gh",
    'ޤ': "q",
}

// VowelMap maps Thaana vowel signs (fili) to Latin vowels.
// These follow the Malé Latin standard for real-world readability.
var VowelMap = map[rune]string{
    'ަ': "a",   // Abafili
    'ާ': "aa",  // Aabbaafili
    'ި': "i",   // Ibifili
    'ީ': "ee",  // Eebeefili
    'ު': "u",   // Ubufili
    'ޫ': "oo",  // Ooboofili
    'ެ': "e",   // Ebefili
    'ޭ': "ey",  // Eybeyfili
    'ޮ': "o",   // Obofili
    'ޯ': "oa",  // Oaboafili
    'ް': "",    // Sukūn (Represents a dead consonant/stop)
}