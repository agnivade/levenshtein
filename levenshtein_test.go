package levenshtein

import "testing"

type testtuple struct {
	first            string
	second           string
	levenshteinValue int
}

var tests = []testtuple{
	{"", "hello", 5},
	{"hello", "", 5},
	{"hello", "hello", 0},
}

func TestSanity(t *testing.T) {
	for i, test := range tests {
		n, err := ComputeDistance(test.first, test.second)
		if err != nil {
			t.Errorf("Test[%d]: Error returned - %s. Value is - %d", i, err, n)
		}
		if n != test.levenshteinValue {
			t.Errorf("Test[%d]: Expected %d, got %d", i, test.levenshteinValue, n)
		}
	}
}

var tests1 = []testtuple{
	{"ab", "aa", 1},
	{"ab", "aaa", 2},
	{"bbb", "a", 3},
	{"kitten", "sitting", 3},
	{"distance", "difference", 5},
	{"levenshtein", "frankenstein", 6},
}

func TestNormal(t *testing.T) {
	for i, test := range tests1 {
		n, err := ComputeDistance(test.first, test.second)
		if err != nil {
			t.Errorf("Test[%d]: Error returned - %s. Value is - %d", i, err, n)
		}
		if n != test.levenshteinValue {
			t.Errorf("Test[%d]: Expected %d, got %d", i, test.levenshteinValue, n)
		}
	}
}
