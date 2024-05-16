package uipack

import "bufio"

type Bundle struct {
	Version Version
	Variant Variant
	Values  []interface{}
}

func (variant *Bundle) Decode(r *bufio.Reader, metadata BundleMetadata) {
	protocol := readUint16(r)
	if protocol != PROTOCOL_VERSION {
		panic("Unsupported protocol version")
	}
	variant.Version.Decode(r)
	variant.Variant = Variant(readUint64(r))
	values := make([]interface{}, len(metadata.Variables))

	for i, variable := range metadata.Variables {
		values[i] = DecodeArgument(variable.Type, r)
	}

	variant.Values = values
}

func (variant *Bundle) Encode(w *bufio.Writer) {
	writeUint16(w, PROTOCOL_VERSION)
	variant.Version.Encode((*bufio.Writer)(w))
	writeUint64(w, uint64(variant.Variant))
	for _, variable := range variant.Values {
		EncodeArgument(w, variable)
	}
}

func DecodeArgument(t VariableType, r *bufio.Reader) interface{} {
	switch VariableType(t) {
	case IntegerType:
		return readUint64(r)
	case BooleanType:
		return readBool(r)
	case FloatType:
		return readFloat64(r)
	case StringType:
		return readString(r)
	case ColorType:
		result := Color{}
		result.Decode(r)
		return result
	case TextStyleType:
		result := TextStyle{}
		result.Decode(r)
		return result
	case LinearGradientType:
		result := LinearGradient{}
		result.Decode(r)
		return result
	case RadialGradientType:
		result := RadialGradient{}
		result.Decode(r)
		return result
	case LabelType:
		result := Label{}
		result.Decode(r)
		return result
	case OffsetType:
		result := Offset{}
		result.Decode(r)
		return result
	case RadiusType:
		result := Radius{}
		result.Decode(r)
		return result
	case BorderRadiusType:
		result := BorderRadius{}
		result.Decode(r)
		return result
	case InstanceType:
		result := Instance{}
		result.Decode(r)
		return result
	default:
		panic("Unknown variable type")
	}
}

func EncodeArgument(w *bufio.Writer, data interface{}) {
	switch v := data.(type) {
	case Color:
		v.Encode((*bufio.Writer)(w))
	case TextStyle:
		v.Encode((*bufio.Writer)(w))
	case LinearGradient:
		v.Encode((*bufio.Writer)(w))
	case RadialGradient:
		v.Encode((*bufio.Writer)(w))
	case Label:
		v.Encode((*bufio.Writer)(w))
	case Offset:
		v.Encode((*bufio.Writer)(w))
	case Instance:
		v.Encode((*bufio.Writer)(w))
	case Radius:
		v.Encode((*bufio.Writer)(w))
	case BorderRadius:
		v.Encode((*bufio.Writer)(w))
	case uint64:
		writeUint64(w, v)
	case bool:
		writeBool(w, v)
	case float64:
		writeFloat64(w, v)
	case string:
		writeString(w, v)
	}
}
