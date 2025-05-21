levenshtein ![Build Status](https://github.com/agnivade/levenshtein/actions/workflows/ci.yml/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/agnivade/levenshtein)](https://goreportcard.com/report/github.com/agnivade/levenshtein) [![PkgGoDev](https://pkg.go.dev/badge/github.com/agnivade/levenshtein)](https://pkg.go.dev/github.com/agnivade/levenshtein)
===========

[Go](http://golang.org) package to calculate the [Levenshtein Distance](http://en.wikipedia.org/wiki/Levenshtein_distance)

The library is fully capable of working with non-ascii strings. But the strings are not normalized. That is left as a user-dependant use case. Please normalize the strings before passing it to the library if you have such a requirement.
- https://blog.golang.org/normalization

#### Limitation

As a performance optimization, the library can handle strings only up to 65536 characters (runes). If you need to handle strings larger than that, please pin to version 1.0.3.

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

Benchmarks and comparisons with other libraries
-----------------------------------------------

Use `make benchAll COUNT=10` to run the benchmarks

```
goos: linux
goarch: amd64
pkg: github.com/agnivade/levenshtein
cpu: AMD Ryzen 7 7840U w/ Radeon  780M Graphics     
                 │ DistanceAgnivade │             DistanceArbovm             │            DistanceDgryski             │
                 │      sec/op      │    sec/op     vs base                  │    sec/op     vs base                  │
*/Edge-16               96.05n ± 1%   498.10n ± 1%   +418.56% (p=0.000 n=10)   511.35n ± 2%   +432.35% (p=0.000 n=10)
*/ASCII-16              163.9n ± 0%    402.5n ± 1%   +145.50% (p=0.000 n=10)    428.3n ± 2%   +161.27% (p=0.000 n=10)
*/French-16             176.8n ± 3%    291.5n ± 1%    +64.88% (p=0.000 n=10)    321.2n ± 1%    +81.70% (p=0.000 n=10)
*/Nordic-16             184.5n ± 0%    590.9n ± 1%   +220.27% (p=0.000 n=10)    613.2n ± 4%   +232.38% (p=0.000 n=10)
*/Tibetan-16            266.6n ± 0%    492.1n ± 0%    +84.53% (p=0.000 n=10)    520.8n ± 1%    +95.31% (p=0.000 n=10)
*/Long_lead-16          259.0n ± 1%   4178.5n ± 1%  +1513.63% (p=0.000 n=10)   4451.5n ± 2%  +1619.06% (p=0.000 n=10)
*/Long_middle-16        292.4n ± 1%   4022.0n ± 1%  +1275.51% (p=0.000 n=10)   4243.5n ± 0%  +1351.27% (p=0.000 n=10)
*/Long_trail-16         375.7n ± 1%   4378.0n ± 1%  +1065.29% (p=0.000 n=10)   4469.5n ± 1%  +1089.65% (p=0.000 n=10)
*/Long_diff-16          3.589µ ± 1%    5.117µ ± 1%    +42.59% (p=0.000 n=10)    5.219µ ± 2%    +45.42% (p=0.000 n=10)
geomean                 288.8n         1.229µ        +325.53%                   1.290µ        +346.70%

                 │ DistanceAgnivade │            DistanceArbovm            │           DistanceDgryski            │
                 │       B/op       │    B/op      vs base                 │    B/op      vs base                 │
*/Edge-16                0.0 ± 0%      344.0 ± 0%         ? (p=0.000 n=10)    344.0 ± 0%         ? (p=0.000 n=10)
*/ASCII-16               0.0 ± 0%      240.0 ± 0%         ? (p=0.000 n=10)    240.0 ± 0%         ? (p=0.000 n=10)
*/French-16              0.0 ± 0%      128.0 ± 0%         ? (p=0.000 n=10)    128.0 ± 0%         ? (p=0.000 n=10)
*/Nordic-16              0.0 ± 0%      192.0 ± 0%         ? (p=0.000 n=10)    192.0 ± 0%         ? (p=0.000 n=10)
*/Tibetan-16             0.0 ± 0%      160.0 ± 0%         ? (p=0.000 n=10)    160.0 ± 0%         ? (p=0.000 n=10)
*/Long_lead-16         544.0 ± 0%     1056.0 ± 0%   +94.12% (p=0.000 n=10)   1056.0 ± 0%   +94.12% (p=0.000 n=10)
*/Long_middle-16       544.0 ± 0%     1056.0 ± 0%   +94.12% (p=0.000 n=10)   1056.0 ± 0%   +94.12% (p=0.000 n=10)
*/Long_trail-16        576.0 ± 0%     1152.0 ± 0%  +100.00% (p=0.000 n=10)   1152.0 ± 0%  +100.00% (p=0.000 n=10)
*/Long_diff-16         752.0 ± 0%     1184.0 ± 0%   +57.45% (p=0.000 n=10)   1184.0 ± 0%   +57.45% (p=0.000 n=10)
geomean                           ¹    429.2       ?                          429.2       ?
¹ summaries must be >0 to compute geomean

                 │ DistanceAgnivade │            DistanceArbovm            │           DistanceDgryski            │
                 │    allocs/op     │ allocs/op   vs base                  │ allocs/op   vs base                  │
*/Edge-16              0.000 ± 0%     9.000 ± 0%        ? (p=0.000 n=10)     9.000 ± 0%        ? (p=0.000 n=10)
*/ASCII-16             0.000 ± 0%     3.000 ± 0%        ? (p=0.000 n=10)     3.000 ± 0%        ? (p=0.000 n=10)
*/French-16            0.000 ± 0%     1.000 ± 0%        ? (p=0.000 n=10)     1.000 ± 0%        ? (p=0.000 n=10)
*/Nordic-16            0.000 ± 0%     1.000 ± 0%        ? (p=0.000 n=10)     1.000 ± 0%        ? (p=0.000 n=10)
*/Tibetan-16           0.000 ± 0%     1.000 ± 0%        ? (p=0.000 n=10)     1.000 ± 0%        ? (p=0.000 n=10)
*/Long_lead-16         2.000 ± 0%     3.000 ± 0%  +50.00% (p=0.000 n=10)     3.000 ± 0%  +50.00% (p=0.000 n=10)
*/Long_middle-16       2.000 ± 0%     3.000 ± 0%  +50.00% (p=0.000 n=10)     3.000 ± 0%  +50.00% (p=0.000 n=10)
*/Long_trail-16        2.000 ± 0%     3.000 ± 0%  +50.00% (p=0.000 n=10)     3.000 ± 0%  +50.00% (p=0.000 n=10)
*/Long_diff-16         3.000 ± 0%     3.000 ± 0%        ~ (p=1.000 n=10) ¹   3.000 ± 0%        ~ (p=1.000 n=10) ¹
geomean                           ²   2.350       ?                          2.350       ?
¹ all samples are equal
² summaries must be >0 to compute geomean

```
