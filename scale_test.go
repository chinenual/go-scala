package scala

import (
	"gotest.tools/v3/assert"
	"testing"
)

// Loading tuning files - Load a 12 tone standard tuning
func TestLoadStandardTuning(t *testing.T) {
	scale, err := ScaleFromSCLFile(testFile("12-intune.scl"))
	assert.NilError(t, err)
	assert.Equal(t, scale.Count, 12)
	// FIXME - write a lot more here obviously
}

// Loading tuning files - Load a 12 tone standard tuning with no description
func TestLoadStandardTuningNoDesc(t *testing.T) {
	scale, err := ScaleFromSCLFile(testFile("12-intune-nodesc.scl"))
	assert.NilError(t, err)
	assert.Equal(t, scale.Count, 12)
	// FIXME - write a lot more here obviously
}
