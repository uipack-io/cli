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

func (o *Label) Encode(writer *bufio.Writer) {
	writeUint32(writer, uint32(len(o.Segments)))
	for _, segment := range o.Segments {
		segment.Encode(writer)
	}
	writeUint32(writer, uint32(len(o.Variables)))
	for _, variable := range o.Variables {
		variable.Encode(writer)
	}
}

func (o *Label) Decode(reader *bufio.Reader) {
	count := readUint32(reader)
	o.Segments = make([]LabelSegment, count)
	for i := 0; i < int(count); i++ {
		o.Segments[i].Decode(reader)
	}
	count = readUint32(reader)
	o.Variables = make([]LabelVariable, count)
	for i := 0; i < int(count); i++ {
		o.Variables[i].Decode(reader)
	}
}

func (o *LabelVariable) Encode(writer *bufio.Writer) {
	writeUint8(writer, o.Id)
	writeUint8(writer, uint8(o.VariableType))
}

func (o *LabelVariable) Decode(reader *bufio.Reader) {
	o.Id = readUint8(reader)
	o.VariableType = LabelVariableType(readUint8(reader))
}

func (o *LabelSegment) Encode(writer *bufio.Writer) {
	if o.Text != nil {
		writeUint8(writer, 1)
		writeString(writer, *o.Text)
	} else {
		writeUint8(writer, 0)
	}
	if o.Style != nil {
		writeUint8(writer, 1)
		writeUint8(writer, uint8(*o.Style))
	} else {
		writeUint8(writer, 0)
	}
	if o.Link != nil {
		writeUint8(writer, 1)
		writeString(writer, *o.Link)
	} else {
		writeUint8(writer, 0)
	}
	if o.Variable != nil {
		writeUint8(writer, 1)
		writeUint8(writer, *o.Variable)
	} else {
		writeUint8(writer, 0)
	}
	if o.Child != nil {
		writeUint8(writer, 1)
		o.Child.Encode(writer)
	} else {
		writeUint8(writer, 0)
	}
}

func (o *LabelSegment) Decode(reader *bufio.Reader) {
	if readUint8(reader) == 1 {
		text := readString(reader)
		o.Text = &text
	}
	if readUint8(reader) == 1 {
		style := LabelStyle(readUint8(reader))
		o.Style = &style
	}
	if readUint8(reader) == 1 {
		link := readString(reader)
		o.Link = &link
	}
	if readUint8(reader) == 1 {
		variable := readUint8(reader)
		o.Variable = &variable
	}
	if readUint8(reader) == 1 {
		child := Label{}
		child.Decode(reader)
		o.Child = &child
	}
}
