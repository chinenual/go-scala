![Go](https://github.com/chinenual/go-scala/workflows/Go/badge.svg)
[![GoReportCard](http://goreportcard.com/badge/github.com/chinenual/go-scala?1)](http://goreportcard.com/report/github.com/chinenual/go-scala)
[![Coverage Status](https://codecov.io/gh/chinenual/go-scala/branch/main/graphs/badge.svg?)](https://app.codecov.io/gh/chinenual/go-scala)
[![Go Reference](https://pkg.go.dev/badge/github.com/chinenual/go-scala.svg)](https://pkg.go.dev/github.com/chinenual/go-scala)

# go-scala
A pure Go library to parse Scala SCL and KBM files to support microtunings

A reimplementation, in Go, of the Surge team's C++ [tuning-library](https://surge-synth-team.org/tuning-library/).
This is mostly a copy of that library, but with some name changes and refactoring to make the library idiomatic Go.  

This version of the Go library (tagged v1.2.0) corresponds to the release_1.1.0 tag of the source C++ library
## Usage

```shell
$ go get github.com/chinenual/go-scala
```

## Example

An example of using the API:

```go
import (
  "github.com/chinenual/go-scala"
)
...
s,_ := scala.ScaleFromSCLFile("./my-scale.scl")
k,_ := scala.KeyboardMappingFromKBMFile("./my-mapping.kbm")
t,_ := scala.TuningFromSCLAndKBM(s,k)
fmt.Printf("The frequency of C4 and A4 are %v and %v\n",
    t.FrequencyForMidiNote(60)
    t.FrequencyForMidiNote(69))
```

## Building and testing the library:

```shell
$ go get -v -t -d ./...
$ go build
$ go test
```
