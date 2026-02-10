package transliterator

import "testing"

func TestTransliteration(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ދިވެހި", "dhivehi"},
		{"ބަސް", "bas"},
		{"ހުނދު", "haṁdu"},
		{"ޮބްއ", "boh"},
	}

	for _, tt := range tests {
		result := Transliterate(tt.input)
		if result != tt.expected {
			t.Errorf("Transliterate(%q) = %q, want %q",
				tt.input, result, tt.expected)
		}
	}
}
