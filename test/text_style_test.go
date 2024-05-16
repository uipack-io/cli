package test

import (
	"aloisdeniel/uipack"
	"bufio"
	"bytes"
	"testing"
)

func TestTextStyleEncodeDecode(t *testing.T) {
	input_text_style := uipack.TextStyle{
		FontFamily:    "Roboto",
		FontSize:      16,
		FontWeight:    4,
		LetterSpacing: 0.5,
		WordSpacing:   0.75,
		LineHeight:    1.5,
		FontVariations: []uipack.FontVariation{
			{
				Axis:  "wght",
				Value: 400,
			},
			{
				Axis:  "wdth",
				Value: 100,
			},
		},
	}

	buffer := bytes.NewBuffer(nil)
	w := bufio.NewWriter(buffer)
	writer := bufio.NewWriter(w)
	input_text_style.Encode(writer)
	writer.Flush()

	r := bytes.NewReader(buffer.Bytes())
	reader := bufio.NewReader(r)

	output_text_style := uipack.TextStyle{}
	output_text_style.Decode(reader)

	if input_text_style.FontFamily != output_text_style.FontFamily {
		t.Error("Expected", input_text_style.FontFamily, "got", output_text_style.FontFamily)
	}

	if input_text_style.FontSize != output_text_style.FontSize {
		t.Error("Expected", input_text_style.FontSize, "got", output_text_style.FontSize)
	}

	if input_text_style.FontWeight != output_text_style.FontWeight {
		t.Error("Expected", input_text_style.FontWeight, "got", output_text_style.FontWeight)
	}

	if input_text_style.LetterSpacing != output_text_style.LetterSpacing {
		t.Error("Expected", input_text_style.LetterSpacing, "got", output_text_style.LetterSpacing)
	}

	if input_text_style.WordSpacing != output_text_style.WordSpacing {
		t.Error("Expected", input_text_style.WordSpacing, "got", output_text_style.WordSpacing)
	}

	if input_text_style.LineHeight != output_text_style.LineHeight {
		t.Error("Expected", input_text_style.LineHeight, "got", output_text_style.LineHeight)
	}

	if len(input_text_style.FontVariations) != len(output_text_style.FontVariations) {
		t.Error("Expected", len(input_text_style.FontVariations), "got", len(output_text_style.FontVariations))
	}

	for i, variation := range input_text_style.FontVariations {
		if variation.Axis != output_text_style.FontVariations[i].Axis {
			t.Error("Expected", variation.Axis, "got", output_text_style.FontVariations[i].Axis)
		}

		if variation.Value != output_text_style.FontVariations[i].Value {
			t.Error("Expected", variation.Value, "got", output_text_style.FontVariations[i].Value)
		}
	}

}
