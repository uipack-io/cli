package importers

import (
	"encoding/json"
	"strings"

	uipack "github.com/uipack-io/cli"
)

// Support for variables exported with this plugin :
//
// Plugin: variables2json
// -> https://www.figma.com/community/plugin/1253571037276959291/variables2json

func DecodeJson(p *uipack.Package, jsn []byte) {
	export := jsonExport{}
	json.Unmarshal([]byte(jsn), &export)
	export.ToPackage(p)
}

// Raw JSON structure from the Figma plugin

type jsonExport struct {
	Collections []jsonCollection `json:"collections"`
}

func (f *jsonExport) ToPackage(p *uipack.Package) {
	p.Metadata = *f.ToBundleMetadata()

	combinations := p.Metadata.GenerateModeCombinations()

	for _, combination := range combinations {
		vid := uipack.Variant(0)
		variant := make(map[string]string)
		for mi, mv := range combination {
			mode := p.Metadata.Modes[mi]
			variant[mode.Name] = mv.Name
			vid = vid.SetMode(mode.Identifier, uipack.Uint4(mv.Identifier))
		}
		p.Bundles = append(p.Bundles, *f.ToBundle(vid, variant))
	}
}

func (f *jsonExport) ToBundleMetadata() *uipack.BundleMetadata {
	result := uipack.BundleMetadata{
		Name: "Figma",
		Version: uipack.Version{
			Major: 1,
			Minor: 0,
		},
	}
	vi := 0
	for ci, collection := range f.Collections {
		mm := uipack.ModeMetadata{
			Identifier: uipack.Uint4(ci),
			Name:       collection.Name,
		}
		for i, mode := range collection.Modes {
			mm.Variants = append(mm.Variants, uipack.ModeVariantMetadata{
				Identifier: uipack.Uint4(i),
				Name:       mode.Name,
			})
			if i == 0 {
				for _, variable := range mode.Variables {
					vm := uipack.VariableMetadata{
						Identifier: uint64(vi),
						Name:       variable.Name,
						Type:       figmaToVariableType(variable.Type),
					}
					result.Variables = append(result.Variables, vm)
					vi = vi + 1
				}
			}
		}

		result.Modes = append(result.Modes, mm)
	}
	return &result
}

func (f *jsonExport) ToBundle(identifier uipack.Variant, variant map[string]string) *uipack.Bundle {
	result := uipack.Bundle{
		Variant: identifier,
	}
	for _, collection := range f.Collections {
		mn := variant[collection.Name]
		mode := collection.FindMode(mn)
		if mode != nil {
			for _, variable := range mode.Variables {
				result.Values = append(result.Values, f.resolveVariable(collection.Name, variant, &variable))
			}
		}
	}
	return &result
}

func (f *jsonExport) FindCollection(name string) *jsonCollection {
	for _, collection := range f.Collections {
		if collection.Name == name {
			return &collection
		}
	}
	return nil
}

func (f *jsonExport) resolveVariable(currentCollection string, variant map[string]string, v *jsonVariable) interface{} {
	if v.IsAlias {
		alias := v.Alias(currentCollection)

		acol := f.FindCollection(alias.Collection)
		if acol != nil {
			mn := variant[alias.Collection]
			mode := acol.FindMode(mn)
			if mode != nil {
				variable := mode.FindVariable(alias.Name)
				if variable != nil {
					return f.resolveVariable(currentCollection, variant, variable)
				}

				panic("Alias variable not found " + alias.Name)
			}

			panic("Alias mode not found " + mn)
		}

		panic("Alias collection not found " + alias.Collection)
	}

	return figmaToVariableValue(v)
}

