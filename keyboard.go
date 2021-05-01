package scala

import "io"

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
