// Package levenshtein is a Go implementation to calculate Levenshtein Distance.
//
// Implementation taken from
// https://gist.github.com/andrei-m/982927#gistcomment-1931258
package levenshtein

// ComputeDistance computes the levenshtein distance between the two
// strings passed as an argument. The return value is the levenshtein distance
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

	// swap to save some memory O(min(a,b)) instead of O(a)
	if len(s1) > len(s2) {
		s1, s2 = s2, s1
	}
	lenS1 := len(s1)
	lenS2 := len(s2)

	// init the row
	x := make([]int, lenS1+1)
	for i := 0; i <= lenS1; i++ {
		x[i] = i
	}

	// fill in the rest
	for i := 1; i <= lenS2; i++ {
		prev := i
		var current int

		for j := 1; j <= lenS1; j++ {

			if s2[i-1] == s1[j-1] {
				current = x[j-1] // match
			} else {
				current = min(x[j-1]+1, prev+1, x[j]+1)
			}
			x[j-1] = prev
			prev = current
		}
		x[lenS1] = prev
	}
	return x[lenS1]
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
