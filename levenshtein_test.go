package levenshtein_test

import (
	"math/rand"
	"testing"
	"time"

	agnivade "github.com/agnivade/levenshtein"
	arbovm "github.com/arbovm/levenshtein"
	dgryski "github.com/dgryski/trifles/leven"
)

// rndSeed is the random seed used for random tests and benchmarks.
const rndSeed = 42

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
			a:    "a very long string that is meant to exceed size of the row",
			b:    "another very long string that is meant to exceed size",
			want: 17,
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

// TestRndInputs tests random inputs, of random lengths, with random changes
// return values are compared to arbovm and dgryski Levenshtein implementations.
func TestRndInputs(t *testing.T) {
	const (
		nbTests    = 1000 // number of tests.
		maxLen     = 100  // maximum length in runes of rune array ra.
		maxChanges = 20   // maximum number of changes from rune array ra to rune array rb.
	)

	rnd := rand.New(rand.NewSource(rndSeed))

	for i := 0; i < nbTests; i++ {
		ra := RandRunes(rnd, maxLen)
		rb := RandRunesChange(rnd, ra, maxChanges)
		a := string(ra)
		b := string(rb)

		da := agnivade.ComputeDistance(a, b)
		dar := arbovm.Distance(a, b)
		ddg := dgryski.Levenshtein(ra, rb)

		if da != dar || da != ddg {
			t.Errorf("ComputeDistance(%s,%s) returned %d, want %d (arbovm) or %d (dgryski)", a, b, da, dar, ddg)
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

// BenchmarkRandom benchmarks random inputs, of random lengths.
func BenchmarkRandom(b *testing.B) {
	const (
		nbParams   = 10000 // number of random parameters.
		maxLen     = 100   // maximum length in runes of rune array ra.
		maxChanges = 90    // maximum number of changes from rune array ra to rune array rb.
	)

	rnd := rand.New(rand.NewSource(rndSeed))

	// create an array of random inputs.
	type param struct{ a, b string }

	params := make([]param, 0, nbParams)

	for i := 0; i < nbParams; i++ {
		ra := RandRunes(rnd, maxLen)
		rb := RandRunesChange(rnd, ra, maxChanges)
		p := param{a: string(ra), b: string(rb)}
		params = append(params, p)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		p := params[n%nbParams]
		sink = agnivade.ComputeDistance(p.a, p.b)
	}
}

// Fuzzing
// ----------------------------------------------

// FuzzComputeDistance is a fuzz test function that tests multiple implementations
// of the Levenshtein distance (agnivade, arbovm, and dgryski). It generates
// random rune sequences, applies a number of changes to one sequence, and then
// tests whether the three implementations produce the same result. The test fails
// if there is any discrepancy between the results of the different algorithms.
// Additionally, it also tests if the levenshtein distance is at most the hamming distance.
func FuzzComputeDistance(f *testing.F) {
	const (
		nbSeeds    = 100 // number of seeds.
		maxLen     = 100 // maximum length in runes of rune array ra.
		maxChanges = 20  // maximum number of changes from rune array ra to rune array rb.
	)

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < nbSeeds; i++ {
		ra := RandRunes(rnd, maxLen)
		rb := RandRunesChange(rnd, ra, maxChanges)

		f.Add(string(ra), string(rb))
	}

	f.Fuzz(func(t *testing.T, a, b string) {
		da := agnivade.ComputeDistance(a, b)
		dar := arbovm.Distance(a, b)
		ddg := dgryski.Levenshtein([]rune(a), []rune(b))

		if da != dar || da != ddg {
			t.Errorf("ComputeDistance(%s,%s) returned %d, want %d (arbovm) or %d (dgryski)", a, b, da, dar, ddg)
		}

		dh := hammingDistance([]rune(a), []rune(b))
		if da > dh {
			t.Errorf("ComputeDistance(%s,%s) returned %d, want at most %d (hamming distance)", a, b, da, dh)
		}
	})
}

func hammingDistance(a, b []rune) int {
	if len(a) > len(b) {
		a, b = b, a
	}
	d := len(b) - len(a)
	for i, r := range a {
		if r != b[i] {
			d++
		}
	}
	return d
}

// Random runes generation functions
// ----------------------------------

// RandRunes generates a random array of runes of maxLen length.
func RandRunes(rnd *rand.Rand, maxLen int) []rune {
	nbRunes := rnd.Intn(maxLen)
	runes := make([]rune, nbRunes)

	for i := 0; i < nbRunes; i++ {
		runes[i] = RandRune(rnd)
	}
	return runes
}

// RandRunesChange randomly makes maxChanges to rune array.
// A change consists of random insert or update of a random rune.
func RandRunesChange(rnd *rand.Rand, runes []rune, maxChanges int) []rune {
	if len(runes) == 0 || maxChanges == 0 {
		return runes
	}

	for i := 0; i < maxChanges; i++ {
		pos := rnd.Intn(len(runes))
		r := RandRune(rnd)

		if rnd.Intn(2) == 0 {
			// Insert
			runes = append(runes[:pos], append([]rune{r}, runes[pos:]...)...)
		} else {
			// Update
			runes[pos] = r
		}
	}
	return runes
}

// RandRune generates a random rune from a random set of runes of different byte sizes.
func RandRune(rnd *rand.Rand) rune {
	var s, e int

	switch rnd.Intn(4) {
	case 0:
		s, e = 0x0000, 0x007f // ASCII (1 byte)
	case 1:
		s, e = 0x0400, 0x04FF // Cyrillic (2 bytes)
	case 2:
		s, e = 0x0F00, 0x0FFF // Tibetan (3 bytes)
	case 3:
		s, e = 0x1F600, 0x1F64F // Emoticons (4 bytes)
	}
	return rune(s + rnd.Intn(e-s))
}
