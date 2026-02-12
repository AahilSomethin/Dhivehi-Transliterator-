package transliterator

import "testing"

func TestTransliteration(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        // Basic Words
        {"ދިވެހި", "dhivehi"}, // Core name
        {"ބަސް", "bas"},     // Simple consonant-vowel-consonant
        {"މާލެ", "maale"},   // Vowel doubling (Aabbaafili)
        {"އަދު", "adhu"},    // Initial Alifu (silent carrier)

        // Special Sukūn Rules (Smart Logic)
        {"ބޮށް", "boh"},      // Shaviyani + Sukun context
        {"މަސް", "mas"},     // Seenu + Sukun
        {"އިށް", "ih"},      // Alifu + Shaviyani + Sukun
        {"ދެށް", "dheh"},    // Dhaalu + Sukun context

        // Nasalization (romanizeNu Rule)
        {"ހަމްދު", "hamdhu"}, // Noonu + Baa shift to 'm'
        {"އަންބަރަ", "ambara"}, // Medial nasalization
        {"ކަން", "kan"},     // Final Noonu (remains 'n')

        // Carrier & Glottal Stop (romanizeAlif Rule)
        {"އިރު", "iru"},      // Initial Alifu (silent)
        {"ބައެއް", "ba'eh"},  // Medial Alifu (glottal stop apostrophe)
        {"ގެއް", "geh"},     // Alifu + Sukun at end (soft 'h')

        // Vowel Nuances
        {"ފޮތް", "foth"},    // Obofili
        {"ބޭނުން", "beynun"}, // Eybeyfili
        {"ކާށި", "kaashi"},  // Shaviyani in middle (remains 'sh')
        {"ގޮއި", "goi"},     // Sequential vowels

        // Dotted Thaana (Arabic Loanwords)
        {"ޝަރުޠު", "sharuthu"}, // Sheenu and Taa
        {"ޤައުމު", "qaumu"},    // Qaafu
        {"ޢާއްމު", "'aammu"},   // Ainu
    }

    for _, tt := range tests {
        result := Transliterate(tt.input)
        if result != tt.expected {
            t.Errorf("Transliterate(%q) = %q, want %q",
                tt.input, result, tt.expected)
        }
    }
}

// TestTransliterationV2 tests v2 features enabled via Options.
func TestTransliterationV2(t *testing.T) {
    t.Run("Gemination", func(t *testing.T) {
        // v2 Gemination: consonant + sukun + same consonant → doubled output
        opts := Options{Gemination: true}

        tests := []struct {
            input    string
            expected string
        }{
            // ބައްބަ = baa + a + baa + sukun + baa + a → "babba" with gemination
            {"ބައްބަ", "babba"},
            // ކައްކަ = kaa + a + kaa + sukun + kaa + a → "kakka" with gemination
            {"ކައްކަ", "kakka"},
        }

        for _, tt := range tests {
            result := TransliterateWithOptions(tt.input, opts)
            if result != tt.expected {
                t.Errorf("TransliterateWithOptions(%q, Gemination) = %q, want %q",
                    tt.input, result, tt.expected)
            }
        }
    })

    t.Run("SuppressGlottalStop", func(t *testing.T) {
        // v2 SuppressGlottalStop: suppress apostrophe between adjacent vowels
        opts := Options{SuppressGlottalStop: true}

        tests := []struct {
            input    string
            expected string
        }{
            // ބައެއް normally produces ba'eh, but with suppression → "baeh"
            {"ބައެއް", "baeh"},
        }

        for _, tt := range tests {
            result := TransliterateWithOptions(tt.input, opts)
            if result != tt.expected {
                t.Errorf("TransliterateWithOptions(%q, SuppressGlottalStop) = %q, want %q",
                    tt.input, result, tt.expected)
            }
        }
    })

    t.Run("CombinedOptions", func(t *testing.T) {
        // Test multiple v2 options together
        opts := Options{
            Gemination:          true,
            SuppressGlottalStop: true,
        }

        tests := []struct {
            input    string
            expected string
        }{
            {"ބައްބަ", "babba"}, // Gemination active
            {"ބައެއް", "baeh"},  // Glottal stop suppressed
        }

        for _, tt := range tests {
            result := TransliterateWithOptions(tt.input, opts)
            if result != tt.expected {
                t.Errorf("TransliterateWithOptions(%q, Combined) = %q, want %q",
                    tt.input, result, tt.expected)
            }
        }
    })
}

// BenchmarkTransliterate measures v1 transliteration performance.
func BenchmarkTransliterate(b *testing.B) {
    input := "ދިވެހި ބަސް މާލެ އަދު ބޮށް އަންބަރަ ބައެއް ގެއް ޝަރުޠު ޤައުމު ޢާއްމު"
    for i := 0; i < b.N; i++ {
        Transliterate(input)
    }
}

// BenchmarkTransliterateWithOptions measures v2 transliteration performance.
func BenchmarkTransliterateWithOptions(b *testing.B) {
    input := "ދިވެހި ބަސް މާލެ އަދު ބޮށް އަންބަރަ ބައެއް ގެއް ޝަރުޠު ޤައުމު ޢާއްމު"
    opts := Options{Gemination: true, SuppressGlottalStop: true}
    for i := 0; i < b.N; i++ {
        TransliterateWithOptions(input, opts)
    }
}

// (from the internal/transliterator/ directory)
// go test -bench Benchmark -benchmem