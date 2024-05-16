package uipack

import (
	"bufio"
	"fmt"
)

// A color value.
type Color struct {
	Red, Green, Blue, Alpha uint8
}

var Black = Color{
	Red:   uint8(0),
	Green: uint8(0),
	Blue:  uint8(0),
	Alpha: uint8(255),
}

func (color *Color) ToHexString() string {
	return fmt.Sprintf("%02x%02x%02x%02x", color.Alpha, color.Red, color.Green, color.Blue)
}

// Binary encoding

func (color *Color) Encode(writer *bufio.Writer) {
	err := writer.WriteByte(byte(color.Red))
	if err != nil {
		panic(err)
	}
	err = writer.WriteByte(byte(color.Green))
	if err != nil {
		panic(err)
	}
	err = writer.WriteByte(byte(color.Blue))
	if err != nil {
		panic(err)
	}
	err = writer.WriteByte(byte(color.Alpha))
	if err != nil {
		panic(err)
	}
}

func (color *Color) Decode(reader *bufio.Reader) {
	color.Red = readUint8(reader)
	color.Green = readUint8(reader)
	color.Blue = readUint8(reader)
	color.Alpha = readUint8(reader)
}
