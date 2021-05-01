package scala

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
