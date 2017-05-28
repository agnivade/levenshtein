levenshtein [![Build Status](https://travis-ci.org/agnivade/levenshtein.svg?branch=master)](https://travis-ci.org/agnivade/levenshtein) [![Go Report Card](https://goreportcard.com/badge/github.com/agnivade/levenshtein)](https://goreportcard.com/report/github.com/agnivade/levenshtein) [![GoDoc](https://godoc.org/github.com/agnivade/levenshtein?status.svg)](https://godoc.org/github.com/agnivade/levenshtein)
===========

[Go](http://golang.org) package to calculate the [Levenshtein Distance](http://en.wikipedia.org/wiki/Levenshtein_distance)

Install
-------

    go get github.com/agnivade/levenshtein

Example
-------

```go
package main

import (
	"fmt"
	"github.com/agnivade/levenshtein"
)

func main() {
	s1 := "kitten"
	s2 := "sitting"
	distance := levenshtein.ComputeDistance(s1, s2)
	fmt.Printf("The distance between %s and %s is %d.\n", s1, s2, distance)
	// Output:
	// The distance between kitten and sitting is 3.
}

```

