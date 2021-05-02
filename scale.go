package scala

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)


const Midi0Freq = 8.17579891564371 // or 440.0 * pow( 2.0, - (69.0/12.0 ) )

type ToneType int
const (
	ToneCents ToneType = iota
	ToneRatio
)

// A Tone is a single entry in an SCL file. It is expressed either in cents or in
// a ratio, as described in the SCL documentation.
//
// In most normal use, you will not use this interface, and it will be internal to a Scale
type Tone struct {
	Type       ToneType
	Cents      float64
	RatioD     int
	RatioN     int
	StringRep  string
	FloatValue float64 // cents / 1200 + 1.
}

// The Scale is the representation of the SCL file. It contains several key
// features. Most importantly it has a count and a vector of Tones.
//
// In most normal use, you will simply pass around instances of this interface
// to a Tunings instance, but in some cases you may want to create
// or inspect this class yourself. Especially if you are displaying this
// object to your end users, you may want to use the rawText or count methods.
type Scale struct {
	Name        string // The name in the SCL file. Informational only
	Description string // The description in the SCL file. Informational only
	RawText     string // The raw text of the SCL file used to create this Scale
	Count       int    // The number of tones.
	Tones       []Tone // The tones
}





func toneFromString(line string, lineno int) (tone Tone, err error) {
	if strings.Contains(line, ".") {
		tone.Type = ToneCents
		if tone.Cents,err = strconv.ParseFloat(strings.TrimSpace(line), 64); err != nil {
			err = errors.Wrapf(err, "Error parsing scale cent: \"%s\", line %d", line, lineno)
			return
		}
	} else {
		var v int64
		tone.Type = ToneRatio
		split := strings.Split(line,"/")
		if split != nil && len(split) == 1 {
			if v, err = strconv.ParseInt(strings.TrimSpace(split[0]), 10, 32); err != nil {
				err = errors.Errorf("Error parsing scale ratio numerator: \"%s\", line %d", split[0], lineno)
				return
			}
			tone.RatioN = int(v)
			tone.RatioD = 1
		} else if split == nil || len(split) != 2 {
			err = errors.Errorf("Error parsing scale ratio: \"%s\", line %d", line, lineno)
			return
		} else {
			if v, err = strconv.ParseInt(strings.TrimSpace(split[0]), 10, 32); err != nil {
				err = errors.Errorf("Error parsing scale ratio numerator: \"%s\", line %d", split[0], lineno)
				return
			}
			tone.RatioN = int(v)
			if v, err = strconv.ParseInt(strings.TrimSpace(split[1]), 10, 32); err != nil {
				err = errors.Errorf("Error parsing scale ratio numerator: \"%s\", line %d", split[1], lineno)
				return
			}
			tone.RatioD = int(v)
		}
		if tone.RatioD == 0 || tone.RatioN == 0 {
			err = errors.Errorf("Error parsing scale ratio - numerator or denominator is zero: \"%s\", line %d", line, lineno)
			return
		}

		// 2^(cents/1200) = n/d
		// cents = 1200 * log(n/d) / log(2)
		tone.Cents = 1200 * math.Log(float64(tone.RatioN)/float64(tone.RatioD)) / math.Log(2.0);
	}
	tone.FloatValue = (tone.Cents / 1200.0) + 1.0
	return
}

// ReadSCLStream returns a Scale from the SCL input stream
func ReadSCLStream(rdr io.Reader) (scale Scale, err error) {
	type stateType int
	const (
		readHeader stateType = iota
		readCount
		readNote
		trailing
	)
	state := readHeader
	scanner := bufio.NewScanner(rdr)
	lineno := 0
	for scanner.Scan() {
		line := scanner.Text()
		scale.RawText = scale.RawText + "\n" + line
		lineno++
		line = strings.TrimRight(line, "\t ")

		//fmt.Printf("DEBUG: l:%d state:%d line:\"%s\":  %#v",lineno,state,line,scale)

		if (state == readNote && len(line) == 0) ||  (len(line)>0 && line[0] == '!') {
			continue
		}
		var v int64
		switch state {
		case readHeader:
			scale.Description = line
			state = readCount
		case readCount:
			if v, err = strconv.ParseInt(strings.TrimSpace(line),10,32); err != nil  {
				err = errors.Wrapf(err, "Error parsing Count: \"%s\", line %d", line,lineno)
				return
			}
			if v < 1 {
				err = errors.Wrapf(err, "Error parsing Count: must be > 0: \"%s\", line %d", line,lineno)
				return
			}
			scale.Count = int(v)
			state = readNote
		case readNote:
			var tone Tone
			if tone,err = toneFromString(line, lineno); err != nil {
				return
			}
			scale.Tones = append(scale.Tones, tone)
			if len(scale.Tones) == scale.Count {
				state = trailing
				break
			}
		}
	}
	if err = scanner.Err(); err != nil {
		return
	}
	if ! (state == readNote || state == trailing) {
		err = errors.Errorf("Incomplete SCL file. Found no notes section in the file.")
		return
	}
	if len(scale.Tones) != scale.Count {
		err = errors.Errorf("Read fewer notes (%d) than count (%d)", len(scale.Tones),scale.Count)
		return
	}
	return
}

// ReadSCLFile returns a Scale from the SCL File in fname
func ReadSCLFile(fname string) (scale Scale, err error) {
	var file *os.File
	if file, err = os.Open(fname); err != nil {
		return
	}
	defer file.Close()
	if scale,err = ReadSCLStream(file); err != nil {
		return
	}
	scale.Name = fname
	return
}

// ParseSCLData returns a scale from the SCL file contents in memory
func ParseSCLData(sclContents string) (scale Scale, err error) {
	rdr := strings.NewReader(sclContents)
	if scale,err = ReadSCLStream(rdr); err != nil {
		return
	}
	scale.Name = "Scale from patch"
	return
}

// evenTemperament12NoteScale provides a utility scale which is
// the "standard tuning" scale
func evenTemperment12NoteScale() (scale Scale, err error) {
	if scale,err = ParseSCLData(`! 12 Tone Equal Temperament.scl
!
12 Tone Equal Temperament | ED2-12 - Equal division of harmonic 2 into 12 parts
 12
!
 100.00000
 200.00000
 300.00000
 400.00000
 500.00000
 600.00000
 700.00000
 800.00000
 900.00000
 1000.00000
 1100.00000
 2/1
`); err != nil {
		return
	}
	return
}

// evenDivisionOfSpanByM provides a scale referred to as "ED2-17" or
// "ED3-24" by dividing the Span into M points. eventDivisionOfSpanByM(2,12)
// should be the evenTemperament12NoteScale
func evenDivisionOfSpanByM(span int, m int) (scale Scale, err error) {
	if span <= 0 {
		err = errors.Errorf("Span must be a positive number: %d", span)
		return
	}
	if m <= 0 {
		err = errors.Errorf("You must divide the period into at least one step: M must be a positive number: %d", m)
		return
	}
	buf := "! Automatically generated ED " + string(span) + "-" + string(m) + " scale\n"
	buf += string(m) + "\n"
	buf += "!\n"

	topCents := 1200.0 * math.Log(float64(span)) / math.Log(2.0)
	dCents := topCents / float64(m)
	for i := 0; i < m; i++ {
		buf += fmt.Sprintf("%f\n", dCents * float64(i))
	}
	buf += string(span) + "/1\n"

	if scale,err = ParseSCLData(buf); err != nil {
		return
	}

	return
}

