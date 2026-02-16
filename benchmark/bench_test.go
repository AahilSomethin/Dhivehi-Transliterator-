package benchmark

import (
	"strings"
	"testing"

	v1 "dhivehi-translit/internal/translit1"
	v2 "dhivehi-translit/internal/translit2"
	v3 "dhivehi-translit/internal/translit3"
)

// Short input: a single common word.
const shortInput = "ދިވެހި"

// Medium input: a sentence with mixed features.
const mediumInput = "ވިސްނުމެއް ނެތި ކޮށްފި ކަމަކުން އެންމެ ފަހަރަކު ދޭހުގައި ގިސްލަމުން ހިތި ކަރުނަ އޮއްސަން ޖެހި ދެޔޭ ޢުމުރަށް މުޅީން"

// Long input: UDHR Article 1 & 2 in Dhivehi (para.txt content).
const longInput = `1 ވަނަ މާއްދާ ހުރިހާ އިންސާނުންވެސް ދުނިޔެއަށް އުފަންވަނީ، މިނިވަންކަމުގައި، ހަމަހަމަ ޙައްޤުތަކަކާއެކު، ހަމަހަމަ ދަރަޖައެއްގައި ކަމޭހިތެވިގެންވާ ބައެއްގެ ގޮތުގައެވެ. ހެޔޮ ވިސްނުމާއި، ހެޔޮބުއްދީގެ ބާރު އެމީހުންނަށް ލިބިގެންވެއެވެ. އަދި އެކަކު އަނެކަކާމެދު އެމީހުން މުޢާމަލާތް ކުރަންވާނީ، އުޚުއްވަތްތެރިކަމުގެ ރޫޙެއްގައެވެ.` + "\n\n" +
	`2 ވަނަ މާއްދާ ހަމަ ކޮންމެ މީހަކަށްމެ، މިޤަރާރުގައި ބަޔާންކޮށްފައިވާ ހުރިހާ ޙައްޤުތަކަކާއި މިނިވަންކަމުގެ މިންގަނޑުތަކެއް ހޯދުމާއި، ލިބިގަތުމުގެ ޙައްޤު ލިބިގެންވެއެވެ.`

// buildRepeated creates a large input by repeating mediumInput n times, separated by spaces.
func buildRepeated(n int) string {
	parts := make([]string, n)
	for i := range parts {
		parts[i] = mediumInput
	}
	return strings.Join(parts, " ")
}

// --- V1 Benchmarks ---

func BenchmarkV1_Short(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v1.Transliterate(shortInput)
	}
}

func BenchmarkV1_Medium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v1.Transliterate(mediumInput)
	}
}

func BenchmarkV1_Long(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v1.Transliterate(longInput)
	}
}

func BenchmarkV1_Repeated100(b *testing.B) {
	input := buildRepeated(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v1.Transliterate(input)
	}
}

// --- V2 Benchmarks ---

func BenchmarkV2_Short(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v2.Transliterate(shortInput)
	}
}

func BenchmarkV2_Medium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v2.Transliterate(mediumInput)
	}
}

func BenchmarkV2_Long(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v2.Transliterate(longInput)
	}
}

func BenchmarkV2_Repeated100(b *testing.B) {
	input := buildRepeated(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v2.Transliterate(input)
	}
}

// --- V3 Benchmarks ---

func BenchmarkV3_Short(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v3.Transliterate(shortInput)
	}
}

func BenchmarkV3_Medium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v3.Transliterate(mediumInput)
	}
}

func BenchmarkV3_Long(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v3.Transliterate(longInput)
	}
}

func BenchmarkV3_Repeated100(b *testing.B) {
	input := buildRepeated(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v3.Transliterate(input)
	}
}
