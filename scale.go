package scala

import "io"


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


// ReadSCLStream returns a Scale from the SCL input stream
func ReadSCLStream(rdr io.Reader) (s Scale, err error) {
	return
}

// ReadSCLFile returns a Scale from the SCL File in fname
func ReadSCLFile(fname string) (s Scale, err error) {
	return
}

// ParseSCLData returns a scale from the SCL file contents in memory
func ParseSCLData(sclContents string) (s Scale, err error) {
	return
}

// EvenTemperament12NoteScale provides a utility scale which is
// the "standard tuning" scale
func EvenTemperament12NoteScale() (s Scale, err error) {
	return
}

// EvenDivisionOfSpanByM provides a scale referred to as "ED2-17" or
// "ED3-24" by dividing the Span into M points. eventDivisionOfSpanByM(2,12)
// should be the evenTemperament12NoteScale
func EvenDivisionOfSpanByM(span int, m int) (s Scale, err error) {
	return
}
