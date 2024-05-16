package uipack

import "bufio"

type Instance struct {
	Type      TypeIdentifier
	Arguments []interface{}
}

func (s *Instance) Decode(r *bufio.Reader) {
	s.Type = TypeIdentifier(readUint64(r))
	argumentsCount := readUint64(r)
	s.Arguments = make([]interface{}, argumentsCount)
	for i := range s.Arguments {
		t := VariableType(readUint8(r))
		s.Arguments[i] = DecodeArgument(t, r)
	}
}

func (s *Instance) Encode(w *bufio.Writer) {
	writeUint64(w, uint64(s.Type))
	writeUint64(w, uint64(len(s.Arguments)))
	for i := range s.Arguments {
		a := s.Arguments[i]
		t := FindVariableType(a)
		writeUint8(w, uint8(t))
		EncodeArgument(w, a)
	}
}