func figmaToVariableValue(v *jsonVariable) interface{} {
	if v.IsAlias {
		panic("Alias not resolved")
	}
	switch v.Type {
	case "color":
		switch v := v.Value.(type) {
		case string:
			result := uipack.Color{}
			result.ParseHexString(v)
			return result
		default:
			panic("Unknown type")
		}
	case "typography":
		switch v := v.Value.(type) {
		case map[string]interface{}:
			fontSize := v["fontSize"].(float64)
			return uipack.TextStyle{
				FontFamily:    v["fontFamily"].(string),
				FontSize:      fontSize,
				FontWeight:    figmaFontWeightToIndex(v["fontWeight"].(string)),
				LetterSpacing: figmaToLetterSpacing(v["letterSpacing"].(float64), fontSize, v["letterSpacingUnit"]),
				LineHeight:    figmaToLineHeight(v["lineHeight"].(float64), fontSize, v["lineHeightUnit"]),
			}
		default:
			panic("typography should be a map")
		}
	case "number":
		switch v := v.Value.(type) {
		case int:
			return float64(v)
		case float64:
			return float64(v)
		default:
			panic("number should be an int or a float")
		}

	case "string":
		switch v := v.Value.(type) {
		case string:
			return v
		default:
			panic("string should be a string")
		}
	default:
		panic("Unknown type " + v.Type)
	}
}

func figmaToLineHeight(v float64, fontSize float64, unit interface{}) float64 {
	switch unit.(string) {
	case "PERCENT":
		return v * 0.01 * fontSize
	case "PIXELS":
		return v
	default:
		return v
	}
}

func figmaToLetterSpacing(v float64, fontSize float64, unit interface{}) float64 {
	switch unit.(string) {
	case "PERCENT":
		return v * 0.01 * fontSize
	case "PIXELS":
		return v
	default:
		return v
	}
}

func figmaToVariableType(t string) uipack.ValueType {
	switch t {
	case "color":
		return uipack.ValueType{
			Type: uipack.ColorType,
		}
	case "typography":
		return uipack.ValueType{
			Type: uipack.TextStyleType,
		}
	case "number":
		return uipack.ValueType{
			Type: uipack.FloatType,
		}
	case "string":
		return uipack.ValueType{
			Type: uipack.StringType,
		}
	default:
		return uipack.ValueType{
			Type: uipack.TextStyleType,
		}
	}
}

func figmaFontWeightToIndex(v string) uint8 {
	switch strings.ToLower(v) {
	case "thin":
		return 1
	case "extralight":
		return 2
	case "light":
		return 3
	case "regular":
		return 4
	case "medium":
		return 5
	case "semibold":
		return 6
	case "bold":
		return 7
	case "extrabold":
		return 8
	case "black":
		return 9
	default:
		return 4
	}
}

type jsonCollection struct {
	Name  string     `json:"name"`
	Modes []jsonMode `json:"modes"`
}

func (f *jsonCollection) FindMode(name string) *jsonMode {
	for _, mode := range f.Modes {
		if mode.Name == name {
			return &mode
		}
	}
	return nil
}

type jsonMode struct {
	Name      string         `json:"name"`
	Variables []jsonVariable `json:"variables"`
}

func (f *jsonMode) FindVariable(name string) *jsonVariable {
	for _, variable := range f.Variables {
		if variable.Name == name {
			return &variable
		}
	}
	return nil
}

type jsonVariable struct {
	Name    string      `json:"name"`
	Type    string      `json:"type"`
	IsAlias bool        `json:"isAlias"`
	Value   interface{} `json:"value"`
}

type jsonAlias struct {
	Collection string
	Name       string
}

func (v *jsonVariable) Alias(currentCollection string) jsonAlias {
	if v.IsAlias {
		switch v := v.Value.(type) {
		case map[string]interface{}:
			switch name := v["name"].(type) {
			case string:
				switch collection := v["collection"].(type) {
				case string:
					return jsonAlias{
						Collection: collection,
						Name:       name,
					}
				default:
					return jsonAlias{
						Collection: currentCollection,
						Name:       name,
					}
				}
			}
		}
	}
	return jsonAlias{}
}
