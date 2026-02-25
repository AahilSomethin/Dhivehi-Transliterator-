package benchmark

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	v1 "dhivehi-translit/internal/translit1"
	v2 "dhivehi-translit/internal/translit2"
	v3 "dhivehi-translit/internal/translit3"
	v4 "dhivehi-translit/internal/translit4"
)

// levenshtein returns the character-level edit distance between a and b.
func levenshtein(a, b string) int {
	ar, br := []rune(a), []rune(b)
	na, nb := len(ar), len(br)
	if na == 0 {
		return nb
	}
	if nb == 0 {
		return na
	}
	// One row of the DP table (we only need previous row).
	prev := make([]int, nb+1)
	curr := make([]int, nb+1)
	for j := 0; j <= nb; j++ {
		prev[j] = j
	}
	for i := 1; i <= na; i++ {
		curr[0] = i
		for j := 1; j <= nb; j++ {
			cost := 1
			if ar[i-1] == br[j-1] {
				cost = 0
			}
			curr[j] = min(curr[j-1]+1, min(prev[j]+1, prev[j-1]+cost))
		}
		prev, curr = curr, prev
	}
	return prev[nb]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// loadGoldenCases reads testdata/golden_cases.txt (tab-separated input, expected).
func loadGoldenCases(t *testing.T) (inputs []string, expected []string) {
	dir, _ := os.Getwd()
	goldenPath := filepath.Join(dir, "testdata", "golden_cases.txt")
	if strings.HasSuffix(dir, "benchmark") {
		goldenPath = filepath.Join(filepath.Dir(dir), "testdata", "golden_cases.txt")
	}
	f, err := os.Open(goldenPath)
	if err != nil {
		t.Skipf("golden file not found: %v (run from module root)", err)
		return nil, nil
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.Index(line, "\t")
		if idx < 0 {
			continue
		}
		inputs = append(inputs, line[:idx])
		expected = append(expected, strings.TrimSpace(line[idx+1:]))
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("read golden: %v", err)
	}
	return inputs, expected
}

// AccuracyReport is the structure written to accuracy_report.json.
type AccuracyReport struct {
	TotalCases int               `json:"total_cases"`
	Versions   map[string]Stats  `json:"versions"`
}

type Stats struct {
	ExactMatchPct   float64 `json:"exact_match_pct"`
	ExactMatches    int     `json:"exact_matches"`
	TotalEditDist   int     `json:"total_character_edit_distance"`
	AvgEditDist     float64 `json:"avg_character_edit_distance"`
}

// TestWriteAccuracyReport runs all four engines on the golden dataset and writes
// accuracy_report.json to the project root.
func TestWriteAccuracyReport(t *testing.T) {
	inputs, expected := loadGoldenCases(t)
	if len(inputs) == 0 {
		return
	}
	n := len(inputs)
	engines := []struct {
		name string
		fn   func(string) string
	}{
		{"translit1", v1.Transliterate},
		{"translit2", v2.Transliterate},
		{"translit3", v3.Transliterate},
		{"translit4", v4.Transliterate},
	}
	report := AccuracyReport{
		TotalCases: n,
		Versions:   make(map[string]Stats),
	}
	for _, e := range engines {
		exact := 0
		totalDist := 0
		for i := 0; i < n; i++ {
			out := e.fn(inputs[i])
			if out == expected[i] {
				exact++
			}
			totalDist += levenshtein(out, expected[i])
		}
		pct := 100.0 * float64(exact) / float64(n)
		avgDist := float64(totalDist) / float64(n)
		report.Versions[e.name] = Stats{
			ExactMatchPct: pct,
			ExactMatches:  exact,
			TotalEditDist: totalDist,
			AvgEditDist:   avgDist,
		}
	}
	dir, _ := os.Getwd()
	jsonPath := filepath.Join(dir, "accuracy_report.json")
	if strings.HasSuffix(dir, "benchmark") {
		jsonPath = filepath.Join(filepath.Dir(dir), "accuracy_report.json")
	}
	enc, _ := json.MarshalIndent(report, "", "  ")
	if err := os.WriteFile(jsonPath, enc, 0644); err != nil {
		t.Fatalf("write accuracy report: %v", err)
	}
	t.Logf("wrote %s", jsonPath)
}
