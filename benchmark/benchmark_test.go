package benchmark_test

import (
	"testing"

	agnivade "github.com/agnivade/levenshtein"
	arbovm "github.com/arbovm/levenshtein"
	dgryski "github.com/dgryski/trifles/leven"
)

var sink int

func BenchmarkSimple(b *testing.B) {
	tests := []struct {
		a, b string
		name string
	}{
		// ASCII
		{a: "levenshtein", b: "frankenstein", name: "ASCII"},
		// Testing acutes and umlauts
		{a: "resumé and café", b: "resumés and cafés", name: "French"},
		{a: "Hafþór Júlíus Björnsson", b: "Hafþor Julius Bjornsson", name: "Nordic"},

		// Long strings
		{
			a:    "a very long string that is meant to exceed",
			b:    "another very long string that is meant to exceed",
			name: "Long lead",
		},
		{
			a:    "a very long string with a word in the middle that is different",
			b:    "a very long string with some text in the middle that is different",
			name: "Long middle",
		},
		{
			a:    "a very long string with some text at the end that is not the same",
			b:    "a very long string with some text at the end that is very different",
			name: "Long trail",
		},
		{
			a:    "+a very long string with different leading and trailing characters+",
			b:    "-a very long string with different leading and trailing characters-",
			name: "Long diff",
		},

		// Only 2 characters are less in the 2nd string
		{a: "།་གམ་འས་པ་་མ།", b: "།་གམའས་པ་་མ", name: "Tibetan"},
	}
	tmp := 0
	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				tmp = agnivade.ComputeDistance(test.a, test.b)
			}
		})
	}
	sink = tmp
}

func BenchmarkAll(b *testing.B) {
	tests := []struct {
		a, b string
		name string
	}{
		// ASCII
		{"levenshtein", "frankenstein", "ASCII"},
		// Testing acutes and umlauts
		{"resumé and café", "resumés and cafés", "French"},
		{"Hafþór Júlíus Björnsson", "Hafþor Julius Bjornsson", "Nordic"},
		// Only 2 characters are less in the 2nd string
		{"།་གམ་འས་པ་་མ།", "།་གམའས་པ་་མ", "Tibetan"},
	}
	tmp := 0
	for _, test := range tests {
		b.Run("case="+test.name, func(b *testing.B) {
			b.Run("impl=agniva", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					tmp = agnivade.ComputeDistance(test.a, test.b)
				}
			})
			b.Run("impl=arbovm", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					tmp = arbovm.Distance(test.a, test.b)
				}
			})
			b.Run("impl=dgryski", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					tmp = dgryski.Levenshtein([]rune(test.a), []rune(test.b))
				}
			})
		})
	}
	sink = tmp
}
