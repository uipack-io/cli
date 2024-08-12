package uipack

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
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

func (color *Color) ParseHexString(v string) {
	// Parsing '#RRGGBB' or '#RRGGBBAA'
	v = strings.TrimPrefix(v, "#")
	switch len(v) {
	case 8:
		values, _ := strconv.ParseUint(string(v), 16, 32)
		color.Alpha = uint8(values & 0xFF)
		color.Red = uint8((values >> 24) & 0xFF)
		color.Green = uint8((values >> 16) & 0xFF)
		color.Blue = uint8((values >> 8) & 0xFF)
	case 6:
		values, _ := strconv.ParseUint(string(v), 16, 24)
		color.Alpha = 255
		color.Red = uint8(values >> 16)
		color.Green = uint8((values >> 8) & 0xFF)
		color.Blue = uint8(values & 0xFF)
	default:
		panic("Invalid color format")
	}
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
