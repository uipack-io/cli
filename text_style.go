package uipack

import (
	"bufio"
)

// A text style.
type TextStyle struct {
	FontFamily     string
	FontSize       float64
	FontWeight     uint8
	LetterSpacing  float64
	WordSpacing    float64
	LineHeight     float64
	FontVariations []FontVariation
}

func (v *TextStyle) Equal(o *TextStyle) bool {

	if v.FontSize != o.FontSize {
		return false
	}

	if v.FontWeight != o.FontWeight {
		return false
	}

	if v.LetterSpacing != o.LetterSpacing {
		return false
	}

	if v.WordSpacing != o.WordSpacing {
		return false
	}

	if v.LineHeight != o.LineHeight {
		return false
	}

	if len(v.FontVariations) != len(o.FontVariations) {
		return false
	}

	for i, variation := range v.FontVariations {
		ov := o.FontVariations[i]
		if variation.Axis != ov.Axis {
			return false
		}

		if variation.Value != ov.Value {
			return false
		}
	}
	return true
}

// A variation for variable fonts
type FontVariation struct {
	Axis  string
	Value float64
}

// Binary encoding

func (v *TextStyle) Encode(writer *bufio.Writer) {
	writeString(writer, v.FontFamily)
	writeFloat64(writer, v.FontSize)
	writeUint8(writer, v.FontWeight)
	writeFloat64(writer, v.LetterSpacing)
	writeFloat64(writer, v.WordSpacing)
	writeFloat64(writer, v.LineHeight)

	writeUint32(writer, uint32(len(v.FontVariations)))
	for _, variation := range v.FontVariations {
		variation.Encode(writer)
	}
}

func (v *TextStyle) Decode(reader *bufio.Reader) {
	v.FontFamily = readString(reader)
	v.FontSize = readFloat64(reader)
	v.FontWeight = readUint8(reader)
	v.LetterSpacing = readFloat64(reader)
	v.WordSpacing = readFloat64(reader)
	v.LineHeight = readFloat64(reader)

	variationsCount := readUint32(reader)
	variations := make([]FontVariation, variationsCount)
	for i := uint32(0); i < variationsCount; i++ {
		variation := FontVariation{}
		variation.Decode(reader)
		variations[i] = variation
	}
	v.FontVariations = variations
}

func (v *FontVariation) Encode(writer *bufio.Writer) {
	writeString(writer, v.Axis)
	writeFloat64(writer, v.Value)
}

func (v *FontVariation) Decode(reader *bufio.Reader) {
	v.Axis = readString(reader)
	v.Value = readFloat64(reader)
}
