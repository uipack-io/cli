package uipack

import "bufio"

type Instance struct {
	Type      CustomTypeIdentifier
	Arguments []interface{}
}

func (s *Instance) Decode(r *bufio.Reader, metadata *BundleMetadata) error {
	t, err := readUint32(r)
	if err != nil {
		return err
	}
	s.Type = CustomTypeIdentifier(t)
	argumentsCount, err := readUint64(r)
	if err != nil {
		return err
	}
	typedef := metadata.FindTypeDefintion(s.Type)
	s.Arguments = make([]interface{}, argumentsCount)
	for i := range s.Arguments {
		argtype := typedef.Properties[i].Type
		s.Arguments[i], err = metadata.DecodeArgument(argtype, r)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Instance) Encode(w *bufio.Writer, metadata *BundleMetadata) error {
	err := writeUint32(w, uint32(s.Type))
	if err != nil {
		return err
	}
	err = writeUint64(w, uint64(len(s.Arguments)))
	if err != nil {
		return err
	}

	for i := range s.Arguments {
		a := s.Arguments[i]
		err = metadata.EncodeArgument(w, a)
		if err != nil {
			return err
		}
	}

	return nil
}
