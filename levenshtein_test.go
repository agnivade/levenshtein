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

// TestComputeDistanceRnd tests ComputeDistance with random inputs, of random lengths, with random changes
// return values are compared to arbovm and dgryski Levenshtein implementations.
func TestComputeDistanceRnd(t *testing.T) {
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
			t.Errorf("ComputeDistance(%s,%s) returned %d, want %d (arbovm) %d (dgryski)",
				a, b, da, dar, ddg)
		}
	}
}

// Benchmarks
// ----------------------------------------------
var benchGroups = []string{"Edge", "ASCII", "French", "Nordic", "Tibetan",
	"Long lead", "Long middle", "Long trail", "Long diff"}

// testCasesForGroup returns the test cases that match the given group.
func testCasesForGroup(group string) testCaseArray {
	var tcg testCaseArray
	for _, tc := range testCases {
		if tc.group == group {
			tcg = append(tcg, tc)
		}
	}

	return tcg
}

func BenchmarkDistanceAgnivade(b *testing.B) {
	for _, bg := range benchGroups {
		tcg := testCasesForGroup(bg)
		b.Run(bg, func(b *testing.B) {
			for _, tc := range tcg {
				for n := 0; n < b.N; n++ {
					_ = agnivade.ComputeDistance(tc.a, tc.b)
				}
			}
		})
	}
}

func BenchmarkDistanceArbovm(b *testing.B) {
	for _, bg := range benchGroups {
		tcg := testCasesForGroup(bg)
		b.Run(bg, func(b *testing.B) {
			for _, tc := range tcg {
				for n := 0; n < b.N; n++ {
					_ = arbovm.Distance(tc.a, tc.b)
				}
			}
		})
	}
}

func BenchmarkDistanceDgryski(b *testing.B) {
	for _, bg := range benchGroups {
		tcg := testCasesForGroup(bg)
		b.Run(bg, func(b *testing.B) {
			for _, tc := range tcg {
				for n := 0; n < b.N; n++ {
					_ = dgryski.Levenshtein([]rune(tc.a), []rune(tc.b))
				}
			}
		})
	}
}

// Fuzzing
// ----------------------------------------------

// FuzzComputeDistance is a fuzz test function that compares current levenshtein function with other implementation.
// It generates random rune sequences for seeds,
// Additionally, it also tests if the levenshtein distance is at most the hamming distance.
func FuzzComputeDistance(f *testing.F) {
	const (
		nbSeeds    = 100 // number of seeds.
		maxLen     = 100 // maximum length in runes of rune array ra.
		maxChanges = 20  // maximum number of changes from rune array ra to rune array rb.
	)

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Add all test cases.
	for _, tc := range testCases {
		f.Add(tc.a, tc.b)
	}

	// Add random seeds.
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
			t.Errorf("ComputeDistance(%s,%s) returned %d, want %d (arbovm), %d (dgryski)", a, b, da, dar, ddg)
		}

		dh := pseudoHammingDistance([]rune(a), []rune(b))
		if da > dh {
			t.Errorf("ComputeDistance(%s,%s) returned %d, want at most %d (hamming distance)", a, b, da, dh)
		}
	})
}

// pseudoHammingDistance returns the hamming distance plus the length difference between 2 rune arrays.
func pseudoHammingDistance(a, b []rune) int {
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
