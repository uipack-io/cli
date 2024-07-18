package test

import (
	"bufio"
	"bytes"
	"testing"

	uipack "github.com/uipack-io/cli"
)

func TestBundleEncodeDecode(t *testing.T) {
	metadata := uipack.BundleMetadata{
		Version: uipack.Version{
			Major: 4,
			Minor: 2,
		},
		Name: "Test Bundle",
		Modes: []uipack.ModeMetadata{
			{
				Name: "theme",
				Variants: []uipack.ModeVariantMetadata{
					{
						Identifier: 0,
						Name:       "dark",
					},
					{
						Identifier: 1,
						Name:       "light",
					},
				},
			},
		},
		Variables: []uipack.VariableMetadata{
			{
				Identifier: 0,
				Name:       "color",
				Type: uipack.ValueType{
					Type: uipack.ColorType,
				},
			},
			{
				Identifier: 1,
				Name:       "text_style",
				Type: uipack.ValueType{
					Type: uipack.TextStyleType,
				},
			},
		},
	}

	input_bundle := uipack.Bundle{
		Version: uipack.Version{
			Major: 4,
			Minor: 2,
		},
		Variant: 0x654,
		Values: []interface{}{
			uipack.Color{
				Red:   32,
				Green: 64,
				Blue:  128,
				Alpha: 255,
			},
			uipack.TextStyle{
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
			},
		},
	}

	buffer := bytes.NewBuffer(nil)
	w := bufio.NewWriter(buffer)
	writer := bufio.NewWriter(w)
	input_bundle.Encode(writer, &metadata)
	writer.Flush()

	r := bytes.NewReader(buffer.Bytes())
	reader := bufio.NewReader(r)

	output_bundle := uipack.Bundle{}
	output_bundle.Decode(reader, &metadata)

	if output_bundle.Version != input_bundle.Version {
		t.Error("Expected", input_bundle.Version, "got", output_bundle.Version)
	}

	if output_bundle.Variant != input_bundle.Variant {
		t.Error("Expected", input_bundle.Variant, "got", output_bundle.Variant)
	}

	if len(input_bundle.Values) != len(output_bundle.Values) {
		t.Error("Expected", len(input_bundle.Values), "got", len(output_bundle.Values))
	}

	for i, value := range input_bundle.Values {
		switch v := value.(type) {
		case uipack.Color:
			if v != output_bundle.Values[i].(uipack.Color) {
				t.Error("Expected", v, "got", output_bundle.Values[i].(uipack.Color))
			}
		case uipack.TextStyle:
			ov := output_bundle.Values[i].(uipack.TextStyle)
			if !v.Equal(&ov) {
				t.Error("Expected", v, "got", output_bundle.Values[i].(uipack.TextStyle))
			}
		}
	}

}
