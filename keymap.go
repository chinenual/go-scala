package scala

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

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
	type stateType int
	const (
		mapSize stateType = iota
		firstMidi
		lastMidi
		middle
		reference
		freq
		degree
		keys
		trailing
	)

	state := mapSize
	scanner := bufio.NewScanner(rdr)
	lineno := 0
	for scanner.Scan() {
		line := scanner.Text()
		kbm.RawText = kbm.RawText + "\n" + line
		lineno++
		line = strings.TrimRight(line, "\t ")
		if len(line) > 0 && line[0] == '!' {
			continue
		}
		if line == "x" {
			line = "-1"
		} else if state != trailing {
			for i := 0; i < len(line); i++ {
				r := rune(line[i])
				// difference vs. C++ - the scanner strips line endings so no need to check for CR and LF
				if !(line[i] == ' ' || unicode.IsDigit(r) || line[i] == '.') {
					err = errors.Errorf("Invalid line %d.  line=\"%s\". Bad character is '%c'/%d",lineno,line,line[i],line[i])
					return
				}
			}
		}
		asInt := func() (i int, err error) {
			var v int64
			if v, err = strconv.ParseInt(line, 10, 32); err != nil {
				err = errors.Wrapf(err, "Invalid line %d.  line=\"%s\". Could not parse as a integer number", lineno, line)
				return
			}
			i = int(v)
			return
		}
		asFloat := func() (v float64, err error) {
			if v, err = strconv.ParseFloat(line, 64); err != nil {
				err = errors.Wrapf(err, "Invalid line %d.  line=\"%s\". Could not parse as a floating point number", lineno, line)
				return
			}
			return
		}
		switch state {
		case mapSize:
			if kbm.Count, err = asInt(); err != nil {
				return
			}
		case firstMidi:
			if kbm.FirstMidi, err = asInt(); err != nil {
				return
			}
		case lastMidi:
			if kbm.LastMidi, err = asInt(); err != nil {
				return
			}
		case middle:
			if kbm.MiddleNote, err = asInt(); err != nil {
				return
			}
		case reference:
			if kbm.TuningConstantNote, err = asInt(); err != nil {
				return
			}
		case freq:
			if kbm.TuningFrequency, err = asFloat(); err != nil {
				return
			}
			kbm.TuningPitch = kbm.TuningFrequency / 8.17579891564371
		case degree:
			if kbm.OctaveDegrees, err = asInt(); err != nil {
				return
			}
		case keys:
			var i int
			if i, err = asInt(); err != nil {
				return
			}
			kbm.Keys = append(kbm.Keys, i)
		case trailing:
		}
		if ! ( state == keys || state == trailing )  {
			state = state + 1
		}
		if state == keys && kbm.Count == 0  {
			state = trailing
		}
	}
	if err = scanner.Err(); err != nil {
		return
	}
	if ! (state == keys || state == trailing) {
		err = errors.Errorf("Incomplete KBM file.  Unable to get keys section of file.")
		return
	}
	if len(kbm.Keys) != kbm.Count {
		err = errors.Errorf("Different number of keys than mapping file indicates. Count is %d and we parsed %d keys",
			kbm.Count, len(kbm.Keys))
		return
	}
	return
}

// ReadKBMFile returns a KeyboardMapping from a KBM file name
func ReadKBMFile(fname string) (kbm KeyboardMapping, err error) {
	var file *os.File
	if file, err = os.Open(fname); err != nil {
		return
	}
	defer file.Close()
	if kbm, err = ReadKBMStream(file); err != nil {
		return
	}
	kbm.Name = fname
	return
}

// ParseKBMData returns a KeyboardMapping from a KBM data in memory
func ParseKBMData(kbmContents string) (kbm KeyboardMapping, err error) {
	rdr := strings.NewReader(kbmContents)
	if kbm, err = ReadKBMStream(rdr); err != nil {
		return
	}
	kbm.Name = "Mapping from patch"
	return
}

// tuneA69To creates a KeyboardMapping which keeps the midi note 69 (A4) set
// to a constant frequency, given
func tuneA69To(freq float64) (kbm KeyboardMapping, err error) {
	kbm, err = tuneNoteTo(69, freq)
	return
}

// tuneNoteTo creates a KeyboardMapping which keeps the midi note given is set
// to a constant frequency, given
func tuneNoteTo(midiNote int, freq float64) (kbm KeyboardMapping, err error) {
	kbm, err = startScaleOnAndTuneNoteTo(69, midiNote, freq)
	return
}

// startScaleOnAndTuneNoteTo generates a KBM where scaleStart is the note 0
// of the scale, where midiNote is the tuned note, and where feq is the frequency
func startScaleOnAndTuneNoteTo(scaleStart int, midiNote int, freq float64) (kbm KeyboardMapping, err error) {
	buf := "! Automatically generated mapping, tuning note " + strconv.Itoa(midiNote) + " to " + fmt.Sprintf("%f", freq) + " Hz\n"
	buf += "!\n"
	buf += "! Size of map\n"
	buf += "0\n"
	buf += "! First and last MIDI notes to map - map the entire keyboard\n"
	buf += "0\n"
	buf += "127\n"
	buf += "! Middle note where the first entry in the scale is mapped.\n"
	buf += strconv.Itoa(scaleStart) + "\n"
	buf += "! Reference note where frequency is fixed\n"
	buf += strconv.Itoa(midiNote) + "\n"
	buf += "! Frequency for MIDI note " + strconv.Itoa(midiNote) + "\n"
	buf += fmt.Sprintf("%f", freq) + "\n"
	buf += "! Scale degree for formal octave. This is am empty mapping, so:\n"
	buf += "0\n"
	buf += "! Mapping. This is an empty mapping so list no keys\n"
	kbm, err = ParseKBMData(buf)
	return
}

func standardKeyboardMapping() (kbm KeyboardMapping, err error) {
	freq := Midi0Freq * 32.0
	kbm,err = ParseKBMData(`! Default KBM file
0
0
127
60
60
` + fmt.Sprintf("%f", freq) + `
0
`)
	return
}