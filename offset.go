package uipack

import "bufio"

type Offset struct {
	X float64
	Y float64
}

// Binary encoding

func (o *Offset) Encode(writer *bufio.Writer) error {
	err := writeFloat64(writer, o.X)
	if err != nil {
		return err
	}
	err = writeFloat64(writer, o.Y)
	if err != nil {
		return err
	}
	return nil
}

func (o *Offset) Decode(reader *bufio.Reader) error {
	var err error
	o.X, err = readFloat64(reader)
	if err != nil {
		return err
	}
	o.Y, err = readFloat64(reader)
	if err != nil {
		return err
	}
	return nil
}
