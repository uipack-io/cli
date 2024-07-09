package uipack

import (
	"bufio"
	"errors"
)

type Bundle struct {
	Version Version
	Variant Variant
	Values  []interface{}
}

func (variant *Bundle) Decode(r *bufio.Reader, metadata *BundleMetadata) error {
	protocol, err := readUint16(r)
	if err != nil {
		return err
	}
	if protocol != PROTOCOL_VERSION {
		return errors.New("unsupported protocol version")
	}
	err = variant.Version.Decode(r)
	if err != nil {
		return err
	}
	v, err := readUint64(r)
	if err != nil {
		return err
	}
	variant.Variant = Variant(v)
	values := make([]interface{}, len(metadata.Variables))

	for i, variable := range metadata.Variables {
		values[i], err = metadata.DecodeArgument(variable.Type, r)
		if err != nil {
			return err
		}
	}

	variant.Values = values
	return nil
}

func (variant *Bundle) Encode(w *bufio.Writer, metadata *BundleMetadata) error {
	err := writeUint16(w, PROTOCOL_VERSION)
	if err != nil {
		return err
	}
	err = variant.Version.Encode((*bufio.Writer)(w))
	if err != nil {
		return err
	}
	err = writeUint64(w, uint64(variant.Variant))
	if err != nil {
		return err
	}
	for _, variable := range variant.Values {
		err = metadata.EncodeArgument(w, variable)
		if err != nil {
			return err
		}
	}
	return nil
}

func (metadata *BundleMetadata) DecodeArgument(t ValueType, r *bufio.Reader) (interface{}, error) {
	switch MainValueType(t.Type) {
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
		err := result.Decode(r)
		return result, err
	case TextStyleType:
		result := TextStyle{}
		err := result.Decode(r)
		return result, err
	case LinearGradientType:
		result := LinearGradient{}
		err := result.Decode(r)
		return result, err
	case RadialGradientType:
		result := RadialGradient{}
		err := result.Decode(r)
		return result, err
	case LabelType:
		result := Label{}
		err := result.Decode(r)
		return result, err
	case OffsetType:
		result := Offset{}
		err := result.Decode(r)
		return result, err
	case RadiusType:
		result := Radius{}
		err := result.Decode(r)
		return result, err
	case BorderRadiusType:
		result := BorderRadius{}
		err := result.Decode(r)
		return result, err
	case CustomType:
		result := Instance{}
		err := result.Decode(r, metadata)
		return result, err
	default:
		return nil, errors.New("unknown type")
	}
}

func (metadata *BundleMetadata) EncodeArgument(w *bufio.Writer, data interface{}) error {
	switch v := data.(type) {
	case Color:
		return v.Encode(w)
	case TextStyle:
		return v.Encode(w)
	case LinearGradient:
		return v.Encode(w)
	case RadialGradient:
		return v.Encode(w)
	case Label:
		return v.Encode(w)
	case Offset:
		return v.Encode(w)
	case Instance:
		return v.Encode(w, metadata)
	case Radius:
		return v.Encode(w)
	case BorderRadius:
		return v.Encode(w)
	case uint64:
		return writeUint64(w, v)
	case bool:
		return writeBool(w, v)
	case float64:
		return writeFloat64(w, v)
	case string:
		return writeString(w, v)
	}

	return errors.New("unknown type")
}
