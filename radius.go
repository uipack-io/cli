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

func (r *Radius) Encode(writer *bufio.Writer) error {
	err := writeFloat64(writer, r.X)
	if err != nil {
		return err
	}
	err = writeFloat64(writer, r.Y)
	if err != nil {
		return err
	}
	return nil
}

func (r *Radius) Decode(reader *bufio.Reader) error {
	var err error
	r.X, err = readFloat64(reader)
	if err != nil {
		return err
	}
	r.Y, err = readFloat64(reader)
	if err != nil {
		return err
	}

	return nil
}

func (b *BorderRadius) Encode(writer *bufio.Writer) error {
	err := b.TopLeft.Encode(writer)
	if err != nil {
		return err
	}
	err = b.TopRight.Encode(writer)
	if err != nil {
		return err
	}
	err = b.BottomRight.Encode(writer)
	if err != nil {
		return err
	}
	err = b.BottomLeft.Encode(writer)
	if err != nil {
		return err
	}
	return nil
}

func (b *BorderRadius) Decode(reader *bufio.Reader) error {
	err := b.TopLeft.Decode(reader)
	if err != nil {
		return err
	}
	err = b.TopRight.Decode(reader)
	if err != nil {
		return err
	}
	err = b.BottomRight.Decode(reader)
	if err != nil {
		return err
	}
	err = b.BottomLeft.Decode(reader)
	if err != nil {
		return err
	}
	return nil
}
