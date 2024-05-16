package uipack

import "bufio"

type TypeIdentifier uint64

type TypeDefinition struct {
	Identifier TypeIdentifier
	Name       string
	Properties []Property
}

type PropertyIdentifier uint64

type Property interface {
	Type() VariableType
}

type StringProperty struct {
	Name         string
	DefaultValue string
}

func (s *StringProperty) Type() VariableType {
	return StringType
}

type IntegerProperty struct {
	Name         string
	DefaultValue int64
}

func (s *IntegerProperty) Type() VariableType {
	return IntegerType
}

type BooleanProperty struct {
	Name         string
	DefaultValue bool
}

func (s *BooleanProperty) Type() VariableType {
	return BooleanType
}

type FloatProperty struct {
	Name         string
	DefaultValue float64
}

func (s *FloatProperty) Type() VariableType {
	return FloatType
}

type InstanceProperty struct {
	Name         string
	InstanceType TypeIdentifier
	DefaultValue Instance
}

func (s *InstanceProperty) Type() VariableType {
	return InstanceType
}

// Binary encoding

func (s *TypeDefinition) Decode(r *bufio.Reader) {
	s.Identifier = TypeIdentifier(readUint64(r))
	s.Name = readString(r)
	propertiesCount := readUint64(r)
	s.Properties = make([]Property, propertiesCount)
	for i := range s.Properties {
		switch VariableType(readUint8(r)) {
		case StringType:
			result := StringProperty{}
			result.Decode(r)
			s.Properties[i] = &result
		case IntegerType:
			result := IntegerProperty{}
			result.Decode(r)
			s.Properties[i] = &result
		case BooleanType:
			result := BooleanProperty{}
			result.Decode(r)
			s.Properties[i] = &result
		case FloatType:
			result := FloatProperty{}
			result.Decode(r)
			s.Properties[i] = &result
		case InstanceType:
			result := InstanceProperty{}
			result.Decode(r)
			s.Properties[i] = &result
		default:
			panic("Unknown property type")
		}
	}
}

func (s *TypeDefinition) Encode(w *bufio.Writer) {
	writeUint64(w, uint64(s.Identifier))
	writeString(w, s.Name)
	writeUint64(w, uint64(len(s.Properties)))
	for i := range s.Properties {
		prop := s.Properties[i]
		writeUint8(w, uint8(prop.Type()))
		switch prop.Type() {
		case StringType:
			prop.(*StringProperty).Encode(w)
		case IntegerType:
			prop.(*IntegerProperty).Encode(w)
		case BooleanType:
			prop.(*BooleanProperty).Encode(w)
		case FloatType:
			prop.(*FloatProperty).Encode(w)
		case InstanceType:
			prop.(*InstanceProperty).Encode(w)
		default:
			panic("Unknown property type")
		}
	}
}

func (s *FloatProperty) Decode(r *bufio.Reader) {
	s.Name = readString(r)
	s.DefaultValue = readFloat64(r)
}

func (s *FloatProperty) Encode(w *bufio.Writer) {
	writeString(w, s.Name)
	writeFloat64(w, s.DefaultValue)
}

func (s *BooleanProperty) Decode(r *bufio.Reader) {
	s.Name = readString(r)
	s.DefaultValue = readBool(r)
}

func (s *BooleanProperty) Encode(w *bufio.Writer) {
	writeString(w, s.Name)
	writeBool(w, s.DefaultValue)
}

func (s *IntegerProperty) Decode(r *bufio.Reader) {
	s.Name = readString(r)
	s.DefaultValue = readInt64(r)
}

func (s *IntegerProperty) Encode(w *bufio.Writer) {
	writeString(w, s.Name)
	writeInt64(w, s.DefaultValue)
}

func (s *StringProperty) Decode(r *bufio.Reader) {
	s.Name = readString(r)
	s.DefaultValue = readString(r)
}

func (s *StringProperty) Encode(w *bufio.Writer) {
	writeString(w, s.Name)
	writeString(w, s.DefaultValue)
}

func (s *InstanceProperty) Decode(r *bufio.Reader) {
	s.Name = readString(r)
	s.InstanceType = TypeIdentifier(readUint64(r))
	s.DefaultValue.Decode(r)
}

func (s *InstanceProperty) Encode(w *bufio.Writer) {
	writeString(w, s.Name)
	writeUint64(w, uint64(s.InstanceType))
	s.DefaultValue.Encode(w)
}
