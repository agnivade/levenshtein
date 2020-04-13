package fuzz

import (
	"strings"

	"github.com/agnivade/levenshtein"
)

func Fuzz(data []byte) int {
	str := string(data)
	if len(str) == 0 {
		return -1
	}
	parts := strings.Split(str, "\n")
	if len(parts) != 2 {
		return -1
	}
	s1 := parts[0]
	s2 := parts[1]
	res := levenshtein.ComputeDistance(s1, s2)
	// definitely an error.
	if res < 0 || res > len(s1) || res > len(s2) {
		return 0
	}
	return 1
}
