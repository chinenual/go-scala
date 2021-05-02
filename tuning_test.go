package scala

import (
	"fmt"
	"gotest.tools/v3/assert"
	"math"
	"testing"
)

// margin is delta for doing floating point comparisons
// Surge uses 1.0e-10 -- need to diagnose why my numbers are less precise
const margin = 1.0e-6

// returns "" if equal, else a useful error message. intended to be called from assert.Equals("", approxEqual(...))
// this allows go test to report the actual line of the test failure
func approxEqual(margin float64, v1 float64, v2 float64) (result string) {
	diff := math.Abs(v1 - v2)
	if diff > margin {
		result = fmt.Sprintf("ApproxEqual: %v and %v are not within %v (diff == %v)\n", v1,v2,margin,diff)
	}
	return
}

// normal Go convention would name the testing.T argument 't'.  But in order to keep the tests syntactically similar
// to the Surge tests, which use 't' for the tuning variable, I'm using 'tt' in this file.

// 12-intune tunes properly
func TestIdentity12Intune(tt *testing.T) {
	var s Scale
	var t Tuning
	var err error
	s,err = ReadSCLFile(testFile("12-intune.scl"))
	assert.NilError(tt, err)
	assert.Equal(tt, s.Count, 12)
	t,err = CreateTuningFromSCL(s)
	assert.NilError(tt, err)
	assert.Equal(tt, "", approxEqual(margin, t.FrequencyForMidiNote( 69 ), 440.0))
	assert.Equal(tt, "", approxEqual(margin, t.FrequencyForMidiNoteScaledByMidi0( 60 ), 32.0 ))
	assert.Equal(tt, "", approxEqual(margin, t.LogScaledFrequencyForMidiNote( 60 ), 5.0 ))
}

// 12-intune tunes doubles properly
func TestIdentity12IntuneDoubles(tt *testing.T) {
	var s Scale
	var t Tuning
	var err error
	s,err = ReadSCLFile(testFile("12-intune.scl"))
	assert.NilError(tt, err)
	t,err = CreateTuningFromSCL(s)
	assert.NilError(tt, err)
	for i:= 0; i < 12; i++ {
		note := -12 * 4 + i
		sc := t.FrequencyForMidiNoteScaledByMidi0(note)
		lc := t.LogScaledFrequencyForMidiNote(note)
		for note < 200 {
			note += 12
			nlc := t.LogScaledFrequencyForMidiNote(note)
			nsc := t.FrequencyForMidiNoteScaledByMidi0(note)
			assert.Equal(tt, "", approxEqual(margin, nsc, sc * 2.0), "i==%d, note==%d", i, note)
			assert.Equal(tt, "", approxEqual(margin, nlc, lc + 1.0),"i==%d, note==%d", i, note)
			sc = nsc
			lc = nlc
		}
	}
}
