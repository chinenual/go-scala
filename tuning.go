package scala

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
type Tuning struct {
	Scale           Scale
	KeyboardMapping KeyboardMapping
}

// CreateStandardTuning constructs a tuning with even temperament and standard mapping
func CreateStandardTuning() (t Tuning, err error) {
	return
}

// CreateTuningFromSCL constructs a tuning for a particular scale.
func CreateTuningFromSCL(s Scale) (t Tuning, err error) {
	return
}

// CreateTuningFromKBD constructs a tuning for a particular mapping.
func CreateTuningFromKBD(k KeyboardMapping) (t Tuning, err error) {
	return
}

// CreateTuningFromSCLAndKBD constructs a tuning for a particular scale and mapping
func CreateTuningFromSCLAndKBD(s Scale, k KeyboardMapping) (t Tuning, err error) {
	return
}

// FrequencyForMidiNote returns the Frequency in HZ for a given midi
// note. In standard tuning, FrequencyForMidiNote(69) will be 440
// and frequencyForMidiNote(60) will be 261.62 - the standard frequencies
// for A and middle C.
func (t Tuning) FrequencyForMidiNote(mn int) float64 {
	return 0
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
func (t Tuning) FrequencyForMidiNoteScaledByMidi0(mn int) float64 {
	return 0
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
func (t Tuning) LogScaledFrequencyForMidiNote(mn int) float64 {
	return 0
}

// ScalePositionForMidiNote returns the space in the logical scale. Note 0 is the root.
// It has a maximum value of count-1. Note that SCL files omit the root internally and so
// this logical scale position is off by 1 from the index in the tones array of the Scale data.
func (t Tuning) ScalePositionForMidiNote(mn int) int {
	return 0
}
