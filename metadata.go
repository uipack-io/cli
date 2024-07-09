package uipack

import (
	"bufio"
	"strings"
)

type MainValueType uint8

type ValueType struct {
	Type       MainValueType
	CustomType CustomTypeIdentifier
}

const (
	DeprecatedType MainValueType = iota
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
	CustomType
)

func FindVariableType(o interface{}) MainValueType {
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
		return CustomType
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

func (b *BundleMetadata) FindTypeDefintion(identifier CustomTypeIdentifier) *TypeDefinition {
	for _, t := range b.Types {
		if t.Identifier == identifier {
			return &t
		}
	}
	return nil
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
	Identifier uint64    // Unique identifier, also the index of the variable.
	Type       ValueType // Type of the variable.
	Name       string    // Name of the variables.
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

func (mv *ModeVariantMetadata) Decode(r *bufio.Reader) error {
	var err error
	mv.Name, err = readString(r)
	if err != nil {
		return err
	}
	return nil
}

func (mv *ModeVariantMetadata) Encode(r *bufio.Writer) error {
	err := writeString(r, mv.Name)
	if err != nil {
		return err
	}
	return nil
}

func (b *BundleMetadata) Decode(r *bufio.Reader) error {
	var err error
	protocol, err := readUint16(r)
	if err != nil {
		return err
	}
	if protocol != PROTOCOL_VERSION {
		panic("Unsupported protocol version")
	}
	err = b.Version.Decode(r)
	if err != nil {
		return err
	}
	b.Name, err = readString(r)
	if err != nil {
		return err
	}
	modesCount, err := readUint8(r)
	if err != nil {
		return err
	}
	modes := make([]ModeMetadata, modesCount)
	for i := uint8(0); i < modesCount; i++ {
		mode := ModeMetadata{
			Identifier: Uint4(i),
		}
		err = mode.Decode(r)
		if err != nil {
			return err
		}
		modes[i] = mode
	}
	b.Modes = modes
	variablesCount, err := readUint64(r)
	if err != nil {
		return err
	}
	variables := make([]VariableMetadata, variablesCount)
	for i := uint64(0); i < variablesCount; i++ {
		variable := VariableMetadata{
			Identifier: i,
		}
		err = variable.Decode(r)
		if err != nil {
			return err
		}
		variables[i] = variable
	}
	b.Variables = variables
	return nil
}

func (b *BundleMetadata) Encode(r *bufio.Writer) error {
	err := writeUint16(r, PROTOCOL_VERSION)
	if err != nil {
		return err
	}
	err = b.Version.Encode(r)
	if err != nil {
		return err
	}
	err = writeString(r, b.Name)
	if err != nil {
		return err
	}
	err = writeUint8(r, uint8(len(b.Modes)))
	if err != nil {
		return err
	}
	for _, mode := range b.Modes {
		err = mode.Encode(r)
		if err != nil {
			return err
		}
	}
	err = writeUint64(r, uint64(len(b.Variables)))
	if err != nil {
		return err
	}
	for _, variable := range b.Variables {
		err = variable.Encode(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *ModeMetadata) Decode(r *bufio.Reader) error {
	var err error
	m.Name, err = readString(r)
	if err != nil {
		return err
	}
	uipackCount, err := readUint8(r)
	if err != nil {
		return err
	}
	variations := make([]ModeVariantMetadata, uipackCount)
	for i := uint8(0); i < uipackCount; i++ {
		variant := ModeVariantMetadata{
			Identifier: Uint4(i),
		}
		err = variant.Decode(r)
		if err != nil {
			return err
		}
		variations[i] = variant
	}
	m.Variants = variations
	return nil
}

func (m *ModeMetadata) Encode(r *bufio.Writer) error {
	err := writeString(r, m.Name)
	if err != nil {
		return err
	}
	err = writeUint8(r, uint8(len(m.Variants)))
	if err != nil {
		return err
	}
	for _, variant := range m.Variants {
		err = variant.Encode(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *VariableMetadata) Decode(r *bufio.Reader) error {
	mtype := ValueType{}
	err := mtype.Decode(r)
	if err != nil {
		return err
	}
	v.Type = mtype
	v.Name, err = readString(r)
	if err != nil {
		return err
	}
	return nil
}

func (v *VariableMetadata) Encode(r *bufio.Writer) error {
	err := v.Type.Encode(r)
	if err != nil {
		return err
	}
	err = writeString(r, v.Name)
	if err != nil {
		return err
	}
	return nil
}

func (v *ValueType) Decode(r *bufio.Reader) error {
	var err error
	t, err := readUint8(r)
	if err != nil {
		return err
	}
	v.Type = MainValueType(t)
	if v.Type == CustomType {
		ct, err := readUint32(r)
		if err != nil {
			return err
		}
		v.CustomType = CustomTypeIdentifier(ct)
	}
	return nil
}

func (v *ValueType) Encode(r *bufio.Writer) error {
	err := writeUint8(r, uint8(v.Type))
	if err != nil {
		return err
	}
	if v.Type == CustomType {
		err = writeUint32(r, uint32(v.CustomType))
		if err != nil {
			return err
		}
	}
	return nil
}
