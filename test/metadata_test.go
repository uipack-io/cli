package test

import (
	"aloisdeniel/uipack"
	"bufio"
	"bytes"
	"testing"
)

func TestMetadataEncodeDecode(t *testing.T) {
	input_metadata := uipack.BundleMetadata{
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
				Type:       uipack.ColorType,
			},
			{
				Identifier: 1,
				Name:       "text_style",
				Type:       uipack.TextStyleType,
			},
		},
	}

	buffer := bytes.NewBuffer(nil)
	w := bufio.NewWriter(buffer)
	writer := bufio.NewWriter(w)
	input_metadata.Encode(writer)
	writer.Flush()

	r := bytes.NewReader(buffer.Bytes())
	reader := bufio.NewReader(r)

	output_metadata := uipack.BundleMetadata{}
	output_metadata.Decode(reader)

	if input_metadata.Version != output_metadata.Version {
		t.Error("Expected", input_metadata.Version, "got", output_metadata.Version)
	}

	if input_metadata.Name != output_metadata.Name {
		t.Error("Expected", input_metadata.Name, "got", output_metadata.Name)
	}

	if len(input_metadata.Modes) != len(output_metadata.Modes) {
		t.Error("Expected", len(input_metadata.Modes), "got", len(output_metadata.Modes))
	}

	for i, mode := range input_metadata.Modes {
		if mode.Name != output_metadata.Modes[i].Name {
			t.Error("Expected", mode.Name, "got", output_metadata.Modes[i].Name)
		}

		if len(mode.Variants) != len(output_metadata.Modes[i].Variants) {
			t.Error("Expected", len(mode.Variants), "got", len(output_metadata.Modes[i].Variants))
		}

		for j, variant := range mode.Variants {
			if variant.Name != output_metadata.Modes[i].Variants[j].Name {
				t.Error("Expected", variant.Name, "got", output_metadata.Modes[i].Variants[j].Name)
			}
		}
	}

	if len(input_metadata.Variables) != len(output_metadata.Variables) {
		t.Error("Expected", len(input_metadata.Variables), "got", len(output_metadata.Variables))
	}

	for i, variable := range input_metadata.Variables {
		if variable.Name != output_metadata.Variables[i].Name {
			t.Error("Expected", variable.Name, "got", output_metadata.Variables[i].Name)
		}

		if variable.Type != output_metadata.Variables[i].Type {
			t.Error("Expected", variable.Type, "got", output_metadata.Variables[i].Type)
		}
	}
}
