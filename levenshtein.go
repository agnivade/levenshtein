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
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	x := make([]int, len(s2)+1)
	y := make([]int, len(s2)+1)
	for i := range x {
		x[i] = i
	}
	for i := range s1 {
		y[0] = i + 1

		for j := range s2 {
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
	return y[len(s2)]
}

// min is a slightly optimised version which calculates minimum of 3 integers.
// The normal version uses much more if-else conditions
func min(a, b, c int) int {
	m := a

	if m > b {
		m = b
	}
	if m > c {
		return c
	}
	return m
}
