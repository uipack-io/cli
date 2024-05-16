package uipack

import "bufio"

type Radius struct {
	X float64
	Y float64
}

type BorderRadius struct {
	TopLeft     Radius
	TopRight    Radius
	BottomRight Radius
	BottomLeft  Radius
}

// Binary encoding

func (r *Radius) Encode(writer *bufio.Writer) {
	writeFloat64(writer, r.X)
	writeFloat64(writer, r.Y)
}

func (r *Radius) Decode(reader *bufio.Reader) {
	r.X = readFloat64(reader)
	r.Y = readFloat64(reader)
}

func (b *BorderRadius) Encode(writer *bufio.Writer) {
	b.TopLeft.Encode(writer)
	b.TopRight.Encode(writer)
	b.BottomRight.Encode(writer)
	b.BottomLeft.Encode(writer)
}

func (b *BorderRadius) Decode(reader *bufio.Reader) {
	b.TopLeft.Decode(reader)
	b.TopRight.Decode(reader)
	b.BottomRight.Decode(reader)
	b.BottomLeft.Decode(reader)
}
