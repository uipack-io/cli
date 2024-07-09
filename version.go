package uipack

import "bufio"

type Version struct {
	Major uint16
	Minor uint16
}

func (v *Version) Encode(r *bufio.Writer) error {
	err := writeUint16(r, v.Major)
	if err != nil {
		return err
	}
	err = writeUint16(r, v.Minor)
	if err != nil {
		return err
	}
	return nil
}

func (v *Version) Decode(r *bufio.Reader) error {
	var err error
	v.Major, err = readUint16(r)
	if err != nil {
		return err
	}
	v.Minor, err = readUint16(r)
	if err != nil {
		return err
	}
	return nil
}
