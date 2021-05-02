package scala

import (
	"fmt"
	"gotest.tools/v3/assert"
	"math"
	"math/rand"
	"path"
	"testing"
)

// margin is delta for doing floating point comparisons
// Surge uses 1.0e-10 -- need to diagnose why my numbers are less precise
const margin = 1.0e-6

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
var testKBMs = []string{
	"empty-note61.kbm",
	"empty-note69.kbm",
	"mapping-a440-constant.kbm",
	"mapping-a442-7-to-12.kbm",
	"mapping-whitekeys-a440.kbm",
	"mapping-whitekeys-c261.kbm",
	"shuffle-a440-constant.kbm",
}

func testFile(f string) string {
	return path.Join(testData, f)
}

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

// Identity Tuning Tests - 12-intune tunes properly
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

// Identity Tuning Tests - 12-intune tunes doubles properly
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

// Identity Tuning Tests - Scaling is constant
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
	t,err = CreateTuningFromKBM(k)
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
	t,err = CreateTuningFromKBM(k)
	assert.NilError(tt, err)
	assert.Equal(tt, "", approxEqual(margin, t.FrequencyForMidiNote(69), 432.0))
	assert.Equal(tt, "", approxEqual(margin, t.FrequencyForMidiNote(60), 261.625565301 * 432.0 / 440.0))
}

// Simple Keyboard Remapping Tunes A69 - Random As Scale Consistently
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
		t, err = CreateTuningFromKBM(k)
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


// Internal Constraints between Measures -- Test All Constraints SCL only
func TestInternalConstraintsSCL(tt *testing.T) {
	for _,sclFname := range testSCLs {
		var s Scale
		var t Tuning
		var err error
		s,err = ReadSCLFile(testFile(sclFname))
		assert.NilError(tt,err)
		t,err = CreateTuningFromSCL(s)
		assert.NilError(tt,err)
		for i := 0; i < 127; i++ {
			assert.Equal(tt, t.FrequencyForMidiNote(i), t.FrequencyForMidiNoteScaledByMidi0(i) * Midi0Freq, "scl:%s",sclFname)
			assert.Equal(tt, t.FrequencyForMidiNoteScaledByMidi0(i), math.Pow(2.0, t.LogScaledFrequencyForMidiNote(i)), "scl:%s",sclFname)
		}
	}
}
// Internal Constraints between Measures -- Test All Constraints KBM only
func TestInternalConstraintsKBM(tt *testing.T) {
	for _,kbmFname := range testKBMs {
		var k KeyboardMapping
		var t Tuning
		var err error
		k,err = ReadKBMFile(testFile(kbmFname))
		assert.NilError(tt,err)
		t,err = CreateTuningFromKBM(k)
		assert.NilError(tt,err)
		for i := 0; i < 127; i++ {
			assert.Equal(tt, t.FrequencyForMidiNote(i), t.FrequencyForMidiNoteScaledByMidi0(i) * Midi0Freq, "scl:%s",kbmFname)
			assert.Equal(tt, t.FrequencyForMidiNoteScaledByMidi0(i), math.Pow(2.0, t.LogScaledFrequencyForMidiNote(i)), "scl:%s",kbmFname)
		}
	}
}
// Internal Constraints between Measures -- Test All Constraints SCL & KBM
func TestInternalConstraintsSCLAndKBM(tt *testing.T) {
	for _,sclFname := range testSCLs {
		for _, kbmFname := range testKBMs {
			var k KeyboardMapping
			var s Scale
			var t Tuning
			var err error
			//fmt.Printf("DEBUG: test %s %s\n", sclFname,kbmFname)
			s, err = ReadSCLFile(testFile(sclFname))
			assert.NilError(tt, err)
			k, err = ReadKBMFile(testFile(kbmFname))
			assert.NilError(tt, err)

			if k.OctaveDegrees > s.Count {
				// don't test this combo; trap it below as an error case
				continue;
			}

			t, err = CreateTuningFromSCLAndKBM(s,k)
			assert.NilError(tt, err)

			for i := 0; i < 127; i++ {
				assert.Equal(tt, t.FrequencyForMidiNote(i), t.FrequencyForMidiNoteScaledByMidi0(i)*Midi0Freq, "scl:%s, kbm:%s", sclFname, kbmFname)
				assert.Equal(tt, t.FrequencyForMidiNoteScaledByMidi0(i), math.Pow(2.0, t.LogScaledFrequencyForMidiNote(i)), "scl:%s, kbm:%s", sclFname, kbmFname)
			}
		}
	}
}

