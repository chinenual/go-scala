package scala

import (
	"github.com/pkg/errors"
	"math"
)

// The Tuning type is the primary place where you will interact with this library.
// It is constructed for a scale and mapping and then gives you the ability to
// determine frequencies across and beyond the midi keyboard. Since modulation
// can force key number well outside the [0,127] range in some of our synths we
// support a midi note range from -256 to + 256 spanning more than the entire frequency
// space reasonable.
//
// To use this type, you construct a fresh instance every time you want to use a
// different Scale and Keyboard. If you want to tune to a different scale or mapping,
// just construct a new instance.
type Tuning interface {
	// FrequencyForMidiNote returns the Frequency in HZ for a given midi
	// note. In standard tuning, FrequencyForMidiNote(69) will be 440
	// and frequencyForMidiNote(60) will be 261.62 - the standard frequencies
	// for A and middle C.
	FrequencyForMidiNote(mn int) float64

	// FrequencyForMidiNoteScaledByMidi0 returns the frequency but with the
	// standard frequency of midi note 0 divided out. So in standard tuning
	// frequencyForMidiNoteScaledByMidi0(0) = 1 and frequencyForMidiNoteScaledByMidi0(60) = 32
	//
	// Both the frequency measures have the feature of doubling when frequency doubles
	// (or when a standard octave is spanned), whereas the log one increase by 1 per frequency double.
	//
	// Depending on your internal pitch model, one of these three methods should allow you
	// to calibrate your oscillators to the appropriate frequency based on the midi note
	// at hand.
	FrequencyForMidiNoteScaledByMidi0(mn int) float64

	// LogScaledFrequencyForMidiNote returns the log base 2 of the scaled frequency.
	// So logScaledFrequencyForMidiNote(0) = 0 and logScaledFrequencyForMidiNote(60) = 5.
	//
	// Both the frequency measures have the feature of doubling when frequency doubles
	// (or when a standard octave is spanned), whereas the log one increase by 1 per frequency double.
	//
	// Depending on your internal pitch model, one of these three methods should allow you
	// to calibrate your oscillators to the appropriate frequency based on the midi note
	// at hand.
	LogScaledFrequencyForMidiNote(mn int) float64

	// ScalePositionForMidiNote returns the space in the logical scale. Note 0 is the root.
	// It has a maximum value of count-1. Note that SCL files omit the root internally and so
	// this logical scale position is off by 1 from the index in the tones array of the Scale data.
	ScalePositionForMidiNote(mn int) int

	// For convenience, the scale and mapping used to construct this are kept as public copies
	Scale() Scale
	KeyboardMapping() KeyboardMapping
}

type TuningImpl struct {
	scale              Scale
	keyboardMapping    KeyboardMapping
	lptable            [numPrecomputed]float64
	ptable             [numPrecomputed]float64
	scalePositionTable [numPrecomputed]int
}

const numPrecomputed = 512

// CreateStandardTuning constructs a tuning with even temperament and standard mapping
func CreateStandardTuning() (t Tuning, err error) {
	var k KeyboardMapping
	var s Scale
	if k, err = standardKeyboardMapping(); err != nil {
		return
	}
	if s, err = evenTemperment12NoteScale(); err != nil {
		return
	}
	if t, err = CreateTuningFromSCLAndKBM(s, k); err != nil {
		return
	}
	return
}

// CreateTuningFromSCL constructs a tuning for a particular scale.
func CreateTuningFromSCL(s Scale) (t Tuning, err error) {
	var k KeyboardMapping
	if k, err = standardKeyboardMapping(); err != nil {
		return
	}
	if t, err = CreateTuningFromSCLAndKBM(s, k); err != nil {
		return
	}
	return
}

// CreateTuningFromKBM constructs a tuning for a particular mapping.
func CreateTuningFromKBM(k KeyboardMapping) (t Tuning, err error) {
	var s Scale
	if s, err = evenTemperment12NoteScale(); err != nil {
		return
	}
	if t, err = CreateTuningFromSCLAndKBM(s, k); err != nil {
		return
	}
	return
}

