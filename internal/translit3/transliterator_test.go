package transliterator

import "testing"

func TestTransliteration(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"އަލަމާރި", "alamaari"},
		{"އަންނާރު", "annaaru"},
		{"އެތެރެ", "ethere"},
		{"އެދުރު", "edhuru"},
		{"އިސްކުރު", "iskuru"},
		{"އިރު", "iru"},
		{"އޮނު", "onu"},
		{"އޮޑާ", "odaa"},
		{"އުތުރު", "uthuru"},
		{"އުރަ", "ura"},
		{"ރީތި", "reethi"},
		{"ހުތުރު", "huthuru"},
		{"ކޫރު", "kooru"},
		{"ބެރެބެދި", "berebedhi"},
		{"ފޭރު", "feyru"},
		{"ބޮކަރު", "bokaru"},
		{"ރޯނު", "roanu"},
		{"ފައި", "fai"},
		{"އަތަރު", "atharu"},
		{"މުނިއަވަސް", "muniavas"},
		{"އާރު", "aaru"},
		{"ކިއާށޭ", "kiaashey"},
		{"އިސްތިރި", "isthiri"},
		{"ކޮއިމަލާ", "koimalaa"},
		{"އީޓު", "eetu"},
		{"ރައީސް", "raees"},
		{"ޢަމަލް", "a'mal"},
		{"ޢާއިލާ", "a'ailaa"},
		{"މުޢީނު", "mue'enu"},
		{"ޢިޝްޤް", "i'sh'q"},
		{"މަސްޢޫދް", "maso'odh"},
		{"ޢުންޥާން", "u'nwaan"},
		{"ނަލަ", "nala"},
		{"ހަމަ", "hama"},
		{"ބާރު", "baaru"},
		{"ނާރު", "naaru"},
		{"ނިޝާން", "nish'aan"},
		{"ޚަލް", "kh'al"},   // خ (khaa) per Qawaaidu: Kh'
		{"ފިލި", "fili"},
		{"ބައްޔެއް", "bayyeh"},
		{"ބައްޕަ", "bappa"},
		{"މަންމަ", "mamma"},
		{"ކަނޑި", "kan'di"},
		{"އަނބު", "an'bu"},
		{"އުކުނު", "ukunu"},
		{"ފައުނު", "faunu"},
		{"އޫރު", "ooru"},
		{"މަސްއޫލް", "masool"},
		{"އެކަތަ", "ekatha"},
		{"ބައެއް", "baeh"},
		{"އޭނާ", "eynaa"},
		{"އޭދަފުށި", "eydhafushi"},
		{"އޮޅު", "olhu"},
		{"ކަޅުއޮއް", "kalhuoh"},
		{"އޯބު", "oabu"},
		{"އޯބަތް", "oabaiy"},
		{"އައިނު", "ainu"},
		{"އައިބު", "aibu"},
		{"މިއީ ޖުމްލައެކެވެ", "miee jumlaekeve"},
		{"ނައިފަރު", "naifaru"},
		{"މީހެއް", "meeheh"},
		{"ކުށް", "kuh"},
		{"އަށް", "ah"},
		{"ފެން", "fen"},
		{"ތުން", "thun"},
		{"ބަތް", "baiy"},
		{"ގާތް", "gaaiy"},
		{"ހިތް", "hiiy"},
		{"ކެތް", "keiy"},
		{"އޭތް", "eyiy"},
		{"ފޮތް", "foiy"},
		{"އޯތް", "oaiy"},
		{"މުތް", "muiy"},
		{"ގަސް", "gas"},
		{"ވިސްނުން", "visnun"},
		{"މިލްކު", "milku"},
		{"ރަމްޒު", "ramzu"},
		{"މާފަންނު", "maafannu"},
		{"ވިސްނުމެއް ނެތި ކޮށްފި ކަމަކުން އެންމެ ފަހަރަކު ދޭހުގައި ގިސްލަމުން ހިތި ކަރުނަ އޮއްސަން ޖެހި ދެޔޭ ޢުމުރަށް މުޅީން", "visnumeh nethi koffi kamakun emme faharaku dheyhugai gislamun hithi karuna ossan jehi dheyey u'murah mulheen"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			result := Transliterate(tt.input)
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGemination(t *testing.T) {
	opts := Options{Gemination: true}

	tests := []struct {
		input    string
		expected string
	}{
		{"ބައްބަ", "babba"},
		{"ކައްކަ", "kakka"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := TransliterateWithOptions(tt.input, opts)
			if result != tt.expected {
				t.Errorf("TransliterateWithOptions(%q, Gemination) = %q, want %q",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestSuppressGlottalStop(t *testing.T) {
	opts := Options{SuppressGlottalStop: true}

	tests := []struct {
		input    string
		expected string
	}{
		{"ބައެއް", "baeh"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := TransliterateWithOptions(tt.input, opts)
			if result != tt.expected {
				t.Errorf("TransliterateWithOptions(%q, SuppressGlottalStop) = %q, want %q",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizeArabic(t *testing.T) {
	opts := Options{NormalizeArabic: true}

	tests := []struct {
		input    string
		expected string
	}{
		{"ޝަރުޠު", "sharuthu"},
		{"ޤައުމު", "qaumu"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := TransliterateWithOptions(tt.input, opts)
			if result != tt.expected {
				t.Errorf("TransliterateWithOptions(%q, NormalizeArabic) = %q, want %q",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestCombinedOptions(t *testing.T) {
	opts := Options{
		Gemination:          true,
		SuppressGlottalStop: true,
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"ބައްބަ", "babba"},
		{"ބައެއް", "baeh"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := TransliterateWithOptions(tt.input, opts)
			if result != tt.expected {
				t.Errorf("TransliterateWithOptions(%q, Combined) = %q, want %q",
					tt.input, result, tt.expected)
			}
		})
	}
}

// --- Benchmarks ---

func BenchmarkTransliterate(b *testing.B) {
	input := "ދިވެހި ބަސް މާލެ އަދު ބޮށް އަންބަރަ ބައެއް ގެއް ޝަރުޠު ޤައުމު ޢާއްމު"
	for i := 0; i < b.N; i++ {
		Transliterate(input)
	}
}

func BenchmarkTransliterateWithOptions(b *testing.B) {
	input := "ދިވެހި ބަސް މާލެ އަދު ބޮށް އަންބަރަ ބައެއް ގެއް ޝަރުޠު ޤައުމު ޢާއްމު"
	opts := Options{Gemination: true, SuppressGlottalStop: true}
	for i := 0; i < b.N; i++ {
		TransliterateWithOptions(input, opts)
	}
}
