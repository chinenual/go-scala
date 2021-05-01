# go-scala
A pure Go library to parse Scala SCL and KBM files to support microtunings

A reimplementation, in Go, of the Surge team's C++ [tuning-library](https://surge-synth-team.org/tuning-library/).

## Example

An example of using the API:

```go
import (
  github.com/chinenual/go-scala scala
  )
...
var s scala.Scale
var k scala.KeyboardMapping
var t scala.Tuning
if s,err = scala.ReadSCLFile("./my-scale.scl"); err != nil {
   fmt.Printf("ERROR! %v\n", err)
}
if k,err = scala.ReadKBMFile("./my-mapping.kbm"); err != nil {
   fmt.Printf("ERROR! %v\n", err)
}
if t,err = scala.CreateTuningFromSCLAndKBD(s,k) err != nil {
   fmt.Printf("ERROR! %v\n", err)
}
fmt.Printf(""The frequency of C4 and A4 are %f and %f\n",
    t.FrequencyForMidiNote(60)
    t.FrequencyForMidiNote(69))
```
