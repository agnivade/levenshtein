// This is a Go implementation to calculate Levenshtein Distance.
// Implementation taken from
// http://en.wikipedia.org/wiki/Levenshtein_distance#Iterative_with_two_matrix_rows
package levenshtein

import (
  //"log"
  "math"
  )

// ComputeDistance computes the levenshtein distance between the two
// strings passed as an argument. The return value is the levenshtein distance
// and error if any
func ComputeDistance(s1, s2 string) (n int, err error) {

  if s1 == s2 {
    return 0, nil
  }
  if len(s1) == 0 {
    return len(s2), nil
  }
  if len(s2) == 0 {
    return len(s1), nil
  }

  x := make([]int, len(s2)+1)
  y := make([]int, len(s2)+1)
  for i, _ := range(x) {
    x[i] = i
  }
  for i, _ := range(s1) {
    y[0] = i+1

    for j, _ := range(s2) {
      var cost int
      if s1[i] == s2[j] {
        cost = 0
      } else {
        cost = 1
      }
      y[j+1] = int(math.Min(float64(y[j]+1),
                        math.Min(float64(x[j+1]+1), float64(x[j]+cost))))
    }
    copy(x, y)
  }
  return y[len(s2)], nil
}
