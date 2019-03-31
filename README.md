levenshtein [![Build Status](https://travis-ci.org/agnivade/levenshtein.svg?branch=master)](https://travis-ci.org/agnivade/levenshtein) [![Go Report Card](https://goreportcard.com/badge/github.com/agnivade/levenshtein)](https://goreportcard.com/report/github.com/agnivade/levenshtein) [![GoDoc](https://godoc.org/github.com/agnivade/levenshtein?status.svg)](https://godoc.org/github.com/agnivade/levenshtein)
===========

[Go](http://golang.org) package to calculate the [Levenshtein Distance](http://en.wikipedia.org/wiki/Levenshtein_distance)

The library is fully capable of working with non-ascii strings. But the strings are not normalized. That is left as a user-dependant use case. Please normalize the strings before passing it to the library if you have such a requirement.
- https://blog.golang.org/normalization

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

Benchmarks
----------

```
name              time/op
Simple/ASCII-4     365ns ± 1%
Simple/French-4    680ns ± 2%
Simple/Nordic-4   1.33µs ± 2%
Simple/Tibetan-4  1.15µs ± 2%

name              alloc/op
Simple/ASCII-4     96.0B ± 0%
Simple/French-4     128B ± 0%
Simple/Nordic-4     192B ± 0%
Simple/Tibetan-4    144B ± 0%

name              allocs/op
Simple/ASCII-4      1.00 ± 0%
Simple/French-4     1.00 ± 0%
Simple/Nordic-4     1.00 ± 0%
Simple/Tibetan-4    1.00 ± 0%
```
