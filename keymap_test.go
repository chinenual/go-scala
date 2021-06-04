package scala

import (
	"gotest.tools/v3/assert"
	"math"
	"math/rand"
	"testing"
)

// Loading tuning files - KBM File from text
func TestKeymapFromString(t *testing.T) {
	kbm, err := KeyboardMappingFromKBMString(`! A scale file
! with zero size
0
! spanning the keyboard
0
127
! With C60 as constant and A as 452
60
69
452
! and an octave might as well be zero
0
`)
	assert.NilError(t, err)
	assert.Equal(t, kbm.Count, 0)
	assert.Equal(t, kbm.FirstMidi, 0)
	assert.Equal(t, kbm.LastMidi, 127)
	assert.Equal(t, kbm.MiddleNote, 60)
	assert.Equal(t, kbm.TuningConstantNote, 69)
	assert.Equal(t, kbm.TuningFrequency, 452.0)
}

// KBM Constructor RawText - KBM
func TestReparseKBMRawText(tt *testing.T) {
	var k KeyboardMapping
	var kparse KeyboardMapping
	var err error
	k, err = KeyboardMappingStandard()
	assert.NilError(tt, err)

	kparse, err = KeyboardMappingFromKBMString(k.RawText)
	assert.NilError(tt, err)

	assert.Equal(tt, k.Count, kparse.Count)
	assert.Equal(tt, k.FirstMidi, kparse.FirstMidi)
	assert.Equal(tt, k.LastMidi, kparse.LastMidi)
	assert.Equal(tt, k.MiddleNote, kparse.MiddleNote)
	assert.Equal(tt, k.TuningConstantNote, kparse.TuningConstantNote)
	assert.Equal(tt, k.TuningFrequency, kparse.TuningFrequency)
	assert.Equal(tt, k.OctaveDegrees, kparse.OctaveDegrees)
}

// Built in Generators - KBM Generator
func TestBuiltinGeneratorsKBMGenerator(tt *testing.T) {
	for i := 0; i < 100; i++ {
		n := int(rand.Uint32()%60 + 30)
		fr := 1000.0 * float64(rand.Uint32()/math.MaxUint32+50)

		var k KeyboardMapping
		var err error

		k, err = KeyboardMappingTuneNoteTo(n, fr)
		assert.NilError(tt, err)
		assert.Equal(tt, k.TuningConstantNote, n)
		assert.Equal(tt, k.TuningFrequency, fr)
		assert.Equal(tt, k.TuningPitch, k.TuningFrequency/midi0Freq)
		assert.Check(tt, len(k.RawText) > 1)
	}
}

// Scala KBMs from Issue 42 (tuning-library #42)
func TestIssue42KBMs(tt *testing.T) {
	var k KeyboardMapping
	var err error
	k, err = KeyboardMappingFromKBMFile(testFile("128.kbm"))
	assert.NilError(tt, err)
	assert.Equal(tt, k.Count, 0)

	k, err = KeyboardMappingFromKBMFile(testFile("piano.kbm"))
	assert.NilError(tt, err)
	assert.Equal(tt, k.Count, 0)
}

