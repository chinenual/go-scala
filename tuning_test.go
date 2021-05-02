package scala

import (
	"fmt"
	"gotest.tools/v3/assert"
	"math"
	"math/rand"
	"testing"
)

// margin is delta for doing floating point comparisons
// Surge uses 1.0e-10 -- need to diagnose why my numbers are less precise
const margin = 1.0e-6

// HACK:
// returns "" if equal, else a useful error message. intended to be called from assert.Equals("", approxEqual(...))
// this allows go test to report the actual line of the test failure, but still report the diff and not just the two
// values
func approxEqual(margin float64, v1 float64, v2 float64) (result string) {
	diff := math.Abs(v1 - v2)
	if diff > margin {
		result = fmt.Sprintf("ApproxEqual: %v and %v are not within %v (diff == %v)\n", v1,v2,margin,diff)
	}
	return
}

// Normal Go convention would name the *testing.T argument 't'.  But in order to keep the tests syntactically similar
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

// Scaling is constant
func TestScalingIsConstant(tt *testing.T) {
	var s Scale
	var t Tuning
	var err error
	s,err = ReadSCLFile(testFile("12-intune.scl"))
	assert.NilError(tt, err)
	t,err = CreateTuningFromSCL(s)
	assert.NilError(tt, err)

	f60 := t.FrequencyForMidiNote(60)
	fs60 := t.FrequencyForMidiNoteScaledByMidi0(60)
	for i := -200; i < 200; i++ {
		f := t.FrequencyForMidiNote(i)
		fs := t.FrequencyForMidiNoteScaledByMidi0(i)
		assert.Equal(tt, f/fs, f60/fs60)
	}
}

// Simple Keyboard Remapping Tunes A69 - A440
func TestKeyboardRemappingA69A440(tt *testing.T) {
	var k KeyboardMapping
	var t Tuning
	var err error
	k, err = tuneA69To(440.0)
	assert.NilError(tt, err)
	t,err = CreateTuningFromKBD(k)
	assert.NilError(tt, err)
	assert.Equal(tt, "", approxEqual(margin, t.FrequencyForMidiNote(69), 440.0))
	assert.Equal(tt, "", approxEqual(margin, t.FrequencyForMidiNote(60), 261.625565301))
}

// Simple Keyboard Remapping Tunes A69 - A432
func TestKeyboardRemappingA69A432(tt *testing.T) {
	var k KeyboardMapping
	var t Tuning
	var err error
	k, err = tuneA69To(432.0)
	assert.NilError(tt, err)
	t,err = CreateTuningFromKBD(k)
	assert.NilError(tt, err)
	assert.Equal(tt, "", approxEqual(margin, t.FrequencyForMidiNote(69), 432.0))
	assert.Equal(tt, "", approxEqual(margin, t.FrequencyForMidiNote(60), 261.625565301 * 432.0 / 440.0))
}

// Random As Scale Consistently
func TestRandomAsScaleConsistently(tt *testing.T) {
	var ut Tuning
	var err error
	ut, err = CreateStandardTuning()
	assert.NilError(tt, err)

	for i := 0; i < 100; i++ {
		fr := 400.0 + 80.0 * float64(rand.Int31() / math.MaxInt32)
		var k KeyboardMapping
		var t Tuning
		k, err = tuneA69To(fr)
		assert.NilError(tt, err)
		t, err = CreateTuningFromKBD(k)
		assert.NilError(tt, err)
		assert.Equal(tt, "", approxEqual(margin, t.FrequencyForMidiNote(69), fr), "i==%d", i)
		assert.Equal(tt, "", approxEqual(margin, t.FrequencyForMidiNote(60), 261.625565301  * fr / 440.0), "i==%d", i)

		ldiff := t.LogScaledFrequencyForMidiNote(69) - ut.LogScaledFrequencyForMidiNote(69)
		ratio := t.FrequencyForMidiNote(69) / ut.FrequencyForMidiNote(69)

		for j := -200; j < 200; j++ {
			assert.Equal(tt, "", approxEqual(margin, t.LogScaledFrequencyForMidiNote(j) - ut.LogScaledFrequencyForMidiNote(j), ldiff), "i==%d, j==%d", i, j)
			assert.Equal(tt, "", approxEqual(margin, t.FrequencyForMidiNote(j) / ut.FrequencyForMidiNote(j), ratio), "i==%d, j==%d", i, j)
		}
	}
}