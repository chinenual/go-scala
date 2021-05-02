package scala

import (
	"gotest.tools/assert"
	"testing"
)
// Load a 12 tone standard tuning
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