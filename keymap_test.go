package scala

import (
	"gotest.tools/v3/assert"
	"testing"
)

// Loading tuning files - KBM File from text
func TestKeymapFromString(t *testing.T) {
	kbm,err := ParseKBMData(`! A scale file
! with zero size
0
! spanning the keybaord
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
	assert.Equal(t, kbm.Count,0)
	assert.Equal(t, kbm.FirstMidi,0)
	assert.Equal(t, kbm.LastMidi,127)
	assert.Equal(t, kbm.MiddleNote,60)
	assert.Equal(t, kbm.TuningConstantNote,69)
	assert.Equal(t, kbm.TuningFrequency,452.0)
}


// KBM Constructor RawText - KBM
func TestReparseKBMRawText(tt *testing.T) {
	var k KeyboardMapping
	var kparse KeyboardMapping
	var err error
	k, err = standardKeyboardMapping()
	assert.NilError(tt, err)

	kparse, err = ParseKBMData(k.RawText)
	assert.NilError(tt, err)

	assert.Equal(tt, k.Count, kparse.Count)
	assert.Equal(tt, k.FirstMidi, kparse.FirstMidi)
	assert.Equal(tt, k.LastMidi, kparse.LastMidi)
	assert.Equal(tt, k.MiddleNote, kparse.MiddleNote)
	assert.Equal(tt, k.TuningConstantNote, kparse.TuningConstantNote)
	assert.Equal(tt, k.TuningFrequency, kparse.TuningFrequency)
	assert.Equal(tt, k.OctaveDegrees, kparse.OctaveDegrees)
}
