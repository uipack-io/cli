package uipack

import (
	"bufio"
	"encoding/binary"
	"math"
)

func writeUint8(writer *bufio.Writer, value uint8) error {
	return writer.WriteByte(byte(value))
}

func readUint8(reader *bufio.Reader) (uint8, error) {
	result, err := reader.ReadByte()
	return uint8(result), err
}

func writeUint16(writer *bufio.Writer, value uint16) error {
	b := make([]byte, 4)
	binary.BigEndian.PutUint16(b, value)
	_, err := writer.Write(b)
	return err
}

func readUint16(reader *bufio.Reader) (uint16, error) {
	b := make([]byte, 4)
	_, err := reader.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(b), nil
}

func writeUint32(writer *bufio.Writer, value uint32) error {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, value)
	_, err := writer.Write(b)
	return err
}

func readUint32(reader *bufio.Reader) (uint32, error) {
	b := make([]byte, 4)
	_, err := reader.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(b), nil
}

func writeInt64(writer *bufio.Writer, value int64) error {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(value))
	_, err := writer.Write(b)
	return err
}

func readInt64(reader *bufio.Reader) (int64, error) {
	b := make([]byte, 8)
	_, err := reader.Read(b)
	if err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(b)), nil
}

func writeUint64(writer *bufio.Writer, value uint64) error {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, value)
	_, err := writer.Write(b)
	return err
}

func readUint64(reader *bufio.Reader) (uint64, error) {
	b := make([]byte, 8)
	_, err := reader.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(b), nil
}

func writeFloat64(writer *bufio.Writer, value float64) error {
	bits := math.Float64bits(value)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, bits)
	_, err := writer.Write(b)
	return err

}

func readFloat64(reader *bufio.Reader) (float64, error) {
	b := make([]byte, 8)
	_, err := reader.Read(b)
	if err != nil {
		return 0, err
	}
	bits := binary.LittleEndian.Uint64(b)
	return math.Float64frombits(bits), nil
}

// Writes a string to the writer, with its size first.
func writeString(writer *bufio.Writer, value string) error {
	writeUint32(writer, uint32(len(value)))
	_, err := writer.WriteString(value)
	if err != nil {
		return err
	}
	return nil
}

// Reads a string from the reader, with its size first.
func readString(reader *bufio.Reader) (string, error) {
	length, err := readUint32(reader)
	if err != nil {
		return "", err
	}
	b := make([]byte, length)
	_, err = reader.Read(b)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func writeBool(writer *bufio.Writer, value bool) error {
	if value {
		return writer.WriteByte(1)
	} else {
		return writer.WriteByte(0)
	}
}

func readBool(reader *bufio.Reader) (bool, error) {
	r, err := reader.ReadByte()
	return r == 1, err
}
