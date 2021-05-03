# go-scala
A pure Go library to parse Scala SCL and KBM files to support microtunings

A reimplementation, in Go, of the Surge team's C++ [tuning-library](https://surge-synth-team.org/tuning-library/).

## Example

An example of using the API:

```go
import (
  "github.com/chinenual/go-scala" scala
)
...
s,_ := scala.ScaleFromSCLFile("./my-scale.scl")
k,_ := scala.KeyboardMappingFromKBMFile("./my-mapping.kbm")
t,_ := scala.TuningFromSCLAndKBM(s,k)
fmt.Printf("The frequency of C4 and A4 are %f and %f\n",
    t.FrequencyForMidiNote(60)
    t.FrequencyForMidiNote(69))
```
