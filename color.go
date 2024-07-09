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

func (color *Color) Encode(writer *bufio.Writer) error {
	err := writer.WriteByte(byte(color.Red))
	if err != nil {
		return err
	}
	err = writer.WriteByte(byte(color.Green))
	if err != nil {
		return err
	}
	err = writer.WriteByte(byte(color.Blue))
	if err != nil {
		return err
	}
	err = writer.WriteByte(byte(color.Alpha))
	if err != nil {
		return err
	}

	return nil
}

func (color *Color) Decode(reader *bufio.Reader) error {
	var err error
	color.Red, err = readUint8(reader)
	if err != nil {
		return err
	}
	color.Green, err = readUint8(reader)
	if err != nil {
		return err
	}
	color.Blue, err = readUint8(reader)
	if err != nil {
		return err
	}
	color.Alpha, err = readUint8(reader)
	if err != nil {
		return err
	}
	return nil
}
