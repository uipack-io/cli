package uipack

import "bufio"

type Label struct {
	Segments  []LabelSegment
	Variables []LabelVariable
}

func ParseLabel(data string) Label {
	return Label{} // TODO
}

type LabelVariable struct {
	Id           uint8             // The identifier of the variable.
	VariableType LabelVariableType // The type of the variable.
}

type LabelSegment struct {
	Text     *string
	Style    *LabelStyle
	Link     *string
	Variable *uint8 // This is a segment that is replaced by a value at runtime.
	Child    *Label
}

type LabelVariableType uint8

const (
	String LabelVariableType = iota
	Number
)

type LabelStyle uint8

const (
	None LabelStyle = iota
	Bold
	Italic
	Underlined
)

// Binary encoding

func (o *Label) Encode(writer *bufio.Writer) error {
	err := writeUint32(writer, uint32(len(o.Segments)))
	if err != nil {
		return err
	}
	for _, segment := range o.Segments {
		err = segment.Encode(writer)
		if err != nil {
			return err
		}
	}
	err = writeUint32(writer, uint32(len(o.Variables)))
	if err != nil {
		return err
	}
	for _, variable := range o.Variables {
		err = variable.Encode(writer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *Label) Decode(reader *bufio.Reader) error {
	count, err := readUint32(reader)
	if err != nil {
		return err
	}
	o.Segments = make([]LabelSegment, count)
	for i := 0; i < int(count); i++ {
		err = o.Segments[i].Decode(reader)
		if err != nil {
			return err
		}
	}
	count, err = readUint32(reader)
	if err != nil {
		return err
	}
	o.Variables = make([]LabelVariable, count)
	for i := 0; i < int(count); i++ {
		err = o.Variables[i].Decode(reader)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *LabelVariable) Encode(writer *bufio.Writer) error {
	err := writeUint8(writer, o.Id)
	if err != nil {
		return err
	}
	err = writeUint8(writer, uint8(o.VariableType))
	if err != nil {
		return err
	}
	return nil
}

func (o *LabelVariable) Decode(reader *bufio.Reader) error {
	var err error
	o.Id, err = readUint8(reader)
	if err != nil {
		return err
	}
	t, err := readUint8(reader)
	if err != nil {
		return err
	}
	o.VariableType = LabelVariableType(t)
	return nil
}

func (o *LabelSegment) Encode(writer *bufio.Writer) error {
	var err error
	if o.Text != nil {
		err = writeUint8(writer, 1)
		if err != nil {
			return err
		}
		err = writeString(writer, *o.Text)
		if err != nil {
			return err
		}
	} else {
		err = writeUint8(writer, 0)
		if err != nil {
			return err
		}
	}
	if o.Style != nil {
		err = writeUint8(writer, 1)
		if err != nil {
			return err
		}
		err = writeUint8(writer, uint8(*o.Style))
		if err != nil {
			return err
		}
	} else {
		err = writeUint8(writer, 0)
		if err != nil {
			return err
		}
	}
	if o.Link != nil {
		err = writeUint8(writer, 1)
		if err != nil {
			return err
		}
		err = writeString(writer, *o.Link)
		if err != nil {
			return err
		}
	} else {
		err = writeUint8(writer, 0)
		if err != nil {
			return err
		}
	}
	if o.Variable != nil {
		err = writeUint8(writer, 1)
		if err != nil {
			return err
		}
		err = writeUint8(writer, *o.Variable)
		if err != nil {
			return err
		}
	} else {
		err = writeUint8(writer, 0)
		if err != nil {
			return err
		}
	}
	if o.Child != nil {
		err = writeUint8(writer, 1)
		if err != nil {
			return err
		}
		err = o.Child.Encode(writer)
		if err != nil {
			return err
		}
	} else {
		err = writeUint8(writer, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *LabelSegment) Decode(reader *bufio.Reader) error {

	t, err := readUint8(reader)
	if err != nil {
		return err
	}
	if t == 1 {
		text, err := readString(reader)
		if err != nil {
			return err
		}
		o.Text = &text
	}

	t, err = readUint8(reader)
	if err != nil {
		return err
	}
	if t == 1 {
		s, err := readUint8(reader)
		if err != nil {
			return err
		}
		style := LabelStyle(s)
		o.Style = &style
	}
	t, err = readUint8(reader)
	if err != nil {
		return err
	}
	if t == 1 {
		link, err := readString(reader)
		if err != nil {
			return err
		}
		o.Link = &link
	}
	t, err = readUint8(reader)
	if err != nil {
		return err
	}
	if t == 1 {
		variable, err := readUint8(reader)
		if err != nil {
			return err
		}
		o.Variable = &variable
	}
	t, err = readUint8(reader)
	if err != nil {
		return err
	}
	if t == 1 {
		child := Label{}
		err = child.Decode(reader)
		if err != nil {
			return err
		}
		o.Child = &child
	}
	return nil
}
