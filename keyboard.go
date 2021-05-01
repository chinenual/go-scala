package scala

import "io"

// KeyboardMapping represents a KBM file. In most cases, the salient
// features are the tuningConstantNote and tuningFrequency, which allow you to
// pick a fixed note in the midi keyboard when retuning. The KBM file can also
// remap individual keys to individual points in a scale, which here is done with the
// keys vector.
//
// Just as with Scale, the rawText member contains the text of the KBM file used.
type KeyboardMapping struct {
	Count              int
	FirstMidi          int
	LastMidi           int
	MiddleNote         int
	TuningConstantNote int
	TuningFrequency    float64
	TuningPitch        float64 // pitch = frequency / MIDI_0_FREQ
	OctaveDegrees      int
	Keys               []int // rather than an 'x' we use a '-1' for skipped keys
	RawText            string
	Name               string
}

// ReadKBMStream returns a KeyboardMapping from a KBM input stream
func ReadKBMStream(rdr io.Reader) (kbm KeyboardMapping, err error) {
	return
}

// ReadKBMFile returns a KeyboardMapping from a KBM file name
func ReadKBMFile(fname string) (kbm KeyboardMapping, err error) {
	return
}

// ParseKBMData returns a KeyboardMapping from a KBM data in memory
func ParseKBMData(kbmContents string) (kbm KeyboardMapping, err error) {
	return
}

// TuneA69To creates a KeyboardMapping which keeps the midi note 69 (A4) set
// to a constant frequency, given
func TuneA69To(freq float64) (kbm KeyboardMapping, err error) {
	return
}

// TuneNoteTo creates a KeyboardMapping which keeps the midi note given is set
// to a constant frequency, given
func TuneNoteTo(midiNote int, freq float64) (kbm KeyboardMapping, err error) {
	return
}

// StartScaleOnAndTuneNoteTo generates a KBM where scaleStart is the note 0
// of the scale, where midiNote is the tuned note, and where feq is the frequency
func StartScaleOnAndTuneNoteTo(scaleStart int, midiNote int, freq float64) (kbm KeyboardMapping, err error) {
	return
}
