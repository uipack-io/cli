package uipack

import "bufio"

type GradientStop struct {
	Offset float64
	Color  Color
}

type LinearGradient struct {
	Stops []GradientStop
	Begin Offset
	End   Offset
}

type RadialGradient struct {
	Stops  []GradientStop
	Center Offset
	Radius float64
}

// Binary encoding

func (o *GradientStop) Encode(writer *bufio.Writer) {
	writeFloat64(writer, o.Offset)
	o.Color.Encode(writer)
}

func (o *GradientStop) Decode(reader *bufio.Reader) {
	o.Offset = readFloat64(reader)
	o.Color.Decode(reader)
}

func (o *LinearGradient) Encode(writer *bufio.Writer) {
	writeUint32(writer, uint32(len(o.Stops)))
	for _, stop := range o.Stops {
		stop.Encode(writer)
	}
	o.Begin.Encode(writer)
	o.End.Encode(writer)
}

func (o *LinearGradient) Decode(reader *bufio.Reader) {
	count := readUint32(reader)
	o.Stops = make([]GradientStop, count)
	for i := 0; i < int(count); i++ {
		o.Stops[i].Decode(reader)
	}
	o.Begin.Decode(reader)
	o.End.Decode(reader)
}

func (o *RadialGradient) Encode(writer *bufio.Writer) {
	writeUint32(writer, uint32(len(o.Stops)))
	for _, stop := range o.Stops {
		stop.Encode(writer)
	}
	o.Center.Encode(writer)
	writeFloat64(writer, o.Radius)
}

func (o *RadialGradient) Decode(reader *bufio.Reader) {
	count := readUint32(reader)
	o.Stops = make([]GradientStop, count)
	for i := 0; i < int(count); i++ {
		o.Stops[i].Decode(reader)
	}
	o.Center.Decode(reader)
	o.Radius = readFloat64(reader)
}
