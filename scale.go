package scala

import "io"

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
