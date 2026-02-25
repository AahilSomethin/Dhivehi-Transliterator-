package benchmark

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	v1 "dhivehi-translit/internal/translit1"
	v2 "dhivehi-translit/internal/translit2"
	v3 "dhivehi-translit/internal/translit3"
	v4 "dhivehi-translit/internal/translit4"
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

// mixedWordList is a set of Dhivehi words/phrases used to build the 10k dataset.
var mixedWordList = []string{
	"ދިވެހި", "ބަސް", "މާލެ", "އަދު", "ބޮށް", "އަންބަރަ", "ބައެއް", "ގެއް",
	"ޝަރުޠު", "ޤައުމު", "ޢާއްމު", "އިރު", "އަލަމާރި", "ރީތި", "ފޭރު", "ރޯނު",
	"އިސްކުރު", "އިސްތިރި", "ކޮއިމަލާ", "ޢަމަލް", "ނިޝާން", "ބައްޕަ", "މަންމަ",
	"ވިސްނުން", "ގަސް", "އަށް", "ކުށް", "ފޮތް", "އޯތް", "މުތް", "މަސްއޫލް",
}

// dataset10k returns a string of at least 10,000 mixed Dhivehi words (space-separated).
func dataset10k() string {
	const targetWords = 10000
	words := make([]string, 0, targetWords+len(mixedWordList))
	for len(words) < targetWords {
		for _, w := range mixedWordList {
			words = append(words, w)
			if len(words) >= targetWords {
				break
			}
		}
	}
	return strings.Join(words[:targetWords], " ")
}

// --- Unified 10k-word benchmarks (same dataset for all versions) ---

func BenchmarkTranslit1(b *testing.B) {
	input := dataset10k()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v1.Transliterate(input)
	}
}

func BenchmarkTranslit2(b *testing.B) {
	input := dataset10k()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v2.Transliterate(input)
	}
}

func BenchmarkTranslit3(b *testing.B) {
	input := dataset10k()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v3.Transliterate(input)
	}
}

func BenchmarkTranslit4(b *testing.B) {
	input := dataset10k()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v4.Transliterate(input)
	}
}

// TestWriteBenchmarkCSV runs the four BenchmarkTranslit* benchmarks and writes
// results to benchmark_results.csv in the project root. Run with:
//   go test ./benchmark/ -run TestWriteBenchmarkCSV -v
func TestWriteBenchmarkCSV(t *testing.T) {
	input := dataset10k()
	benches := []struct {
		name string
		fn   func(b *testing.B)
	}{
		{"translit1", func(b *testing.B) { for i := 0; i < b.N; i++ { v1.Transliterate(input) } }},
		{"translit2", func(b *testing.B) { for i := 0; i < b.N; i++ { v2.Transliterate(input) } }},
		{"translit3", func(b *testing.B) { for i := 0; i < b.N; i++ { v3.Transliterate(input) } }},
		{"translit4", func(b *testing.B) { for i := 0; i < b.N; i++ { v4.Transliterate(input) } }},
	}
	var out strings.Builder
	out.WriteString("version,ns_per_op,allocs_per_op,bytes_per_op\n")
	for _, bench := range benches {
		res := testing.Benchmark(bench.fn)
		out.WriteString(bench.name + ",")
		out.WriteString(strconv.FormatInt(res.NsPerOp(), 10) + ",")
		out.WriteString(strconv.FormatInt(int64(res.AllocsPerOp()), 10) + ",")
		out.WriteString(strconv.FormatInt(res.AllocedBytesPerOp(), 10) + "\n")
	}
	// Write to project root (parent of benchmark/)
	dir, _ := os.Getwd()
	csvPath := filepath.Join(dir, "benchmark_results.csv")
	if strings.HasSuffix(dir, "benchmark") {
		csvPath = filepath.Join(filepath.Dir(dir), "benchmark_results.csv")
	}
	if err := os.WriteFile(csvPath, []byte(out.String()), 0644); err != nil {
		t.Fatalf("write CSV: %v", err)
	}
	t.Logf("wrote %s", csvPath)
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
