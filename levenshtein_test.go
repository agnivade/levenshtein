package levenshtein_test

import (
	"fmt"
	"strconv"
	"testing"

	agnivade "github.com/agnivade/levenshtein"
	arbovm "github.com/arbovm/levenshtein"
	dgryski "github.com/dgryski/trifles/leven"
	psadac "github.com/psadac/levenshtein"
)

func TestSanity(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"", "hello", 5},
		{"hello", "", 5},
		{"hello", "hello", 0},
		{"ab", "aa", 1},
		{"ab", "aaa", 2},
		{"bbb", "a", 3},
		{"kitten", "sitting", 3},
		{"distance", "difference", 5},
		{"levenshtein", "frankenstein", 6},
		{"resume and cafe", "resumes and cafes", 2},
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
		{"resumé and café", "resumés and cafés", 2},
		{"resume and cafe", "resumé and café", 2},
		{"Hafþór Júlíus Björnsson", "Hafþor Julius Bjornsson", 4},
		// Only 2 characters are less in the 2nd string
		{"།་གམ་འས་པ་་མ།", "།་གམའས་པ་་མ", 2},
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
		{"levenshtein", "frankenstein", "ASCII"},
		// Testing acutes and umlauts
		{"resumé and café", "resumés and cafés", "French"},
		{"Hafþór Júlíus Björnsson", "Hafþor Julius Bjornsson", "Nordic"},
		// Only 2 characters are less in the 2nd string
		{"།་གམ་འས་པ་་མ།", "།་གམའས་པ་་མ", "Tibetan"},
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

// GenStrings generate all possible strings with a given length composed of '0' and '1'.
func GenStrings(length int) []string {
	strs := make([]string, 0, 2<<length)

	for s := 1; s <= length; s++ {
		format := "%0" + strconv.Itoa(s) + "b"

		for i := 0; i < 1<<s; i++ {
			strs = append(strs, fmt.Sprintf(format, i))
		}
	}

	return strs
}

// TestWithoutCacheInit does not work.
func TestWithoutCacheInit(t *testing.T) {
	length := 4
	strs := GenStrings(length)

	for i := 0; i < len(strs); i++ {
		for j := 0; j < len(strs); j++ {
			vRef := agnivade.ComputeDistance(strs[i], strs[j])
			vTst := psadac.ComputeDistance(strs[i], strs[j])
			if vRef != vTst {
				t.Errorf("levenshtein(\"%s\",\"%s\") : want result to be %d, got %d",
					strs[i], strs[j], vRef, vTst)
			}
		}
	}
}
