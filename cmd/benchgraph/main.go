// Program benchgraph reads benchmark_results.csv and accuracy_report.json
// (in the project root) and generates benchmarks_speed.png and
// benchmarks_accuracy.png. Uses only standard library.
// Run from project root after benchmark and accuracy report exist.
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
)

const (
	imgW     = 800
	imgH     = 480
	marginL  = 70
	marginR  = 30
	marginT  = 40
	marginB  = 70
	barColor = 0x3498db
	textColor = 0x2c3e50
)

func main() {
	dir, _ := os.Getwd()
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	csvPath := filepath.Join(dir, "benchmark_results.csv")
	jsonPath := filepath.Join(dir, "accuracy_report.json")

	if err := genSpeedGraph(csvPath, filepath.Join(dir, "benchmarks_speed.png")); err != nil {
		fmt.Fprintf(os.Stderr, "speed graph: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("wrote benchmarks_speed.png")
	if err := genAccuracyGraph(jsonPath, filepath.Join(dir, "benchmarks_accuracy.png")); err != nil {
		fmt.Fprintf(os.Stderr, "accuracy graph: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("wrote benchmarks_accuracy.png")
}

func genSpeedGraph(csvPath, outPath string) error {
	f, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer f.Close()
	recs, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}
	if len(recs) < 2 {
		return fmt.Errorf("not enough rows in CSV")
	}
	var versions []string
	var nsPerOp []float64
	for _, row := range recs[1:] {
		if len(row) < 2 {
			continue
		}
		versions = append(versions, row[0])
		n, _ := strconv.ParseFloat(row[1], 64)
		nsPerOp = append(nsPerOp, n)
	}
	return drawBarPNG(outPath, "Speed (ns/op)", versions, nsPerOp)
}

func genAccuracyGraph(jsonPath, outPath string) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return err
	}
	var report struct {
		Versions map[string]struct {
			ExactMatchPct float64 `json:"exact_match_pct"`
		} `json:"versions"`
	}
	if err := json.Unmarshal(data, &report); err != nil {
		return err
	}
	order := []string{"translit1", "translit2", "translit3", "translit4"}
	var versions []string
	var pcts []float64
	for _, v := range order {
		if s, ok := report.Versions[v]; ok {
			versions = append(versions, v)
			pcts = append(pcts, s.ExactMatchPct)
		}
	}
	return drawBarPNG(outPath, "Accuracy (%)", versions, pcts)
}

func drawBarPNG(outPath, yLabel string, labels []string, values []float64) error {
	if len(labels) != len(values) || len(labels) == 0 {
		return fmt.Errorf("labels and values length mismatch or empty")
	}
	plotW := imgW - marginL - marginR
	plotH := imgH - marginT - marginB

	maxVal := values[0]
	for _, v := range values[1:] {
		if v > maxVal {
			maxVal = v
		}
	}
	if maxVal <= 0 {
		maxVal = 1
	}

	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	white := color.RGBA{255, 255, 255, 255}
	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			img.Set(x, y, white)
		}
	}

	n := len(labels)
	barWidth := (plotW / n) * 7 / 10
	gap := (plotW / n) / 10
	x0 := marginL + gap
	yBase := marginT + plotH

	barC := color.RGBA{(barColor >> 16) & 0xff, (barColor >> 8) & 0xff, barColor & 0xff, 255}
	textC := color.RGBA{(textColor >> 16) & 0xff, (textColor >> 8) & 0xff, textColor & 0xff, 255}

	for i, v := range values {
		barH := int(float64(plotH) * v / maxVal)
		if barH > plotH {
			barH = plotH
		}
		x1 := x0 + barWidth
		y1 := yBase - barH
		for yy := y1; yy < yBase; yy++ {
			for xx := x0; xx < x1; xx++ {
				if xx >= 0 && xx < imgW && yy >= 0 && yy < imgH {
					img.Set(xx, yy, barC)
				}
			}
		}
		// Version label below (draw string as simple 6px-wide chars)
		drawString(img, x0, yBase+12, labels[i], textC)
		// Value above bar
		valStr := strconv.FormatFloat(v, 'f', 1, 64)
		if len(valStr) > 10 {
			valStr = valStr[:10]
		}
		drawString(img, x0, y1-4, valStr, textC)
		x0 += barWidth + gap
	}

	// Y-axis label (left side)
	drawString(img, 8, marginT+plotH/2-20, yLabel, textC)

	return savePNG(outPath, img)
}

// drawString draws a string using a minimal 5x7 pixel pattern per character.
func drawString(img *image.RGBA, x, y int, s string, c color.Color) {
	const w, h = 5, 7
	for _, r := range s {
		drawRune(img, x, y, r, c)
		x += w + 1
		if x+w >= img.Rect.Dx() {
			return
		}
	}
}

func drawRune(img *image.RGBA, x, y int, r rune, c color.Color) {
	const w, h = 5, 7
	// Minimal 5x7 bitmap for 0-9, ., %
	var bits []uint8
	switch r {
	case '0':
		bits = []uint8{0x1f, 0x11, 0x11, 0x11, 0x1f}
	case '1':
		bits = []uint8{0x04, 0x0c, 0x04, 0x04, 0x0e}
	case '2':
		bits = []uint8{0x1f, 0x01, 0x1f, 0x10, 0x1f}
	case '3':
		bits = []uint8{0x1f, 0x01, 0x0f, 0x01, 0x1f}
	case '4':
		bits = []uint8{0x11, 0x11, 0x1f, 0x01, 0x01}
	case '5':
		bits = []uint8{0x1f, 0x10, 0x1f, 0x01, 0x1f}
	case '6':
		bits = []uint8{0x1f, 0x10, 0x1f, 0x11, 0x1f}
	case '7':
		bits = []uint8{0x1f, 0x01, 0x02, 0x04, 0x04}
	case '8':
		bits = []uint8{0x1f, 0x11, 0x1f, 0x11, 0x1f}
	case '9':
		bits = []uint8{0x1f, 0x11, 0x1f, 0x01, 0x1f}
	case '.':
		bits = []uint8{0x00, 0x00, 0x00, 0x00, 0x04}
	case '%':
		bits = []uint8{0x18, 0x19, 0x02, 0x04, 0x13}
	case 't':
		bits = []uint8{0x04, 0x0e, 0x04, 0x04, 0x04}
	case 'r':
		bits = []uint8{0x00, 0x0e, 0x10, 0x10, 0x10}
	case 'a':
		bits = []uint8{0x00, 0x0f, 0x11, 0x11, 0x0f}
	case 'n':
		bits = []uint8{0x00, 0x1e, 0x11, 0x11, 0x11}
	case 's':
		bits = []uint8{0x00, 0x0f, 0x10, 0x0e, 0x01}
	case 'l':
		bits = []uint8{0x0c, 0x04, 0x04, 0x04, 0x0e}
	case 'i':
		bits = []uint8{0x04, 0x00, 0x04, 0x04, 0x04}
	case ' ':
		bits = []uint8{0x00, 0x00, 0x00, 0x00, 0x00}
	default:
		bits = []uint8{0x15, 0x0a, 0x15, 0x0a, 0x15}
	}
	for row := 0; row < h && row < len(bits); row++ {
		b := bits[row]
		for col := 0; col < w; col++ {
			if (b>>(4-uint(col)))&1 != 0 {
				img.Set(x+col, y+row, c)
			}
		}
	}
}

func savePNG(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
