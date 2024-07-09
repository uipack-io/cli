package uipack

import "bufio"

type GradientStop struct {
	Offset float64
	Color  Color
}

type LinearGradient struct {
	Stops []GradientStop
	Begin Offset
	End   Offset
}

type RadialGradient struct {
	Stops  []GradientStop
	Center Offset
	Radius float64
}

// Binary encoding

func (o *GradientStop) Encode(writer *bufio.Writer) error {
	writeFloat64(writer, o.Offset)
	return o.Color.Encode(writer)
}

func (o *GradientStop) Decode(reader *bufio.Reader) error {
	var err error
	o.Offset, err = readFloat64(reader)
	if err != nil {
		return err
	}
	err = o.Color.Decode(reader)
	return err
}

func (o *LinearGradient) Encode(writer *bufio.Writer) error {
	err := writeUint32(writer, uint32(len(o.Stops)))
	if err != nil {
		return err
	}
	for _, stop := range o.Stops {
		err = stop.Encode(writer)
		if err != nil {
			return err
		}
	}
	err = o.Begin.Encode(writer)
	if err != nil {
		return err
	}
	err = o.End.Encode(writer)
	if err != nil {
		return err
	}
	return nil
}

func (o *LinearGradient) Decode(reader *bufio.Reader) error {
	count, err := readUint32(reader)
	if err != nil {
		return err
	}
	o.Stops = make([]GradientStop, count)
	for i := 0; i < int(count); i++ {
		err = o.Stops[i].Decode(reader)
		if err != nil {
			return err
		}
	}
	err = o.Begin.Decode(reader)
	if err != nil {
		return err
	}
	err = o.End.Decode(reader)
	return err
}

func (o *RadialGradient) Encode(writer *bufio.Writer) error {
	err := writeUint32(writer, uint32(len(o.Stops)))
	if err != nil {
		return err
	}
	for _, stop := range o.Stops {
		err = stop.Encode(writer)
		if err != nil {
			return err
		}
	}
	err = o.Center.Encode(writer)
	if err != nil {
		return err
	}
	err = writeFloat64(writer, o.Radius)
	if err != nil {
		return err
	}
	return nil
}

func (o *RadialGradient) Decode(reader *bufio.Reader) error {
	count, err := readUint32(reader)
	if err != nil {
		return err
	}
	o.Stops = make([]GradientStop, count)
	for i := 0; i < int(count); i++ {
		err = o.Stops[i].Decode(reader)
		if err != nil {
			return err
		}
	}
	err = o.Center.Decode(reader)
	if err != nil {
		return err
	}

	o.Radius, err = readFloat64(reader)
	if err != nil {
		return err
	}
	return nil
}
