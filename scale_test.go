package scala

import (
	"gotest.tools/v3/assert"
	"path"
	"testing"
)

const testData = "testdata"
var testSCLs = []string{
	"12-intune.scl",
	"12-shuffled.scl" ,
	"31edo.scl",
	"6-exact.scl" ,
	"marvel12.scl" ,
	"zeus22.scl",
	"ED4-17.scl",
	"ED3-17.scl",
	"31edo_dos_lineends.scl",
}

func testFile(f string) string {
	return path.Join(testData, f)
}

// Load a 12 tone standard tuning
func TestLoadStandardTuning(t *testing.T) {
	scale,err := ReadSCLFile(testFile("12-intune.scl"))
	assert.NilError(t, err)
	assert.Equal(t, scale.Count, 12)
	// FIXME - write a lot more here obviously
}

// Load a 12 tone standard tuning with no description
func TestLoadStandardTuningNoDesc(t *testing.T) {
	scale,err := ReadSCLFile(testFile("12-intune-nodesc.scl"))
	assert.NilError(t, err)
	assert.Equal(t, scale.Count, 12)
	// FIXME - write a lot more here obviously
}