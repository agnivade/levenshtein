package levenshtein_test

import (
	"testing"

	agnivade "github.com/agnivade/levenshtein"
	arbovm "github.com/arbovm/levenshtein"
	dgryski "github.com/dgryski/trifles/leven"
)

func TestSanity(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{a: "", b: "hello", want: 5},
		{a: "hello", b: "", want: 5},
		{a: "hello", b: "hello"},
		{a: "ab", b: "aa", want: 1},
		{a: "ab", b: "ba", want: 2},
		{a: "ab", b: "aaa", want: 2},
		{a: "bbb", b: "a", want: 3},
		{a: "kitten", b: "sitting", want: 3},
		{a: "distance", b: "difference", want: 5},
		{a: "levenshtein", b: "frankenstein", want: 6},
		{a: "resume and cafe", b: "resumes and cafes", want: 2},
		{
			a:    "a very long string that is meant to exceed",
			b:    "another very long string that is meant to exceed",
			want: 6,
		},
	}
	for i, d := range tests {
		n := agnivade.ComputeDistance(d.a, d.b)
		if n != d.want {
			t.Errorf("Test[%d]: ComputeDistance(%q,%q) returned %v, want %v",
				i, d.a, d.b, n, d.want)
		}
	}
}

func TestUnicode(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		// Testing acutes and umlauts
		{a: "resumé and café", b: "resumés and cafés", want: 2},
		{a: "resume and cafe", b: "resumé and café", want: 2},
		{a: "Hafþór Júlíus Björnsson", b: "Hafþor Julius Bjornsson", want: 4},
		// Only 2 characters are less in the 2nd string
		{a: "།་གམ་འས་པ་་མ།", b: "།་གམའས་པ་་མ", want: 2},
	}
	for i, d := range tests {
		n := agnivade.ComputeDistance(d.a, d.b)
		if n != d.want {
			t.Errorf("Test[%d]: ComputeDistance(%q,%q) returned %v, want %v",
				i, d.a, d.b, n, d.want)
		}
	}
}

// Benchmarks
// ----------------------------------------------
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
		{a: "levenshtein", b: "frankenstein", name: "ASCII"},
		// Testing acutes and umlauts
		{a: "resumé and café", b: "resumés and cafés", name: "French"},
		{a: "Hafþór Júlíus Björnsson", b: "Hafþor Julius Bjornsson", name: "Nordic"},
		// Only 2 characters are less in the 2nd string
		{a: "།་གམ་འས་པ་་མ།", b: "།་གམའས་པ་་མ", name: "Tibetan"},
	}
	tmp := 0
	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			b.Run("agniva", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					tmp = agnivade.ComputeDistance(test.a, test.b)
				}
			})
			b.Run("arbovm", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					tmp = arbovm.Distance(test.a, test.b)
				}
			})
			b.Run("dgryski", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					tmp = dgryski.Levenshtein([]rune(test.a), []rune(test.b))
				}
			})
		})
	}
	sink = tmp
}

// Fuzzing
// ----------------------------------------------

func FuzzComputeDistanceDifferent(f *testing.F) {
	testcases := []struct{ a, b string }{
		{"levenshtein", "frankenstein"},
		{"resumé and café", "resumés and cafés"},
		{"Hafþór Júlíus Björnsson", "Hafþor Julius Bjornsson"},
		{"།་གམ་འས་པ་་མ།", "།་གམའས་པ་་མ"},
		{`_p~𕍞`, `b잖PwN`},
		{`7ȪJR`, `6L)wӝ`},
		{`_p~𕍞`, `Y>q8օ݌`},
	}
	for _, tc := range testcases {
		f.Add(tc.a, tc.b)
	}
	f.Fuzz(func(t *testing.T, a, b string) {
		n := agnivade.ComputeDistance(a, b)
		if n < 0 {
			t.Errorf("Distance can not be negative: %d, a: %q, b: %q", n, a, b)
		}
		if n > len(a)+len(b) {
			t.Errorf("Distance can not be greater than sum of lengths of a and b: %d, a: %q, b: %q", n, a, b)
		}
	})
}

func FuzzComputeDistanceEqual(f *testing.F) {
	testcases := []string{
		"levenshtein", "frankenstein",
		"resumé and café", "resumés and cafés",
		"Hafþór Júlíus Björnsson", "Hafþor Julius Bjornsson",
		"།་གམ་འས་པ་་མ།", "།་གམའས་པ་་མ",
	}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, a string) {
		n := agnivade.ComputeDistance(a, a)
		if n != 0 {
			t.Errorf("Distance must be zero: %d, a: %q", n, a)
		}
	})
}