// CreateTuningFromSCLAndKBM constructs a tuning for a particular scale and mapping
func CreateTuningFromSCLAndKBM(s Scale, k KeyboardMapping) (tuning Tuning, err error) {
	var t TuningImpl

	t.scale = s
	t.keyboardMapping = k
	if s.Count == 0 {
		err = errors.Errorf("Unable to tune to a scale with no notes. Your scale provided 0 notes.")
		return
	}
	var pitches [numPrecomputed]float64

	posPitch0 := 256 + k.TuningConstantNote
	posScale0 := 256 + k.MiddleNote

	pitchMod := math.Log(k.TuningPitch)/math.Log(2.0) - 1.0

	scalePositionOfTuningNote := k.TuningConstantNote - k.MiddleNote

	if k.Count > 0 {
		scalePositionOfTuningNote = k.Keys[scalePositionOfTuningNote]
	}
	tuningCenterPitchOffset := 0.0

	if scalePositionOfTuningNote != 0 {
		tshift := 0.0
		dt := s.Tones[s.Count-1].FloatValue - 1.0
		for scalePositionOfTuningNote < 0 {
			scalePositionOfTuningNote += s.Count
			tshift += dt
		}
		for scalePositionOfTuningNote > s.Count {
			scalePositionOfTuningNote -= s.Count
			tshift -= dt
		}

		if scalePositionOfTuningNote == 0 {
			tuningCenterPitchOffset = -tshift
		} else {
			tuningCenterPitchOffset = s.Tones[scalePositionOfTuningNote-1].FloatValue - 1.0 - tshift
		}
	}

	for i := 0; i < numPrecomputed; i++ {

		// TODO: ScaleCenter and PitchCenter are now two different notes.
		distanceFromPitch0 := i - posPitch0
		distanceFromScale0 := i - posScale0

		if distanceFromPitch0 == 0 {
			pitches[i] = 1
			t.lptable[i] = pitches[i] + pitchMod
			t.ptable[i] = math.Pow(2.0, t.lptable[i])
			t.scalePositionTable[i] = scalePositionOfTuningNote % s.Count
		} else {
			/*
			   We used to have this which assumed 1-12
			   Now we have our note number, our distance from the
			   center note, and the key remapping
			   int rounds = (distanceFromScale0-1) / s.count
			   int thisRound = (distanceFromScale0-1) % s.count
			*/

			var rounds int
			var thisRound int
			disable := false

			if k.Count == 0 {
				rounds = (distanceFromScale0 - 1) / s.Count
				thisRound = (distanceFromScale0 - 1) % s.Count
			} else {
				/*
				 ** Now we have this situation. We are at note i so we
				 ** are m away from the center note which is distanceFromScale0
				 **
				 ** If we mod that by the mapping size we know which note we are on
				 */
				mappingKey := distanceFromScale0 % k.Count
				if mappingKey < 0 {
					mappingKey += k.Count
				}
				// Now have we gone off the end
				rotations := 0
				dt := distanceFromScale0
				if dt > 0 {
					for dt >= k.Count {
						dt -= k.Count
						rotations++
					}
				} else {
					for dt < 0 {
						dt += k.Count
						rotations--
					}
				}

				cm := k.Keys[mappingKey]
				push := 0
				if cm < 0 {
					disable = true
				} else {
					push = mappingKey - cm
				}

				if k.OctaveDegrees > 0 && k.OctaveDegrees != k.Count {
					rounds = rotations
					thisRound = cm - 1
					if thisRound < 0 {
						thisRound = (k.OctaveDegrees - 1) % s.Count
						rounds--
					}
				} else {
					rounds = (distanceFromScale0 - push - 1) / s.Count
					thisRound = (distanceFromScale0 - push - 1) % s.Count
				}
			}

			if thisRound < 0 {
				thisRound += s.Count
				rounds -= 1
			}

			if disable {
				pitches[i] = 0
				t.scalePositionTable[i] = -1
			} else {
				pitches[i] = s.Tones[thisRound].FloatValue + float64(rounds)*(s.Tones[s.Count-1].FloatValue-1.0) - tuningCenterPitchOffset
				t.scalePositionTable[i] = (thisRound + 1) % s.Count
			}

			t.lptable[i] = pitches[i] + pitchMod
			t.ptable[i] = math.Pow(2.0, pitches[i]+pitchMod)

		}
	}
	tuning = t
	return
}

func imin(x int, y int) int {
	if x < y {
		return x
	}
	return y
}
func imax(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

// FrequencyForMidiNote returns the Frequency in HZ for a given midi
// note. In standard tuning, FrequencyForMidiNote(69) will be 440
// and frequencyForMidiNote(60) will be 261.62 - the standard frequencies
// for A and middle C.
func (t TuningImpl) FrequencyForMidiNote(mn int) float64 {
	mni := imin(imax(0, mn+256), numPrecomputed-1)
	return t.ptable[mni] * Midi0Freq
}

// FrequencyForMidiNoteScaledByMidi0 returns the frequency but with the
// standard frequency of midi note 0 divided out. So in standard tuning
// frequencyForMidiNoteScaledByMidi0(0) = 1 and frequencyForMidiNoteScaledByMidi0(60) = 32
// 
// Both the frequency measures have the feature of doubling when frequency doubles
// (or when a standard octave is spanned), whereas the log one increase by 1 per frequency double.
// 
// Depending on your internal pitch model, one of these three methods should allow you
// to calibrate your oscillators to the appropriate frequency based on the midi note
// at hand.
func (t TuningImpl) FrequencyForMidiNoteScaledByMidi0(mn int) float64 {
	mni := imin(imax(0, mn+256), numPrecomputed-1)
	return t.ptable[mni]
}

// LogScaledFrequencyForMidiNote returns the log base 2 of the scaled frequency.
// So logScaledFrequencyForMidiNote(0) = 0 and logScaledFrequencyForMidiNote(60) = 5.
// 
// Both the frequency measures have the feature of doubling when frequency doubles
// (or when a standard octave is spanned), whereas the log one increase by 1 per frequency double.
// 
// Depending on your internal pitch model, one of these three methods should allow you
// to calibrate your oscillators to the appropriate frequency based on the midi note
// at hand.
func (t TuningImpl) LogScaledFrequencyForMidiNote(mn int) float64 {
	mni := imin(imax(0, mn+256), numPrecomputed-1)
	return t.lptable[mni]
}

// ScalePositionForMidiNote returns the space in the logical scale. Note 0 is the root.
// It has a maximum value of count-1. Note that SCL files omit the root internally and so
// this logical scale position is off by 1 from the index in the tones array of the Scale data.
func (t TuningImpl) ScalePositionForMidiNote(mn int) int {
	mni := imin(imax(0, mn+256), numPrecomputed-1)
	return t.scalePositionTable[mni]
}

func (t TuningImpl) Scale() Scale {
	return t.scale
}

func (t TuningImpl) KeyboardMapping() KeyboardMapping {
	return t.keyboardMapping
}
