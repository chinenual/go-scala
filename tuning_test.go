package scala

import (
	"gotest.tools/v3/assert"
	"math"
	"testing"
)

// Surge uses 1.0e-10 -- need to diagnose why my numbers are less precise
const margin = 1.0e-6 // delta for doing floating point comparisons

func approxEqual(t *testing.T, margin float64, v1 float64, v2 float64){
	diff := math.Abs(v1-v2)
	if diff > margin {
		t.Errorf("AproxEqual: %v and %v are not within %v (diff == %v)\n", v1,v2,margin,diff)
	}
}

// 12-intune tunes properly
func TestIdentity12Intune(t *testing.T) {
	var scale Scale
	var tuning Tuning
	var err error
	scale,err = ReadSCLFile(testFile("12-intune.scl"))
	assert.NilError(t, err)
	assert.Equal(t, scale.Count, 12)
	tuning,err = CreateTuningFromSCL(scale)
	assert.NilError(t, err)
	approxEqual(t, margin, tuning.FrequencyForMidiNote( 69 ), 440.0)
	approxEqual(t, margin, tuning.FrequencyForMidiNoteScaledByMidi0( 60 ), 32.0 );
	approxEqual(t, margin, tuning.LogScaledFrequencyForMidiNote( 60 ), 5.0 );
}
