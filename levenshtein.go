// Package levenshtein is a Go implementation to calculate Levenshtein Distance.
//
// Implementation taken from
// http://en.wikipedia.org/wiki/Levenshtein_distance#Iterative_with_two_matrix_rows
package levenshtein

// ComputeDistance computes the levenshtein distance between the two
// strings passed as an argument. The return value is the levenshtein distance
// and error (currently always nil) if any.
//
// Works on runes (Unicode code points) but does not normalize
// the input strings. See https://blog.golang.org/normalization
// and the golang.org/x/text/unicode/norm pacage.
func ComputeDistance(a, b string) int {
	if a == b {
		return 0
	}

	// Converting to []rune is simple but requires extra
	// storage and time which may be an issue for long strings.
	// This could be avoided by using utf8.RuneCountInString
	// (one pass through the string) and then careful use of
	// Go's string ranging (i.e. ignoring the byte offset
	// index returned and using our own rune index counter).
	s1 := []rune(a)
	s2 := []rune(b)

	lenS1 := len(s1)
	lenS2 := len(s2)

	if lenS1 == 0 {
		return lenS2
	}
	if lenS2 == 0 {
		return lenS1
	}

	x := make([]int, lenS2+1)
	y := make([]int, lenS2+1)
	for i := 0; i < lenS2+1; i++ {
		x[i] = i
	}
	for i := 0; i < lenS1; i++ {
		y[0] = i + 1

		for j := 0; j < lenS2; j++ {

			var cost int
			if s1[i] == s2[j] {
				cost = 0
			} else {
				cost = 1
			}
			y[j+1] = min(y[j]+1, x[j+1]+1, x[j]+cost)
		}
		copy(x, y)
	}
	return y[lenS2]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
}
