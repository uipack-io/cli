package uipack

import (
	"bufio"
	"strings"
)

type VariableType uint8

const (
	DeprecatedType VariableType = iota
	ColorType
	TextStyleType
	LinearGradientType
	RadialGradientType
	LabelType
	StringType
	IntegerType
	BooleanType
	FloatType
	OffsetType
	RadiusType
	BorderRadiusType
	InstanceType
)

func FindVariableType(o interface{}) VariableType {
	switch o.(type) {
	case Color:
		return ColorType
	case TextStyle:
		return TextStyleType
	case LinearGradient:
		return LinearGradientType
	case RadialGradient:
		return RadialGradientType
	case Label:
		return LabelType
	case string:
		return StringType
	case int64:
		return IntegerType
	case bool:
		return BooleanType
	case float64:
		return FloatType
	case Offset:
		return OffsetType
	case Radius:
		return RadiusType
	case BorderRadius:
		return BorderRadiusType
	case Instance:
		return InstanceType
	default:
		return DeprecatedType
	}
}

type BundleMetadata struct {
	Version   Version            // The version of the bundle.
	Name      string             // The name of the bundle.
	Modes     []ModeMetadata     // The modes of the bundle. Maximum 16 modes.
	Variables []VariableMetadata // The variables of the bundle.
	Types     []TypeDefinition   // The custom type definitions
}

func (b *BundleMetadata) GenerateModeCombinations() [][]ModeVariantMetadata {
	var result [][]ModeVariantMetadata
	combineModesHelper(b.Modes, 0, []ModeVariantMetadata{}, &result)
	return result
}

func combineModesHelper(modes []ModeMetadata, index int, current []ModeVariantMetadata, result *[][]ModeVariantMetadata) {
	if index == len(modes) {
		combination := make([]ModeVariantMetadata, len(current))
		copy(combination, current)
		*result = append(*result, combination)
		return
	}

	for _, variant := range modes[index].Variants {
		combineModesHelper(modes, index+1, append(current, variant), result)
	}
}

type ModeMetadata struct {
	Identifier Uint4                 // Unique identifier, also the index of the mode.
	Name       string                // Name of the mode.
	Variants   []ModeVariantMetadata // uipack of the mode. Maximum 16 uipack.
}

type ModeVariantMetadata struct {
	Identifier Uint4  // Unique identifier, also the index of the mode variant.
	Name       string // Name of the mode value.
}

type VariableMetadata struct {
	Identifier uint64       // Unique identifier, also the index of the variable.
	Type       VariableType // Type of the variable.
	Name       string       // Name of the variables.
}

// Creates a tree of the variables in the bundle by grouping them by following the path in the variable name.
func (b BundleMetadata) BuildTree() VariableCollection {
	return *b.buildTreeNode("", "")
}

func (b BundleMetadata) buildTreeNode(prefix string, name string) *VariableCollection {
	result := VariableCollection{
		Name: name,
	}

	for _, variable := range b.Variables {
		if len(prefix) == 0 || strings.HasPrefix(variable.Name, prefix) {

			relative_name := strings.TrimPrefix(variable.Name, prefix)
			split := strings.Split(relative_name, "/")
			for i := 0; i < len(split); {
				if split[i] == "" {
					split = append(split[:i], split[i+1:]...)
				} else {
					i++
				}
			}

			if len(split) == 1 {
				entry := VariableCollectionVariable{
					Variable: variable,
					Name:     split[0],
				}
				result.Variables = append(result.Variables, entry)
			} else {
				collectionName := split[0]
				found := false
				for _, collection := range result.Collections {
					if collection.Name == collectionName {
						found = true
						break
					}
				}
				if !found {
					subprefix := prefix + collectionName + "/"
					subcollection := b.buildTreeNode(subprefix, collectionName)
					result.Collections = append(result.Collections, *subcollection)
				}
			}
		}
	}

	return &result
}

type VariableCollection struct {
	Name        string
	Variables   []VariableCollectionVariable
	Collections []VariableCollection
}

type VariableCollectionVariable struct {
	Variable VariableMetadata
	Name     string
}

// Binary encoding

func (mv *ModeVariantMetadata) Decode(r *bufio.Reader) {
	mv.Name = readString(r)
}

func (mv *ModeVariantMetadata) Encode(r *bufio.Writer) {
	writeString(r, mv.Name)
}

func (b *BundleMetadata) Decode(r *bufio.Reader) {
	protocol := readUint16(r)
	if protocol != PROTOCOL_VERSION {
		panic("Unsupported protocol version")
	}
	b.Version.Decode(r)
	b.Name = readString(r)
	modesCount := readUint8(r)
	modes := make([]ModeMetadata, modesCount)
	for i := uint8(0); i < modesCount; i++ {
		mode := ModeMetadata{
			Identifier: Uint4(i),
		}
		mode.Decode(r)
		modes[i] = mode
	}
	b.Modes = modes
	variablesCount := readUint64(r)
	variables := make([]VariableMetadata, variablesCount)
	for i := uint64(0); i < variablesCount; i++ {
		variable := VariableMetadata{
			Identifier: i,
		}
		variable.Decode(r)
		variables[i] = variable
	}
	b.Variables = variables
}

func (b *BundleMetadata) Encode(r *bufio.Writer) {
	writeUint16(r, PROTOCOL_VERSION)
	b.Version.Encode(r)
	writeString(r, b.Name)
	writeUint8(r, uint8(len(b.Modes)))
	for _, mode := range b.Modes {
		mode.Encode(r)
	}
	writeUint64(r, uint64(len(b.Variables)))
	for _, variable := range b.Variables {
		variable.Encode(r)
	}
}

func (m *ModeMetadata) Decode(r *bufio.Reader) {
	m.Name = readString(r)
	uipackCount := readUint8(r)
	variations := make([]ModeVariantMetadata, uipackCount)
	for i := uint8(0); i < uipackCount; i++ {
		variant := ModeVariantMetadata{
			Identifier: Uint4(i),
		}
		variant.Decode(r)
		variations[i] = variant
	}
	m.Variants = variations
}

func (m *ModeMetadata) Encode(r *bufio.Writer) {
	writeString(r, m.Name)
	writeUint8(r, uint8(len(m.Variants)))
	for _, variant := range m.Variants {
		variant.Encode(r)
	}
}

func (v *VariableMetadata) Decode(r *bufio.Reader) {
	v.Type = VariableType(readUint8(r))
	v.Name = readString(r)
}

func (v *VariableMetadata) Encode(r *bufio.Writer) {
	writeUint8(r, uint8(v.Type))
	writeString(r, v.Name)
}
