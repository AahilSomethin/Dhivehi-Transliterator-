package transliterator

// applyConsonant returns the base Latin for a consonant from the map
func applyConsonant(r rune) string {
    return ConsonantMap[r]
}

// applyVowel attaches a vowel to the previous consonant
func applyVowel(prev string, vowel rune) string {
    return prev + VowelMap[vowel]
}

// romanizeSha handles the Malé Latin standard for Shaviyani (ށ)
func romanizeSha(isDoubling bool) string {
    if isDoubling {
        return "sh-sh"
    }
    return "sh" 
}

// romanizeNu handles the phonetic shift of Noonu (ނ) to "m" before labial consonants
func romanizeNu(nextRune rune) string {
    // If Noonu is followed immediately by Baa (ބ) or Paviyani (ޕ), it sounds like 'm'
    if nextRune == 'ބ' || nextRune == 'ޕ' {
        return "m"
    }
    return "n"
}

// romanizeAlif handles the carrier/stop logic for Alifu (އ)
func romanizeAlif(position int, hasVowel bool, finalWithSukun bool) string {
    // 1. Initial carrier: silent if starting a word with a vowel
    if position == 0 && hasVowel {
        return ""
    }

    // 2. Medial glottal stop: represented by an apostrophe
    if position > 0 && hasVowel {
        return "'"
    }

    // 3. Final stop: Alifu with sukun often sounds like 'h'
    if finalWithSukun {
        return "h"
    }

    return ""
}