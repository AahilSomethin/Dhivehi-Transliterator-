package transliterator

var ConsonantMap = map[rune]string{
    'ހ': "h",   
    'ށ': "sh",   // Reset to standard
    'ނ': "n",   
    'ރ': "r",   
    'ބ': "b",   
    'ޅ': "lh",  
    'ކ': "k",   
    'އ': "",    // Carrier remains empty for rule-based handling
    'ވ': "v",   
    'މ': "m",   
    'ފ': "f",   
    'ދ': "dh",   // Reset to standard
    'ތ': "th",  
    'ލ': "l",   
    'ގ': "g",   
    'ޏ': "gn",  
    'ސ': "s",   
    'ޑ': "d",   
    'ޖ': "j",   
    'ޒ': "z",   
    'ޓ': "t",   
    'ޕ': "p",   
    'ޔ': "y",   
    'ޗ': "ch",  

    // Arabic-derived (Normalized to standard Latin)
    'ޘ': "th",
    'ޙ': "h",   
    'ޚ': "kh",
    'ޛ': "dh",
    'ޜ': "z",
    'ޝ': "sh",  
    'ޞ': "s",   
    'ޟ': "d",   
    'ޠ': "t",   
    'ޡ': "z",   
    'ޢ': "'",   
    'ޣ': "gh",
    'ޤ': "q",
}

var VowelMap = map[rune]string{
    'ަ': "a",   
    'ާ': "aa",  
    'ި': "i",   
    'ީ': "ee",  
    'ު': "u",   
    'ޫ': "oo",  
    'ެ': "e",   
    'ޭ': "ey",  
    'ޮ': "o",   
    'ޯ': "oa",  
    'ް': "",   // Handled by Sukun logic
}