// Did not import the Spanish locale tests.  Revisit this.

// Internal Constraints between Measures -- Mappings bigger than Scales Throw
func TestInternalConstraintsSCLAndKBMMisatched(tt *testing.T) {
	testedAtLeastOne := false
	for _,sclFname := range testSCLs {
		for _, kbmFname := range testKBMs {
			var k KeyboardMapping
			var s Scale
			var err error
			//fmt.Printf("DEBUG: test %s %s\n", sclFname,kbmFname)
			s, err = ReadSCLFile(testFile(sclFname))
			assert.NilError(tt, err)
			k, err = ReadKBMFile(testFile(kbmFname))
			assert.NilError(tt, err)

			if k.OctaveDegrees <= s.Count {
				// don't test this combo; we only want to test the error cases
				continue;
			}
			testedAtLeastOne = true
			_, err = CreateTuningFromSCLAndKBM(s,k)
			assert.ErrorContains(tt, err, "Unable to apply mapping of size","scl:%s, kbm:%s", sclFname, kbmFname)
		}
	}
	assert.Assert(tt, testedAtLeastOne)
}

// Several Sample Scales - Non Monotonic 12 note
func TestSampleScalesNonMonotonic12Note(tt *testing.T) {
	var s Scale
	var err error
	var t Tuning
	s,err = ReadSCLFile(testFile("12-shuffled.scl"))
	assert.NilError(tt,err)
	t,err = CreateTuningFromSCL(s)
	assert.NilError(tt,err)
	assert.Equal(tt, s.Count, 12)
	assert.Equal(tt, "", approxEqual(margin, t.LogScaledFrequencyForMidiNote(60), 5.0))
	order := []int{0, 2, 1, 3, 5, 4, 6, 7, 8, 10, 9, 11, 12}
	l60 := t.LogScaledFrequencyForMidiNote(60)
	for i,oi := range order {
		li := t.LogScaledFrequencyForMidiNote(60+i)
		assert.Equal(tt,"", approxEqual(margin, li-l60, float64(oi)/12.0), "order %d",oi)
	}
}

// Several Sample Scales - 31 edo
func TestSampleScales31Edo(tt *testing.T) {
	var s Scale
	var err error
	var t Tuning
	s,err = ReadSCLFile(testFile("31edo.scl"))
	assert.NilError(tt,err)
	t,err = CreateTuningFromSCL(s)
	assert.NilError(tt,err)
	assert.Equal(tt, s.Count, 31)
	assert.Equal(tt, "", approxEqual(margin, t.LogScaledFrequencyForMidiNote(60), 5.0))
	prev := t.LogScaledFrequencyForMidiNote(60)
	for i := 1; i<31; i++ {
		curr := t.LogScaledFrequencyForMidiNote(60+i)
		assert.Equal(tt,"", approxEqual(margin, curr-prev, 1.0/31.0), "i %d",i)
		prev = curr
	}
}

// Several Sample Scales - ED3-17
func TestSampleScalesEd317(tt *testing.T) {
	var s Scale
	var err error
	var t Tuning
	s,err = ReadSCLFile(testFile("ED3-17.scl"))
	assert.NilError(tt,err)
	t,err = CreateTuningFromSCL(s)
	assert.NilError(tt,err)
	assert.Equal(tt, s.Count, 17)
	assert.Equal(tt, "", approxEqual(margin, t.LogScaledFrequencyForMidiNote(60), 5.0))
	prev := t.LogScaledFrequencyForMidiNote(60)
	for i := 1; i<40; i++ {
		curr := t.LogScaledFrequencyForMidiNote(60+i)
		assert.Equal(tt,"", approxEqual(margin, math.Pow(2.0, 17.0*(curr-prev)), 3.0), "i %d",i)
		prev = curr
	}
}

