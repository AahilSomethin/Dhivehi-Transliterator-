package main

import (
	"os"
	"testing"

	translit1 "dhivehi-translit/internal/translit1"
	translit2 "dhivehi-translit/internal/translit2"
	translit3 "dhivehi-translit/internal/translit3"
	translit4 "dhivehi-translit/internal/translit4"
)

var input string

func TestMain(m *testing.M) {
	data, err := os.ReadFile("para.txt")
	if err != nil {
		panic(err)
	}
	input = string(data)
	os.Exit(m.Run())
}

var sink string

func BenchmarkV1(b *testing.B) {
	for b.Loop() {
		sink = translit1.Transliterate(input)
	}
}

func BenchmarkV2(b *testing.B) {
	for b.Loop() {
		sink = translit2.Transliterate(input)
	}
}

func BenchmarkV3(b *testing.B) {
	for b.Loop() {
		sink = translit3.Transliterate(input)
	}
}

func BenchmarkV4(b *testing.B) {
	for b.Loop() {
		sink = translit4.Transliterate(input)
	}
}
