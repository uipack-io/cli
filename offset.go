package uipack

import "bufio"

type Offset struct {
	X float64
	Y float64
}

// Binary encoding

func (o *Offset) Encode(writer *bufio.Writer) {
	writeFloat64(writer, o.X)
	writeFloat64(writer, o.Y)
}

func (o *Offset) Decode(reader *bufio.Reader) {
	o.X = readFloat64(reader)
	o.Y = readFloat64(reader)
}
