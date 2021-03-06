package scala

import (
	"fmt"
	"gotest.tools/v3/assert"
	"math"
	"math/rand"
	"path"
	"testing"
)

const testData = "testdata"

var testSCLs = []string{
	"12-intune.scl",
	"12-shuffled.scl",
	"31edo.scl",
	"6-exact.scl",
	"marvel12.scl",
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
		result = fmt.Sprintf("ApproxEqual: %v and %v are not within %v (diff == %v)\n", v1, v2, margin, diff)
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
	s, err = ScaleFromSCLFile(testFile("12-intune.scl"))
	assert.NilError(tt, err)
	assert.Equal(tt, s.Count, 12)
	t, err = TuningFromSCL(s)
	assert.NilError(tt, err)
	assert.Equal(tt, "", approxEqual(1e-6 /*orig:1e-10*/, t.FrequencyForMidiNote(69), 440.0))
	assert.Equal(tt, "", approxEqual(1e-6 /*orig:==*/, t.FrequencyForMidiNoteScaledByMidi0(60), 32.0))
	assert.Equal(tt, "", approxEqual(1e-6 /*orig:==*/, t.LogScaledFrequencyForMidiNote(60), 5.0))
}

// Identity Tuning Tests - 12-intune tunes doubles properly
func TestIdentity12IntuneDoubles(tt *testing.T) {
	var s Scale
	var t Tuning
	var err error
	s, err = ScaleFromSCLFile(testFile("12-intune.scl"))
	assert.NilError(tt, err)
	t, err = TuningFromSCL(s)
	assert.NilError(tt, err)
	for i := 0; i < 12; i++ {
		note := -12*4 + i
		sc := t.FrequencyForMidiNoteScaledByMidi0(note)
		lc := t.LogScaledFrequencyForMidiNote(note)
		for note < 200 {
			note += 12
			nlc := t.LogScaledFrequencyForMidiNote(note)
			nsc := t.FrequencyForMidiNoteScaledByMidi0(note)
			assert.Equal(tt, "", approxEqual(1e-10, nsc, sc*2.0), "i==%d, note==%d", i, note)
			assert.Equal(tt, "", approxEqual(1e-10, nlc, lc+1.0), "i==%d, note==%d", i, note)
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
	s, err = ScaleFromSCLFile(testFile("12-intune.scl"))
	assert.NilError(tt, err)
	t, err = TuningFromSCL(s)
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
	k, err = KeyboardMappingTuneA69To(440.0)
	assert.NilError(tt, err)
	t, err = TuningFromKBM(k)
	assert.NilError(tt, err)
	assert.Equal(tt, "", approxEqual(1e-10, t.FrequencyForMidiNote(69), 440.0))
	assert.Equal(tt, "", approxEqual(1e-9 /*orig:1-10*/, t.FrequencyForMidiNote(60), 261.625565301))
}

// Simple Keyboard Remapping Tunes A69 - A432
func TestKeyboardRemappingA69A432(tt *testing.T) {
	var k KeyboardMapping
	var t Tuning
	var err error
	k, err = KeyboardMappingTuneA69To(432.0)
	assert.NilError(tt, err)
	t, err = TuningFromKBM(k)
	assert.NilError(tt, err)
	assert.Equal(tt, "", approxEqual(1e-10, t.FrequencyForMidiNote(69), 432.0))
	assert.Equal(tt, "", approxEqual(1e-9 /*orig:1e-10*/, t.FrequencyForMidiNote(60), 261.625565301*432.0/440.0))
}

// Simple Keyboard Remapping Tunes A69 - Random As Scale Consistently
func TestRandomAsScaleConsistently(tt *testing.T) {
	var ut Tuning
	var err error
	ut, err = TuningEvenStandard()
	assert.NilError(tt, err)

	for i := 0; i < 100; i++ {
		fr := 400.0 + 80.0*float64(rand.Int31()/math.MaxInt32)
		var k KeyboardMapping
		var t Tuning
		k, err = KeyboardMappingTuneA69To(fr)
		assert.NilError(tt, err)
		t, err = TuningFromKBM(k)
		assert.NilError(tt, err)
		assert.Equal(tt, "", approxEqual(1e-10, t.FrequencyForMidiNote(69), fr), "i==%d", i)
		assert.Equal(tt, "", approxEqual(1e-9 /*orig:1e-10*/, t.FrequencyForMidiNote(60), 261.625565301*fr/440.0), "i==%d", i)

		ldiff := t.LogScaledFrequencyForMidiNote(69) - ut.LogScaledFrequencyForMidiNote(69)
		ratio := t.FrequencyForMidiNote(69) / ut.FrequencyForMidiNote(69)

		for j := -200; j < 200; j++ {
			assert.Equal(tt, "", approxEqual(1e-8, t.LogScaledFrequencyForMidiNote(j)-ut.LogScaledFrequencyForMidiNote(j), ldiff), "i==%d, j==%d", i, j)
			assert.Equal(tt, "", approxEqual(1e-8, t.FrequencyForMidiNote(j)/ut.FrequencyForMidiNote(j), ratio), "i==%d, j==%d", i, j)
		}
	}
}

// Internal Constraints between Measures -- Test All Constraints SCL only
func TestInternalConstraintsSCL(tt *testing.T) {
	for _, sclFname := range testSCLs {
		var s Scale
		var t Tuning
		var err error
		s, err = ScaleFromSCLFile(testFile(sclFname))
		assert.NilError(tt, err)
		t, err = TuningFromSCL(s)
		assert.NilError(tt, err)
		for i := 0; i < 127; i++ {
			assert.Equal(tt, t.FrequencyForMidiNote(i), t.FrequencyForMidiNoteScaledByMidi0(i)*midi0Freq, "scl:%s", sclFname)
			assert.Equal(tt, t.FrequencyForMidiNoteScaledByMidi0(i), math.Pow(2.0, t.LogScaledFrequencyForMidiNote(i)), "scl:%s", sclFname)
		}
	}
}

// Internal Constraints between Measures -- Test All Constraints KBM only
func TestInternalConstraintsKBM(tt *testing.T) {
	for _, kbmFname := range testKBMs {
		var k KeyboardMapping
		var t Tuning
		var err error
		k, err = KeyboardMappingFromKBMFile(testFile(kbmFname))
		assert.NilError(tt, err)
		t, err = TuningFromKBM(k)
		assert.NilError(tt, err)
		for i := 0; i < 127; i++ {
			assert.Equal(tt, t.FrequencyForMidiNote(i), t.FrequencyForMidiNoteScaledByMidi0(i)*midi0Freq, "scl:%s", kbmFname)
			assert.Equal(tt, t.FrequencyForMidiNoteScaledByMidi0(i), math.Pow(2.0, t.LogScaledFrequencyForMidiNote(i)), "scl:%s", kbmFname)
		}
	}
}

// Internal Constraints between Measures -- Test All Constraints SCL & KBM
func TestInternalConstraintsSCLAndKBM(tt *testing.T) {
	for _, sclFname := range testSCLs {
		for _, kbmFname := range testKBMs {
			var k KeyboardMapping
			var s Scale
			var t Tuning
			var err error

			s, err = ScaleFromSCLFile(testFile(sclFname))
			assert.NilError(tt, err)
			k, err = KeyboardMappingFromKBMFile(testFile(kbmFname))
			assert.NilError(tt, err)

			if k.OctaveDegrees > s.Count {
				// don't test this combo; trap it below as an error case
				continue
			}

			t, err = TuningFromSCLAndKBM(s, k)
			assert.NilError(tt, err)

			for i := 0; i < 127; i++ {
				assert.Equal(tt, t.FrequencyForMidiNote(i), t.FrequencyForMidiNoteScaledByMidi0(i)*midi0Freq, "scl:%s, kbm:%s", sclFname, kbmFname)
				assert.Equal(tt, t.FrequencyForMidiNoteScaledByMidi0(i), math.Pow(2.0, t.LogScaledFrequencyForMidiNote(i)), "scl:%s, kbm:%s", sclFname, kbmFname)
			}
		}
	}
}

// Did not import the Spanish locale tests.  Revisit this.

// Internal Constraints between Measures -- Mappings bigger than Scales Throw
func TestInternalConstraintsSCLAndKBMMisatched(tt *testing.T) {
	testedAtLeastOne := false
	for _, sclFname := range testSCLs {
		for _, kbmFname := range testKBMs {
			var k KeyboardMapping
			var s Scale
			var err error

			s, err = ScaleFromSCLFile(testFile(sclFname))
			assert.NilError(tt, err)
			k, err = KeyboardMappingFromKBMFile(testFile(kbmFname))
			assert.NilError(tt, err)

			if k.OctaveDegrees <= s.Count {
				// don't test this combo; we only want to test the error cases
				continue
			}
			testedAtLeastOne = true
			_, err = TuningFromSCLAndKBM(s, k)
			assert.ErrorContains(tt, err, "Unable to apply mapping of size", "scl:%s, kbm:%s", sclFname, kbmFname)
		}
	}
	assert.Assert(tt, testedAtLeastOne)
}

// Several Sample Scales - Non Monotonic 12 note
func TestSampleScalesNonMonotonic12Note(tt *testing.T) {
	var s Scale
	var err error
	var t Tuning
	s, err = ScaleFromSCLFile(testFile("12-shuffled.scl"))
	assert.NilError(tt, err)
	t, err = TuningFromSCL(s)
	assert.NilError(tt, err)
	assert.Equal(tt, s.Count, 12)
	assert.Equal(tt, "", approxEqual(1e-6, t.LogScaledFrequencyForMidiNote(60), 5.0))
	order := []int{0, 2, 1, 3, 5, 4, 6, 7, 8, 10, 9, 11, 12}
	l60 := t.LogScaledFrequencyForMidiNote(60)
	for i, oi := range order {
		li := t.LogScaledFrequencyForMidiNote(60 + i)
		assert.Equal(tt, "", approxEqual(1e-6, li-l60, float64(oi)/12.0), "order %d", oi)
	}
}

// Several Sample Scales - 31 edo
func TestSampleScales31Edo(tt *testing.T) {
	var s Scale
	var err error
	var t Tuning
	s, err = ScaleFromSCLFile(testFile("31edo.scl"))
	assert.NilError(tt, err)
	t, err = TuningFromSCL(s)
	assert.NilError(tt, err)
	assert.Equal(tt, s.Count, 31)
	assert.Equal(tt, "", approxEqual(1e-6, t.LogScaledFrequencyForMidiNote(60), 5.0))
	prev := t.LogScaledFrequencyForMidiNote(60)
	for i := 1; i < 31; i++ {
		curr := t.LogScaledFrequencyForMidiNote(60 + i)
		assert.Equal(tt, "", approxEqual(1e-6, curr-prev, 1.0/31.0), "i %d", i)
		prev = curr
	}
}

// Several Sample Scales - ED3-17
func TestSampleScalesEd317(tt *testing.T) {
	var s Scale
	var err error
	var t Tuning
	s, err = ScaleFromSCLFile(testFile("ED3-17.scl"))
	assert.NilError(tt, err)
	t, err = TuningFromSCL(s)
	assert.NilError(tt, err)
	assert.Equal(tt, s.Count, 17)
	assert.Equal(tt, "", approxEqual(1e-6, t.LogScaledFrequencyForMidiNote(60), 5.0))
	prev := t.LogScaledFrequencyForMidiNote(60)
	for i := 1; i < 40; i++ {
		curr := t.LogScaledFrequencyForMidiNote(60 + i)
		assert.Equal(tt, "", approxEqual(1e-6, math.Pow(2.0, 17.0*(curr-prev)), 3.0), "i %d", i)
		prev = curr
	}
}

// Several Sample Scales - ED4-17
func TestSampleScalesEd417(tt *testing.T) {
	var s Scale
	var err error
	var t Tuning
	s, err = ScaleFromSCLFile(testFile("ED4-17.scl"))
	assert.NilError(tt, err)
	t, err = TuningFromSCL(s)
	assert.NilError(tt, err)
	assert.Equal(tt, s.Count, 17)
	assert.Equal(tt, "", approxEqual(1e-6, t.LogScaledFrequencyForMidiNote(60), 5.0))
	prev := t.LogScaledFrequencyForMidiNote(60)
	for i := 1; i < 40; i++ {
		curr := t.LogScaledFrequencyForMidiNote(60 + i)
		assert.Equal(tt, "", approxEqual(1e-6, math.Pow(2.0, 17.0*(curr-prev)), 4.0), "i %d", i)
		prev = curr
	}
}

// Several Sample Scales - 6 exact
func TestSampleScales6Exact(tt *testing.T) {
	var s Scale
	var err error
	var t Tuning
	s, err = ScaleFromSCLFile(testFile("6-exact.scl"))
	assert.NilError(tt, err)
	t, err = TuningFromSCL(s)
	assert.NilError(tt, err)
	assert.Equal(tt, s.Count, 6)
	assert.Equal(tt, "", approxEqual(1e-5, t.LogScaledFrequencyForMidiNote(60), 5.0))
	knownValues := []float64{0, 0.22239, 0.41504, 0.58496, 0.73697, 0.87447, 1.0}
	for i, v := range knownValues {
		assert.Equal(tt, "", approxEqual(1e-5, t.LogScaledFrequencyForMidiNote(60+i),
			t.LogScaledFrequencyForMidiNote(60)+v), "i %d", i)
	}
}

// Several Sample Scales - Carlos Alpha (one step scale)
func TestSampleScalesCarlosAlpha(tt *testing.T) {
	var s Scale
	var err error
	var t Tuning
	s, err = ScaleFromSCLFile(testFile("carlos-alpha.scl"))
	assert.NilError(tt, err)
	t, err = TuningFromSCL(s)
	assert.NilError(tt, err)
	assert.Equal(tt, s.Count, 1)
	assert.Equal(tt, "", approxEqual(1e-5, t.LogScaledFrequencyForMidiNote(60), 5.0))
	diff := math.Pow(2.0, 78.0/1200.0)

	for i := 30; i < 80; i++ {
		assert.Equal(tt, "", approxEqual(1e-5, t.FrequencyForMidiNoteScaledByMidi0(i)*diff,
			t.FrequencyForMidiNoteScaledByMidi0(i+1)), "i %d", i)
	}
}

// Remapping frequency with non-12-length scales - 6 exact
func TestRemappingFreqWithNon12Scales6Exact(tt *testing.T) {
	var s Scale
	var err error
	var t Tuning
	s, err = ScaleFromSCLFile(testFile("6-exact.scl"))
	assert.NilError(tt, err)
	t, err = TuningFromSCL(s)
	assert.NilError(tt, err)
	for i := 0; i < 100; i++ {
		mn := int(rand.Uint32()%40 + 40)
		freq := 150.0 + 300.0*float64(rand.Uint32())/float64(math.MaxUint32)
		var k KeyboardMapping
		k, err = KeyboardMappingTuneNoteTo(mn, freq)
		assert.NilError(tt, err)
		var mapped Tuning
		mapped, err = TuningFromSCLAndKBM(s, k)
		assert.NilError(tt, err)

		assert.Equal(tt, "", approxEqual(1e-6, mapped.FrequencyForMidiNote(mn), freq), "mn:%v, freq:%v", mn, freq)

		// This scale is monotonic so test monotonicity still
		for ii := 1; ii < 127; ii++ {
			if mapped.FrequencyForMidiNote(ii) > 1.0 {
				assert.Check(tt, mapped.FrequencyForMidiNote(ii) > mapped.FrequencyForMidiNote(ii-1), "mn:%v, freq:%v, ii:%v", mn, freq, ii)
			}
		}
		n60ldiff := t.LogScaledFrequencyForMidiNote(60) - mapped.LogScaledFrequencyForMidiNote(60)
		for j := 0; j < 128; j++ {
			assert.Equal(tt, "", approxEqual(1e-6, t.LogScaledFrequencyForMidiNote(j)-mapped.LogScaledFrequencyForMidiNote(j),
				n60ldiff), "mn:%v, freq:%v, j:%v,tl:%v,mappedl:%v", mn, freq, j, t.LogScaledFrequencyForMidiNote(j), mapped.LogScaledFrequencyForMidiNote(j))
		}
	}
}

// Remapping frequency with non-12-length scales - 31 edo
func TestRemappingFreqWithNon12Scales31Edo(tt *testing.T) {
	var s Scale
	var err error
	var t Tuning
	s, err = ScaleFromSCLFile(testFile("31edo.scl"))
	assert.NilError(tt, err)
	t, err = TuningFromSCL(s)
	assert.NilError(tt, err)
	for i := 0; i < 100; i++ {
		mn := int(rand.Uint32()%20 + 50)
		freq := 150.0 + 300.0*float64(rand.Uint32())/float64(math.MaxUint32)
		var k KeyboardMapping
		k, err = KeyboardMappingTuneNoteTo(mn, freq)
		assert.NilError(tt, err)
		var mapped Tuning
		mapped, err = TuningFromSCLAndKBM(s, k)
		assert.NilError(tt, err)

		assert.Equal(tt, "", approxEqual(1e-6, mapped.FrequencyForMidiNote(mn), freq), "mn:%v, freq:%v", mn, freq)

		// This scale is monotonic so test monotonicity still
		for ii := 1; ii < 127; ii++ {
			if mapped.FrequencyForMidiNote(ii) > 1.0 {
				assert.Check(tt, mapped.FrequencyForMidiNote(ii) > mapped.FrequencyForMidiNote(ii-1), "mn:%v, freq:%v, ii:%v", mn, freq, ii)
			}
		}
		n60ldiff := t.LogScaledFrequencyForMidiNote(60) - mapped.LogScaledFrequencyForMidiNote(60)
		for j := 0; j < 128; j++ {
			assert.Equal(tt, "", approxEqual(1e-6, t.LogScaledFrequencyForMidiNote(j)-mapped.LogScaledFrequencyForMidiNote(j),
				n60ldiff), "mn:%v, freq:%v, j:%v", mn, freq, j)
		}
	}
}

// Remapping frequency with non-12-length scales - ED4-17
func TestRemappingFreqWithNon12ScalesED417(tt *testing.T) {
	var s Scale
	var err error
	var t Tuning
	s, err = ScaleFromSCLFile(testFile("ED4-17.scl"))
	assert.NilError(tt, err)
	t, err = TuningFromSCL(s)
	assert.NilError(tt, err)
	for i := 0; i < 100; i++ {
		mn := int(rand.Uint32()%40 + 40)
		freq := 150.0 + 300.0*float64(rand.Uint32())/float64(math.MaxUint32)
		var k KeyboardMapping
		k, err = KeyboardMappingTuneNoteTo(mn, freq)
		assert.NilError(tt, err)
		var mapped Tuning
		mapped, err = TuningFromSCLAndKBM(s, k)
		assert.NilError(tt, err)

		assert.Equal(tt, "", approxEqual(1e-6, mapped.FrequencyForMidiNote(mn), freq), "mn:%v, freq:%v", mn, freq)

		// This scale is monotonic so test monotonicity still
		for ii := 1; ii < 127; ii++ {
			if mapped.FrequencyForMidiNote(ii) > 1.0 {
				assert.Check(tt, mapped.FrequencyForMidiNote(ii) > mapped.FrequencyForMidiNote(ii-1), "mn:%v, freq:%v, ii:%v", mn, freq, ii)
			}
		}
		n60ldiff := t.LogScaledFrequencyForMidiNote(60) - mapped.LogScaledFrequencyForMidiNote(60)
		for j := 0; j < 128; j++ {
			assert.Equal(tt, "", approxEqual(1e-6, t.LogScaledFrequencyForMidiNote(j)-mapped.LogScaledFrequencyForMidiNote(j),
				n60ldiff), "mn:%v, freq:%v, j:%v", mn, freq, j)
		}
	}
}

// Remapping frequency with non-12-length scales - ED3-17
func TestRemappingFreqWithNon12ScalesED317(tt *testing.T) {
	var s Scale
	var err error
	var t Tuning
	s, err = ScaleFromSCLFile(testFile("ED3-17.scl"))
	assert.NilError(tt, err)
	t, err = TuningFromSCL(s)
	assert.NilError(tt, err)
	for i := 0; i < 100; i++ {
		mn := int(rand.Uint32()%40 + 40)
		freq := 150.0 + 300.0*float64(rand.Uint32())/float64(math.MaxUint32)
		var k KeyboardMapping
		k, err = KeyboardMappingTuneNoteTo(mn, freq)
		assert.NilError(tt, err)
		var mapped Tuning
		mapped, err = TuningFromSCLAndKBM(s, k)
		assert.NilError(tt, err)

		assert.Equal(tt, "", approxEqual(1e-6, mapped.FrequencyForMidiNote(mn), freq), "mn:%v, freq:%v", mn, freq)

		// This scale is monotonic so test monotonicity still
		for ii := 1; ii < 127; ii++ {
			if mapped.FrequencyForMidiNote(ii) > 1.0 {
				assert.Check(tt, mapped.FrequencyForMidiNote(ii) > mapped.FrequencyForMidiNote(ii-1), "mn:%v, freq:%v, ii:%v", mn, freq, ii)
			}
		}
		n60ldiff := t.LogScaledFrequencyForMidiNote(60) - mapped.LogScaledFrequencyForMidiNote(60)
		for j := 0; j < 128; j++ {
			assert.Equal(tt, "", approxEqual(1e-6, t.LogScaledFrequencyForMidiNote(j)-mapped.LogScaledFrequencyForMidiNote(j),
				n60ldiff), "mn:%v, freq:%v, j:%v", mn, freq, j)
		}
	}
}

// KBMs with Gaps - 12 Intune with Gap
func TestKBMsWithGaps12Intune(tt *testing.T) {
	var err error
	var s Scale
	var k KeyboardMapping
	var t, tm Tuning

	s, err = ScaleFromSCLFile(testFile("12-intune.scl"))
	assert.NilError(tt, err)
	k, err = KeyboardMappingFromKBMFile(testFile("mapping-whitekeys-c261.kbm"))
	assert.NilError(tt, err)

	t, err = TuningFromSCL(s)
	assert.NilError(tt, err)
	tm, err = TuningFromSCLAndKBM(s, k)
	assert.NilError(tt, err)

	assert.Equal(tt, s.Count, 12)
	assert.Equal(tt, "", approxEqual(1e-6, t.FrequencyForMidiNote(69), 440.0))

	maps := map[int]int{
		60: 60,
		61: 62,
		62: 64,
		63: 65,
		64: 67,
		65: 69,
		66: 71,
	}
	for first, second := range maps {
		assert.Equal(tt, "", approxEqual(1e-5, t.LogScaledFrequencyForMidiNote(first),
			tm.LogScaledFrequencyForMidiNote(second)), "first:%d,second:%d", first, second)

	}
}

// KBM ReOrdering - Non Monotonic KBM note
func TestKBMReorderingNoMonotonicKBMNote(tt *testing.T) {
	var err error
	var s Scale
	var k KeyboardMapping
	var t Tuning

	s, err = ScaleFromSCLFile(testFile("12-intune.scl"))
	assert.NilError(tt, err)
	k, err = KeyboardMappingFromKBMFile(testFile("shuffle-a440-constant.kbm"))
	assert.NilError(tt, err)

	t, err = TuningFromSCLAndKBM(s, k)
	assert.NilError(tt, err)

	assert.Equal(tt, s.Count, 12)
	assert.Equal(tt, "", approxEqual(1e-6, t.FrequencyForMidiNote(69), 440.0))

	order := []float64{0, 2, 1, 3, 4, 6, 5, 7, 8, 9, 11, 10, 12}
	l60 := t.LogScaledFrequencyForMidiNote(60)
	for i, oi := range order {
		li := t.LogScaledFrequencyForMidiNote(60 + i)
		assert.Equal(tt, "", approxEqual(1e-6, li-l60, oi/12.0), "i:%d,oi:%v", i, oi)

	}
}

// Exceptions and Bad Files - Read Non-present files
func TestBadFilesNoPresent(tt *testing.T) {
	var err error
	_, err = ScaleFromSCLFile("blahlfdsfds")
	assert.ErrorContains(tt, err, "Unable to open file")
	_, err = KeyboardMappingFromKBMFile("blahlfdsfds")
	assert.ErrorContains(tt, err, "Unable to open file")
}

// Exceptions and Bad Files - Bad SCL
func TestBadFilesBadSCL(tt *testing.T) {
	var err error

	// Trailing data is OK
	_, err = ScaleFromSCLFile(testFile("bad/extraline.scl"))
	assert.NilError(tt, err)

	_, err = ScaleFromSCLFile(testFile("bad/badnote.scl"))
	assert.ErrorContains(tt, err, "Unable to parse file")
	_, err = ScaleFromSCLFile(testFile("bad/blanknote.scl"))
	assert.ErrorContains(tt, err, "Unable to parse file")
	_, err = ScaleFromSCLFile(testFile("bad/missingnote.scl"))
	assert.ErrorContains(tt, err, "Unable to parse file")
}

// Exceptions and Bad Files - Bad KBM
func TestBadFilesBadKBM(tt *testing.T) {
	var err error
	_, err = KeyboardMappingFromKBMFile(testFile("bad/blank-line.kbm"))
	assert.ErrorContains(tt, err, "Unable to parse file")
	_, err = KeyboardMappingFromKBMFile(testFile("bad/empty-bad.kbm"))
	assert.ErrorContains(tt, err, "Unable to parse file")
	_, err = KeyboardMappingFromKBMFile(testFile("bad/garbage-key.kbm"))
	assert.ErrorContains(tt, err, "Unable to parse file")

	_, err = KeyboardMappingFromKBMFile(testFile("bad/empty-extra.kbm"))
	assert.NilError(tt, err)
	_, err = KeyboardMappingFromKBMFile(testFile("bad/extraline-long.kbm"))
	assert.NilError(tt, err)

	_, err = KeyboardMappingFromKBMFile(testFile("bad/missing-note.kbm"))
	assert.ErrorContains(tt, err, "Unable to parse file")
}

// Built in Generators - ED2
func TestBuiltinGeneratorsED2(tt *testing.T) {
	var err error
	var s Scale
	s, err = ScaleEvenDivisionOfSpanByM(2, 12)
	assert.NilError(tt, err)
	assert.Equal(tt, s.Count, 12)
	assert.Check(tt, len(s.RawText) > 1)
	var t, ut Tuning
	ut, err = TuningEvenStandard()
	assert.NilError(tt, err)
	t, err = TuningFromSCL(s)
	assert.NilError(tt, err)
	//fmt.Printf("t: %#v\n\nut: %#v\n",t,ut)
	for i := 0; i < 128; i++ {
		//fmt.Printf("i:%v, tl:%v tf:%v ul:%v uf:%v\n", i,
		//	t.LogScaledFrequencyForMidiNote(i), t.FrequencyForMidiNote(i),
		//	ut.LogScaledFrequencyForMidiNote(i), ut.FrequencyForMidiNote(i))
		assert.Equal(tt, "", approxEqual(1e-6, t.LogScaledFrequencyForMidiNote(i), ut.LogScaledFrequencyForMidiNote(i)), "i:%d", i)
	}
}

// Built in Generators - ED3-17
func TestBuiltinGeneratorsED317(tt *testing.T) {
	var err error
	var s, sf Scale
	s, err = ScaleEvenDivisionOfSpanByM(3, 17)
	assert.NilError(tt, err)
	sf, err = ScaleFromSCLFile(testFile("ED3-17.scl"))
	assert.NilError(tt, err)
	var t, ut Tuning
	ut, err = TuningFromSCL(sf)
	assert.NilError(tt, err)
	t, err = TuningFromSCL(s)
	assert.NilError(tt, err)
	for i := 0; i < 128; i++ {
		assert.Equal(tt, "", approxEqual(1e-6, t.LogScaledFrequencyForMidiNote(i), ut.LogScaledFrequencyForMidiNote(i)), "i:%d", i)
	}
}

// Built in Generators - ED4-17
func TestBuiltinGeneratorsED417(tt *testing.T) {
	var err error
	var s, sf Scale
	s, err = ScaleEvenDivisionOfSpanByM(4, 17)
	assert.NilError(tt, err)
	sf, err = ScaleFromSCLFile(testFile("ED4-17.scl"))
	assert.NilError(tt, err)
	var t, ut Tuning
	ut, err = TuningFromSCL(sf)
	assert.NilError(tt, err)
	t, err = TuningFromSCL(s)
	assert.NilError(tt, err)
	for i := 0; i < 128; i++ {
		assert.Equal(tt, "", approxEqual(1e-6, t.LogScaledFrequencyForMidiNote(i), ut.LogScaledFrequencyForMidiNote(i)), "i:%d", i)
	}
}

// Built in Generators - Constraints on random EDN-M
func TestBuiltinGeneratorsRandonEDNMConstraints(tt *testing.T) {
	for i := 0; i < 100; i++ {
		span := int(rand.Uint32()%7 + 2)
		m := int(rand.Uint32()%50 + 3)

		var s Scale
		var t Tuning
		var err error

		s, err = ScaleEvenDivisionOfSpanByM(span, m)
		assert.NilError(tt, err)
		assert.Equal(tt, s.Count, m)
		assert.Check(tt, len(s.RawText) > 1)

		t, err = TuningFromSCL(s)
		assert.NilError(tt, err)
		assert.Equal(tt, "", approxEqual(1e-7, t.FrequencyForMidiNoteScaledByMidi0(60)*float64(span),
			t.FrequencyForMidiNoteScaledByMidi0(60+m)), "i:%v,span:%v,m:%v", i, span, m)
		d0 := t.LogScaledFrequencyForMidiNote(1) - t.LogScaledFrequencyForMidiNote(0)
		for j := 1; j < 128; j++ {
			d := t.LogScaledFrequencyForMidiNote(j) - t.LogScaledFrequencyForMidiNote(j-1)
			assert.Equal(tt, "", approxEqual(1e-7, d, d0), "i:%v,j:%v,span:%v,m:%v", i, j, span, m)
		}
	}
}

// Built in Generators - EDMN Errors
func TestBuiltinGeneratorsEDMNError(tt *testing.T) {
	var err error
	_, err = ScaleEvenDivisionOfSpanByM(0, 12)
	assert.ErrorContains(tt, err, "Span must be a positive number")
	_, err = ScaleEvenDivisionOfSpanByM(2, 0)
	assert.ErrorContains(tt, err, "must divide the period into at least one step")
	_, err = ScaleEvenDivisionOfSpanByM(0, 0)
	assert.ErrorContains(tt, err, "Span must be a positive number")

	_, err = ScaleEvenDivisionOfSpanByM(-1, 12)
	assert.ErrorContains(tt, err, "Span must be a positive number")
	_, err = ScaleEvenDivisionOfSpanByM(2, -1)
	assert.ErrorContains(tt, err, "must divide the period into at least one step")
	_, err = ScaleEvenDivisionOfSpanByM(-1, -1)
	assert.ErrorContains(tt, err, "Span must be a positive number")
}

// Dos Line Endings and Blanks - SCL
func TestDosSCL(tt *testing.T) {
	var err error
	_, err = ScaleFromSCLFile(testFile("12-intune-dosle.scl"))
	assert.NilError(tt, err)
}

// Dos Line Endings and Blanks - Properly read a file with DOS line endings
func TestDosDosEndings(tt *testing.T) {
	var err error
	var s Scale
	s, err = ScaleFromSCLFile(testFile("31edo_dos_lineends.scl"))
	assert.NilError(tt, err)
	assert.Equal(tt, s.Count, 31)
	// should not include \r:
	assert.Equal(tt, s.Description, "31 equal divisions of octave")

	// the parsing should ive the same floatvalues independent of crlf status obviously
	var q Scale
	q, err = ScaleFromSCLFile(testFile("31edo.scl"))
	assert.NilError(tt, err)
	for i := 0; i < q.Count; i++ {
		assert.Equal(tt, q.Tones[i].FloatValue, s.Tones[i].FloatValue)
	}
}

// Dos Line Endings and Blanks - KBM
func TestDosKBM(tt *testing.T) {
	var err error
	var k KeyboardMapping
	k, err = KeyboardMappingFromKBMFile(testFile("empty-note69-dosle.kbm"))
	assert.NilError(tt, err)
	assert.Equal(tt, k.TuningConstantNote, 69)
}

// Dos Line Endings and Blanks - Blank SCL
func TestDosBlankSCL(tt *testing.T) {
	var err error
	_, err = ScaleFromSCLString("")
	assert.ErrorContains(tt, err, "Incomplete SCL")

	// but what if we do construct a bad one?
	var s Scale
	s.Count = 0
	s.Tones = nil
	_, err = TuningFromSCL(s)
	assert.ErrorContains(tt, err, "Your scale provided 0 notes")
}

// Tone API - Valid Tones
func TestToneAPIValid(tt *testing.T) {
	var err error
	var t Tone
	t, err = toneFromString("130.0", 1)
	assert.NilError(tt, err)
	assert.Equal(tt, t.Type, ToneCents)
	assert.Equal(tt, t.Cents, 130.0)
	assert.Equal(tt, t.FloatValue, 130.0/1200.0+1.0)

	t, err = toneFromString("7/2", 1)
	assert.NilError(tt, err)
	assert.Equal(tt, t.Type, ToneRatio)
	assert.Equal(tt, t.RatioN, 7)
	assert.Equal(tt, t.RatioD, 2)
	assert.Equal(tt, t.FloatValue, math.Log(7.0/2.0)/math.Log(2.0)+1.0)

	t, err = toneFromString("3", 1)
	assert.NilError(tt, err)
	assert.Equal(tt, t.Type, ToneRatio)
	assert.Equal(tt, t.RatioN, 3)
	assert.Equal(tt, t.RatioD, 1)
	assert.Equal(tt, t.FloatValue, math.Log(3.0/1.0)/math.Log(2.0)+1.0)
}

// Tone API - Error Tones
func TestToneAPIErrors(tt *testing.T) {
	var err error
	_, err = toneFromString("Not a number", 1)
	assert.ErrorContains(tt, err, "Error parsing")

	// the following are commented out in the C++, but test cleanly for Go:
	_, err = toneFromString("100.200 with extra stuff", 1)
	assert.ErrorContains(tt, err, "Error parsing")
	_, err = toneFromString("7/4/2", 1)
	assert.ErrorContains(tt, err, "Error parsing")
	_, err = toneFromString("7*2", 1)
	assert.ErrorContains(tt, err, "Error parsing")
}

// Scale Position - Untuned
func TestScalePositionUntuned(tt *testing.T) {
	var t Tuning
	var err error
	t, err = TuningEvenStandard()
	assert.NilError(tt, err)
	for i := 0; i < 127; i++ {
		assert.Equal(tt, t.ScalePositionForMidiNote(i), i%12, "i:%d", i)
	}
}

// Scale Position - Untuned, Mapped
func TestScalePositionUntunedMapped(tt *testing.T) {
	var err error
	{
		var t Tuning
		var k KeyboardMapping
		k, err = KeyboardMappingStartScaleOnAndTuneNoteTo(60, 69, 440.0)
		assert.NilError(tt, err)
		t, err = TuningFromKBM(k)
		assert.NilError(tt, err)
		for i := 0; i < 127; i++ {
			assert.Equal(tt, t.ScalePositionForMidiNote(i), i%12, "i:%d", i)
		}
	}
	for j := 0; j < 100; j++ {
		n := int(rand.Uint32()%60 + 30)
		var t Tuning
		var k KeyboardMapping
		k, err = KeyboardMappingStartScaleOnAndTuneNoteTo(n, 69, 440.0)
		assert.NilError(tt, err)
		t, err = TuningFromKBM(k)
		assert.NilError(tt, err)
		for i := 0; i < 127; i++ {
			assert.Equal(tt, t.ScalePositionForMidiNote(i), (i+12-n%12)%12, "n:%d,i:%d", n, i)
		}
	}
	{
		// Check whitekeys
		var t Tuning
		var k KeyboardMapping
		k, err = KeyboardMappingFromKBMFile(testFile("mapping-whitekeys-c261.kbm"))
		assert.NilError(tt, err)
		t, err = TuningFromKBM(k)
		assert.NilError(tt, err)
		maps := map[int]int{
			0:  0,
			2:  1,
			4:  2,
			5:  3,
			7:  4,
			9:  5,
			11: 6,
		}
		for i := 0; i < 127; i++ {
			spn := t.ScalePositionForMidiNote(i)
			expected, exists := maps[i%12]
			if !exists {
				expected = -1
			}
			assert.Equal(tt, spn, expected, "i:%d", i)
		}
	}
}

// Scale Position - Tuned, Unmapped
func TestScalePositionTunedUnmapped(tt *testing.T) {
	var err error
	// Check longer and shorter scales
	{
		var t Tuning
		var s Scale
		s, err = ScaleFromSCLFile(testFile("zeus22.scl"))
		assert.NilError(tt, err)
		t, err = TuningFromSCL(s)
		assert.NilError(tt, err)
		off := 60
		for off > 0 {
			off -= s.Count
		}
		for i := 0; i < 127; i++ {
			assert.Equal(tt, t.ScalePositionForMidiNote(i), (i-off)%s.Count, "i:%d", i)
		}
	}
	// Check longer and shorter scales
	{
		var t Tuning
		var s Scale
		s, err = ScaleFromSCLFile(testFile("6-exact.scl"))
		assert.NilError(tt, err)
		t, err = TuningFromSCL(s)
		assert.NilError(tt, err)
		off := 60
		for off > 0 {
			off -= s.Count
		}
		for i := 0; i < 127; i++ {
			assert.Equal(tt, t.ScalePositionForMidiNote(i), (i-off)%s.Count, "i:%d", i)
		}
	}
}

// Scale Position - Tuned, Mapped
func TestScalePositionTunedMapped(tt *testing.T) {
	// And check some combos
	var err error
	var t Tuning
	var k KeyboardMapping
	var s Scale
	for j := 0; j < 100; j++ {
		n := int(rand.Uint32()%60 + 30)

		s, err = ScaleFromSCLFile(testFile("zeus22.scl"))
		assert.NilError(tt, err)
		k, err = KeyboardMappingStartScaleOnAndTuneNoteTo(n, 69, 440.0)
		assert.NilError(tt, err)
		t, err = TuningFromSCLAndKBM(s, k)
		assert.NilError(tt, err)
		off := n
		for off > 0 {
			off -= s.Count
		}
		for i := 0; i < 127; i++ {
			assert.Equal(tt, t.ScalePositionForMidiNote(i), (i-off)%s.Count, "n:%d,i:%d", n, i)
		}
	}
}

// Default KBM Constructor has Right Base - All Scales with Default KBM
func TestDefaultKBMHasRightBase(tt *testing.T) {
	for _, fname := range testSCLs {
		var s Scale
		var t Tuning
		var err error
		s, err = ScaleFromSCLFile(testFile(fname))
		assert.NilError(tt, err, fname)
		t, err = TuningFromSCL(s)
		assert.NilError(tt, err, fname)
		assert.Equal(tt, "", approxEqual(1e-6, t.FrequencyForMidiNoteScaledByMidi0(60), 32.0), fname)
	}
}

// Different KBM period from Scale period - 31Edo with mean Tone mapping
func TestDiffPeriods32EdoMeanTone(tt *testing.T) {
	var s Scale
	var k KeyboardMapping
	var t Tuning
	var err error
	s, err = ScaleFromSCLFile(testFile("31edo.scl"))
	assert.NilError(tt, err)
	k, err = KeyboardMappingFromKBMFile(testFile("31edo_meantone.kbm"))
	assert.NilError(tt, err)
	t, err = TuningFromSCLAndKBM(s, k)
	assert.NilError(tt, err)

	assert.Equal(tt, "", approxEqual(1e-6, t.FrequencyForMidiNote(69), 440.0))
	assert.Equal(tt, "", approxEqual(1e-6, t.FrequencyForMidiNote(69+12), 880.0))
}

// Different KBM period from Scale period - Perfect 5th UnMapped
func TestDiffPeriodsPerfect5thUnmapped(tt *testing.T) {
	var s Scale
	var t Tuning
	var err error
	s, err = ScaleFromSCLFile(testFile("12-ET-P5.scl"))
	assert.NilError(tt, err)
	t, err = TuningFromSCL(s)
	assert.NilError(tt, err)

	for i := 60 - 36; i < 127; i += 12 {
		f := t.FrequencyForMidiNote(i)
		f5 := t.FrequencyForMidiNote(i + 7)
		assert.Equal(tt, "", approxEqual(1e-6, f5, f*1.5), "i:%d", i)
	}
}

// Different KBM period from Scale period - Perfect 5th 07 mapping
func TestDiffPeriodsPerfect5th07Mapping(tt *testing.T) {
	var s Scale
	var k KeyboardMapping
	var t Tuning
	var err error
	s, err = ScaleFromSCLFile(testFile("12-ET-P5.scl"))
	assert.NilError(tt, err)
	k, err = KeyboardMappingFromKBMFile(testFile("mapping-n60-fifths.kbm"))
	assert.NilError(tt, err)
	t, err = TuningFromSCLAndKBM(s, k)
	assert.NilError(tt, err)

	for i := 60; i < 70; i += 2 {
		f := t.FrequencyForMidiNote(i)
		f5 := t.FrequencyForMidiNote(i + 1)
		assert.Equal(tt, "", approxEqual(1e-6, f5, f*1.5), "i:%d", i)
	}
}

// Skipped Note API tests
// Default Tuning skips Nothing
func TestSkippedNoteDefaultTuningSkipsNothing(tt *testing.T) {
	t, err := TuningEvenStandard()
	assert.NilError(tt, err)
	for i := 0; i < 128; i++ {
		assert.Assert(tt, t.IsMidiNoteMapped(i), i)
	}
}

// SCL-only Tuning skips Nothing
func TestSkippedNoteSCLonlyTuningSkipsNothing(tt *testing.T) {
	s, err := ScaleFromSCLFile(testFile("12-intune.scl"))
	assert.NilError(tt, err)
	t, err := TuningFromSCL(s)
	assert.NilError(tt, err)
	for i := 0; i < 128; i++ {
		assert.Assert(tt, t.IsMidiNoteMapped(i), i)
	}
}

// KBM-only Tuning absent skips skips Nothing
func TestSkippedNoteKBMonlyTuningAbsentSkipsSkipsNothing(tt *testing.T) {
	k, err := KeyboardMappingFromKBMFile(testFile("empty-note69.kbm"))
	assert.NilError(tt, err)
	t, err := TuningFromKBM(k)
	assert.NilError(tt, err)
	for i := 0; i < 128; i++ {
		assert.Assert(tt, t.IsMidiNoteMapped(i), i)
	}
}

// Fully Mapped
func TestSkippedNoteFullyMapped(tt *testing.T) {
	s, err := ScaleFromSCLFile(testFile("12-intune.scl"))
	assert.NilError(tt, err)
	k, err := KeyboardMappingFromKBMFile(testFile("empty-note69.kbm"))
	assert.NilError(tt, err)
	t, err := TuningFromSCLAndKBM(s, k)
	assert.NilError(tt, err)
	for i := 0; i < 128; i++ {
		assert.Assert(tt, t.IsMidiNoteMapped(i), i)
	}
}

// Gaps in the Maps
func TestSkippedNoteGapsInTheMaps(tt *testing.T) {
	s, err := ScaleFromSCLFile(testFile("12-intune.scl"))
	assert.NilError(tt, err)
	k, err := KeyboardMappingFromKBMFile(testFile("mapping-whitekeys-a440.kbm"))
	assert.NilError(tt, err)
	t, err := TuningFromSCLAndKBM(s, k)
	assert.NilError(tt, err)
	for k := 0; k < 128; k++ {
		i := k % 12
		isOn := i == 0 || i == 2 || i == 4 || i == 5 || i == 7 || i == 9 || i == 11
		assert.Assert(tt, t.IsMidiNoteMapped(k) == isOn, k)
	}
}

// Gaps in the Maps KBM Only
func TestSkippedNoteGapsInTheMapsKBMOnly(tt *testing.T) {
	k, err := KeyboardMappingFromKBMFile(testFile("mapping-whitekeys-a440.kbm"))
	assert.NilError(tt, err)
	t, err := TuningFromKBM(k)
	assert.NilError(tt, err)
	for k := 0; k < 128; k++ {
		i := k % 12
		isOn := i == 0 || i == 2 || i == 4 || i == 5 || i == 7 || i == 9 || i == 11
		assert.Assert(tt, t.IsMidiNoteMapped(k) == isOn, k)
	}
}

// Tuning with Gaps
func TestSkippedNoteTuningWithGaps(tt *testing.T) {
	s, err := ScaleFromSCLFile(testFile("12-intune.scl"))
	assert.NilError(tt, err)
	k, err := KeyboardMappingFromKBMFile(testFile("mapping-whitekeys-a440.kbm"))
	assert.NilError(tt, err)
	t, err := TuningFromSCLAndKBM(s, k)
	assert.NilError(tt, err)
	for k := 2; k < 128; k++ {
		i := k % 12
		isOn := i == 0 || i == 2 || i == 4 || i == 5 || i == 7 || i == 9 || i == 11
		var priorOn int
		if i == 0 || i == 5 {
			priorOn = 1
		} else {
			priorOn = 2
		}
		if isOn {
			assert.Assert(tt, t.LogScaledFrequencyForMidiNote(k) > t.LogScaledFrequencyForMidiNote(k-priorOn), k)
		}
	}
}

// Tuning with Gaps and Interpolation
func TestSkippedNoteTuningWithGapsAndInterpolation(tt *testing.T) {
	s, err := ScaleFromSCLFile(testFile("12-intune.scl"))
	assert.NilError(tt, err)
	k, err := KeyboardMappingFromKBMFile(testFile("mapping-whitekeys-a440.kbm"))
	assert.NilError(tt, err)
	t, err := TuningFromSCLAndKBM(s, k)
	t = t.WithSkippedNotesInterpolated()
	assert.NilError(tt, err)
	for k := 2; k < 128; k++ {
		// Now we have filled in, all the APIs should be monotonic
		i := k % 12
		isOn := i == 0 || i == 2 || i == 4 || i == 5 || i == 7 || i == 9 || i == 11
		assert.Assert(tt, t.IsMidiNoteMapped(k) == isOn, k)
		assert.Assert(tt, t.LogScaledFrequencyForMidiNote(k) > t.LogScaledFrequencyForMidiNote(k-1), k)
		assert.Assert(tt, t.FrequencyForMidiNote(k) > t.FrequencyForMidiNote(k-1), k)
		assert.Assert(tt, t.FrequencyForMidiNoteScaledByMidi0(k) > t.FrequencyForMidiNoteScaledByMidi0(k-1), k)
	}
}
