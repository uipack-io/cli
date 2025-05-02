package test

import (
	"bufio"
	"bytes"
	"testing"

	uipack "github.com/uipack-io/cli"
)

func TestColorEncodeDecode(t *testing.T) {
	input_color := uipack.Color{
		Red:   0.5,
		Green: 0.6,
		Blue:  0.2,
		Alpha: 1.0,
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
