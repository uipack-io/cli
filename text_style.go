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

func (v *TextStyle) Encode(writer *bufio.Writer) error {
	err := writeString(writer, v.FontFamily)
	if err != nil {
		return err
	}
	err = writeFloat64(writer, v.FontSize)
	if err != nil {
		return err
	}
	err = writeUint8(writer, v.FontWeight)
	if err != nil {
		return err
	}
	err = writeFloat64(writer, v.LetterSpacing)
	if err != nil {
		return err
	}
	err = writeFloat64(writer, v.WordSpacing)
	if err != nil {
		return err
	}
	err = writeFloat64(writer, v.LineHeight)
	if err != nil {
		return err
	}

	err = writeUint32(writer, uint32(len(v.FontVariations)))
	if err != nil {
		return err
	}
	for _, variation := range v.FontVariations {
		err = variation.Encode(writer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *TextStyle) Decode(reader *bufio.Reader) error {
	var err error
	v.FontFamily, err = readString(reader)
	if err != nil {
		return err
	}
	v.FontSize, err = readFloat64(reader)
	if err != nil {
		return err
	}
	v.FontWeight, err = readUint8(reader)
	if err != nil {
		return err
	}
	v.LetterSpacing, err = readFloat64(reader)
	if err != nil {
		return err
	}
	v.WordSpacing, err = readFloat64(reader)
	if err != nil {
		return err
	}
	v.LineHeight, err = readFloat64(reader)
	if err != nil {
		return err
	}

	variationsCount, err := readUint32(reader)
	if err != nil {
		return err
	}
	variations := make([]FontVariation, variationsCount)
	for i := uint32(0); i < variationsCount; i++ {
		variation := FontVariation{}
		err = variation.Decode(reader)
		variations[i] = variation
		if err != nil {
			return err
		}
	}
	v.FontVariations = variations
	return nil
}

func (v *FontVariation) Encode(writer *bufio.Writer) error {
	err := writeString(writer, v.Axis)
	if err != nil {
		return err
	}
	err = writeFloat64(writer, v.Value)
	if err != nil {
		return err
	}
	return nil
}

func (v *FontVariation) Decode(reader *bufio.Reader) error {
	var err error
	v.Axis, err = readString(reader)
	if err != nil {
		return err
	}
	v.Value, err = readFloat64(reader)
	if err != nil {
		return err
	}

	return nil
}