// Several Sample Scales - ED4-17
func TestSampleScalesEd417(tt *testing.T) {
	var s Scale
	var err error
	var t Tuning
	s,err = ReadSCLFile(testFile("ED4-17.scl"))
	assert.NilError(tt,err)
	t,err = CreateTuningFromSCL(s)
	assert.NilError(tt,err)
	assert.Equal(tt, s.Count, 17)
	assert.Equal(tt, "", approxEqual(1e-6, t.LogScaledFrequencyForMidiNote(60), 5.0))
	prev := t.LogScaledFrequencyForMidiNote(60)
	for i := 1; i<40; i++ {
		curr := t.LogScaledFrequencyForMidiNote(60+i)
		assert.Equal(tt,"", approxEqual(1e-6, math.Pow(2.0, 17.0*(curr-prev)), 4.0), "i %d",i)
		prev = curr
	}
}

// Several Sample Scales - 6 exact
func TestSampleScales6Exact(tt *testing.T) {
	var s Scale
	var err error
	var t Tuning
	s,err = ReadSCLFile(testFile("6-exact.scl"))
	assert.NilError(tt,err)
	t,err = CreateTuningFromSCL(s)
	assert.NilError(tt,err)
	assert.Equal(tt, s.Count, 6)
	assert.Equal(tt, "", approxEqual(1e-5, t.LogScaledFrequencyForMidiNote(60), 5.0))
	knownValues := []float64{0, 0.22239,  0.41504, 0.58496, 0.73697, 0.87447, 1.0}
	for i,v := range knownValues {
		assert.Equal(tt,"", approxEqual(1e-5, t.LogScaledFrequencyForMidiNote(60+i),
			t.LogScaledFrequencyForMidiNote(60) + v),"i %d",i)
	}
}

// Several Sample Scales - Carlos Alpha (one step scale)
func TestSampleScalesCarlosAlpha(tt *testing.T) {
	var s Scale
	var err error
	var t Tuning
	s,err = ReadSCLFile(testFile("carlos-alpha.scl"))
	assert.NilError(tt,err)
	t,err = CreateTuningFromSCL(s)
	assert.NilError(tt,err)
	assert.Equal(tt, s.Count, 1)
	assert.Equal(tt, "", approxEqual(1e-5, t.LogScaledFrequencyForMidiNote(60), 5.0))
	diff := math.Pow( 2.0, 78.0 / 1200.0 )

	for i:= 30; i<80; i++ {
		assert.Equal(tt,"", approxEqual(1e-5, t.FrequencyForMidiNoteScaledByMidi0(i) * diff,
			t.FrequencyForMidiNoteScaledByMidi0(i+1)),"i %d",i)
	}
}

// Remapping frequency with non-12-length scales - 6 exact
// Remapping frequency with non-12-length scales - 31 edo
// Remapping frequency with non-12-length scales - ED4-17
// Remapping frequency with non-12-length scales - ED3-17

// KBMs with Gaps - 12 Intune with Gap

// KBM ReOrdering - Non Monotonic KBM note

// Exceptions and Bad Files - Read Non-present files
// Exceptions and Bad Files - Bad SCL
// Exceptions and Bad Files - Bad KBM

// Built in Generators - ED2
// Built in Generators - ED3-17
// Built in Generators - ED4-17
// Built in Generators - Constraints on random EDN-M
// Built in Generators - EDMN Errors
// Built in Generators - KBM Generator

// Dos Line Endings and Blanks - SCL
// Dos Line Endings and Blanks - Properly read a file with DOS line endings
// Dos Line Endings and Blanks - KBM
// Dos Line Endings and Blanks - Blank SCL

// Tone API - Valid Tones
// Tone API - Error Tones

// Scale Position - Untuned
// Scale Position - Untuned, Mapped
// Scale Position - Tuned, Unmapped
// Scale Position - Tuned, Mapped

// Default KBM Constructor has Right Base - All Scales with Default KBM

// Different KBM period from Scale period - 31Edo with mean tone mapping
// Different KBM period from Scale period - Perfect 5th UnMapped
// Different KBM period from Scale period - Perfect 5th 07 mapping







