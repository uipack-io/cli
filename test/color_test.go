package test

import (
	"bufio"
	"bytes"
	"testing"

	uipack "github.com/uipack-io/cli"
)

func TestColorEncodeDecode(t *testing.T) {
	input_color := uipack.Color{
		Red:   32,
		Green: 64,
		Blue:  128,
		Alpha: 255,
	}

	buffer := bytes.NewBuffer(nil)
	w := bufio.NewWriter(buffer)
	writer := bufio.NewWriter(w)
	input_color.Encode(writer)
	writer.Flush()

	r := bytes.NewReader(buffer.Bytes())
	reader := bufio.NewReader(r)

	output_color := uipack.Color{}
	output_color.Decode(reader)

	if input_color != output_color {
		t.Error("Expected", input_color, "got", output_color)
	}
}
