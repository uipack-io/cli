package uipack

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type ColorSpace uint8

const (
	SRGB = iota
	ExtendedSRGB
	DisplayP3
)

// A color value.
type Color struct {
	Red, Green, Blue, Alpha float64
	ColorSpace              ColorSpace
}

var Black = Color{
	Red:        0.0,
	Green:      0.0,
	Blue:       0.0,
	Alpha:      1.0,
	ColorSpace: DisplayP3,
}

func (color *Color) ParseHexString(v string) {
	// Parsing '#RRGGBB' or '#RRGGBBAA'
	v = strings.TrimPrefix(v, "#")
	switch len(v) {
	case 8:
		values, _ := strconv.ParseUint(string(v), 16, 32)
		color.Alpha = float64(uint8(values&0xFF)) / 255.0
		color.Red = float64(uint8((values>>24)&0xFF)) / 255.0
		color.Green = float64(uint8((values>>16)&0xFF)) / 255.0
		color.Blue = float64(uint8((values>>8)&0xFF)) / 255.0
	case 6:
		values, _ := strconv.ParseUint(string(v), 16, 24)
		color.Alpha = 1.0
		color.Red = float64(uint8(values>>16)) / 255.0
		color.Green = float64(uint8((values>>8)&0xFF)) / 255.0
		color.Blue = float64(uint8(values&0xFF)) / 255.0
	default:
		panic("Invalid color format")
	}
}

func (color *Color) ToHexString() string {
	toBytes := func(v float64) byte {
		return byte(uint32(v * 255.0))
	}
	return fmt.Sprintf("%02x%02x%02x%02x", toBytes(color.Alpha), toBytes(color.Red), toBytes(color.Green), toBytes(color.Blue))
}

// Binary encoding

func (color *Color) Encode(writer *bufio.Writer) error {
	err := writer.WriteByte(byte(color.ColorSpace))
	if err != nil {
		return err
	}
	err = writeFloat64(writer, color.Red)
	if err != nil {
		return err
	}
	err = writeFloat64(writer, color.Green)
	if err != nil {
		return err
	}
	err = writeFloat64(writer, color.Blue)
	if err != nil {
		return err
	}
	err = writeFloat64(writer, color.Alpha))
	if err != nil {
		return err
	}

	return nil
}

func (color *Color) Decode(reader *bufio.Reader) error {
	spaceByte, err := reader.ReadByte()
	if err != nil {
		return err
	}
	color.ColorSpace = ColorSpace(spaceByte)

	color.Red, err = readFloat64(reader)
	if err != nil {
		return err
	}
	color.Green, err = readFloat64(reader)
	if err != nil {
		return err
	}
	color.Blue, err = readFloat64(reader)
	if err != nil {
		return err
	}
	color.Alpha, err = readFloat64(reader)
	if err != nil {
		return err
	}
	return nil
}
