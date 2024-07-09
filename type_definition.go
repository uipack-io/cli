package uipack

import (
	"bufio"
)

type CustomTypeIdentifier uint32

type TypeDefinition struct {
	Identifier CustomTypeIdentifier
	Name       string
	Properties []PropertyDefinition
}

type PropertyDefinition struct {
	Name         string
	Type         ValueType
	DefaultValue interface{}
}

// Binary encoding

func (s *TypeDefinition) Decode(r *bufio.Reader, metadata *BundleMetadata) error {
	id, err := readUint64(r)
	if err != nil {
		return err
	}
	s.Identifier = CustomTypeIdentifier(id)
	s.Name, err = readString(r)
	if err != nil {
		return err
	}
	propertiesCount, err := readUint64(r)
	if err != nil {
		return err
	}
	s.Properties = make([]PropertyDefinition, propertiesCount)
	for i := range s.Properties {
		t := ValueType{}
		err = t.Decode(r)
		if err != nil {
			return err
		}
		s.Properties[i].Type = t
		s.Properties[i].Name, err = readString(r)
		if err != nil {
			return err
		}
		s.Properties[i].DefaultValue, err = metadata.DecodeArgument(t, r)
	}
	return nil
}

func (s *TypeDefinition) Encode(w *bufio.Writer, metadata *BundleMetadata) error {
	err := writeUint64(w, uint64(s.Identifier))
	if err != nil {
		return err
	}
	err = writeString(w, s.Name)
	if err != nil {
		return err
	}
	err = writeUint64(w, uint64(len(s.Properties)))
	if err != nil {
		return err
	}
	for i := range s.Properties {
		prop := s.Properties[i]
		err = prop.Type.Encode(w)
		if err != nil {
			return err
		}
		err = writeString(w, prop.Name)
		if err != nil {
			return err
		}
		err = metadata.EncodeArgument(w, prop.DefaultValue)
		if err != nil {
			return err
		}
	}
	return nil
}
