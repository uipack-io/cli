package uipack

import "bufio"

type Version struct {
	Major uint16
	Minor uint16
}

func (v *Version) Encode(r *bufio.Writer) {
	writeUint16(r, v.Major)
	writeUint16(r, v.Minor)
}

func (v *Version) Decode(r *bufio.Reader) {
	v.Major = readUint16(r)
	v.Minor = readUint16(r)
}
