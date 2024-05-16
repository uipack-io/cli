package uipack

import (
	"bufio"
	"encoding/binary"
	"math"
)

func writeUint8(writer *bufio.Writer, value uint8) {
	writeByte(writer, byte(value))
}

func readUint8(reader *bufio.Reader) uint8 {
	return uint8(readByte(reader))
}

func writeByte(writer *bufio.Writer, value byte) {
	err := writer.WriteByte(value)
	if err != nil {
		panic(err)
	}
}

func readByte(reader *bufio.Reader) byte {
	r, err := reader.ReadByte()
	if err != nil {
		panic(err)
	}
	return r
}

func writeUint16(writer *bufio.Writer, value uint16) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint16(b, value)
	_, err := writer.Write(b)
	if err != nil {
		panic(err)
	}
}

func readUint16(reader *bufio.Reader) uint16 {
	b := make([]byte, 4)
	_, err := reader.Read(b)
	if err != nil {
		panic(err)
	}
	return binary.BigEndian.Uint16(b)
}

func writeUint32(writer *bufio.Writer, value uint32) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, value)
	_, err := writer.Write(b)
	if err != nil {
		panic(err)
	}
}

func readUint32(reader *bufio.Reader) uint32 {
	b := make([]byte, 4)
	_, err := reader.Read(b)
	if err != nil {
		panic(err)
	}
	return binary.BigEndian.Uint32(b)
}

func writeInt64(writer *bufio.Writer, value int64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(value))
	_, err := writer.Write(b)
	if err != nil {
		panic(err)
	}
}

func readInt64(reader *bufio.Reader) int64 {
	b := make([]byte, 8)
	_, err := reader.Read(b)
	if err != nil {
		panic(err)
	}
	return int64(binary.BigEndian.Uint64(b))
}

func writeUint64(writer *bufio.Writer, value uint64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, value)
	_, err := writer.Write(b)
	if err != nil {
		panic(err)
	}
}

func readUint64(reader *bufio.Reader) uint64 {
	b := make([]byte, 8)
	_, err := reader.Read(b)
	if err != nil {
		panic(err)
	}
	return binary.BigEndian.Uint64(b)
}

func writeFloat64(writer *bufio.Writer, value float64) {
	bits := math.Float64bits(value)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, bits)
	_, err := writer.Write(b)
	if err != nil {
		panic(err)
	}
}

func readFloat64(reader *bufio.Reader) float64 {
	b := make([]byte, 8)
	_, err := reader.Read(b)
	if err != nil {
		panic(err)
	}
	bits := binary.LittleEndian.Uint64(b)
	return math.Float64frombits(bits)
}

// Writes a string to the writer, with its size first.
func writeString(writer *bufio.Writer, value string) {
	writeUint32(writer, uint32(len(value)))
	_, err := writer.WriteString(value)
	if err != nil {
		panic(err)
	}
}

// Reads a string from the reader, with its size first.
func readString(reader *bufio.Reader) string {
	length := readUint32(reader)
	b := make([]byte, length)
	_, err := reader.Read(b)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func writeBool(writer *bufio.Writer, value bool) {
	if value {
		writeByte(writer, 1)
	} else {
		writeByte(writer, 0)
	}
}

func readBool(reader *bufio.Reader) bool {
	return readByte(reader) == 1
}
