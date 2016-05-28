package levenshtein_test

import (
	"fmt"

	"github.com/agnivade/levenshtein"
)

func Example() {
	s1 := "kitten"
	s2 := "sitting"
	distance, err := levenshtein.ComputeDistance(s1, s2)
	if err != nil {
		// handle error
	}
	fmt.Printf("The distance between %s and %s is %d.\n", s1, s2, distance)
	// Output:
	// The distance between kitten and sitting is 3.
}
