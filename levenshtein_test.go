package levenshtein_test

import (
	"testing"

	agnivade "github.com/agnivade/levenshtein"
	arbovm "github.com/arbovm/levenshtein"
	dgryski "github.com/dgryski/trifles/leven"
)

type testCaseArray []struct {
	group string // group of the test case.
	a, b  string // inputs.
	want  int    // expected result.
}

var testCases = testCaseArray{
	{group: "Edge", a: "", b: "", want: 0},
	{group: "Edge", a: "hello", b: "hello", want: 0},
	{group: "Edge", a: "hello 😊", b: "hello 😊", want: 0},
	{group: "Edge", a: "", b: "hello", want: 5},
	{group: "Edge", a: "", b: "hello 😊", want: 7},
	{group: "Edge", a: "hello", b: "hello world", want: 6},
	{group: "Edge", a: "hello", b: "hello world 😊", want: 8},
	{group: "Edge", a: "hello", b: "", want: 5},
	{group: "Edge", a: "hello 😊", b: "", want: 7},
	{group: "ASCII", a: "kitten", b: "sitting", want: 3},
	{group: "ASCII", a: "distance", b: "difference", want: 5},
	{group: "ASCII", a: "levenshtein", b: "frankenstein", want: 6},
	{group: "French", a: "resume and cafe", b: "résumé and café", want: 3},
	{group: "Nordic", a: "Hafþór Júlíus Björnsson", b: "Hafþor Julius Bjornsson", want: 4},
	{group: "Tibetan", a: "།་གམ་འས་པ་་མ།", b: "།་གམའས་པ་་མ", want: 2},
	{
		group: "Long lead",
		a:     "a very long string where the leading words are very different",
		b:     "another very long string where the leading words are very different",
		want:  6,
	},
	{
		group: "Long middle",
		a:     "a very long string with a word in the middle that is different",
		b:     "a very long string with some text in the middle that is different",
		want:  8,
	},
	{
		group: "Long trail",
		a:     "a very long string with some text at the end that is not the same",
		b:     "a very long string with some text at the end that is very different",
		want:  13,
	},
	{
		group: "Long diff",
		a:     "a very long string with different leading and trailing characters",
		b:     "this is a very long string with different leading and trailing characters.",
		want:  9,
	},
	{group: "Other", a: "some text", b: "😊😊😊some tex😊t😊", want: 5},
	{group: "Other", a: "so😊me text", b: "😊😊some tex😊t😊", want: 5},
	{group: "Other", a: "so😊me text", b: "😊😊some tex😊x😊t", want: 6},
}

// TestComputeDistance
func TestComputeDistance(t *testing.T) {
	for _, tc := range testCases {
		a := tc.a
		b := tc.b
		da := agnivade.ComputeDistance(a, b)
		dar := arbovm.Distance(a, b)
		ddg := dgryski.Levenshtein([]rune(a), []rune(b))

		if da != tc.want || da != dar || da != ddg {
			t.Errorf("ComputeDistance(%s,%s) returned %d, want %d,  %d (arbovm), %d (dgryski)",
				a, b, da, tc.want, dar, ddg)
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
