import '../encoding/binary_encoding.dart';
import 'version.dart';
import 'variant.dart';
import 'metadata.dart';
import 'color.dart';
import 'text_style.dart';

class Bundle {
  Version version;
  Variant variant;
  List<dynamic> values;

  Bundle({
    required this.version,
    required this.variant,
    required this.values,
  });

  void decode(ByteDataReader reader, BundleMetadata metadata) {
    final protocol = reader.readUint16();
    if (protocol != protocolVersion) {
      throw Exception('Unsupported protocol version: $protocol');
    }

    version = Version(major: 0, minor: 0);
    version.decode(reader);
    
    final variantValue = reader.readUint64();
    variant = Variant(variantValue);
    
    values = List.generate(metadata.variables.length, (i) {
      final variable = metadata.variables[i];
      return _decodeArgument(variable.type, reader, metadata);
    });
  }

  void encode(ByteDataWriter writer, BundleMetadata metadata) {
    writer.writeUint16(protocolVersion);
    version.encode(writer);
    writer.writeUint64(variant.value);
    
    for (final value in values) {
      _encodeArgument(writer, value, metadata);
    }
  }

  dynamic _decodeArgument(ValueType type, ByteDataReader reader, BundleMetadata metadata) {
    switch (type.type) {
      case MainValueType.integer:
        return reader.readUint64();
      case MainValueType.boolean:
        return reader.readBool();
      case MainValueType.float:
        return reader.readFloat64();
      case MainValueType.string:
        return reader.readString();
      case MainValueType.color:
        final color = Color(
          red: 0.0,
          green: 0.0,
          blue: 0.0,
          alpha: 1.0,
          colorSpace: ColorSpace.srgb,
        );
        color.decode(reader);
        return color;
      case MainValueType.textStyle:
        final textStyle = TextStyle(
          fontFamily: '',
          fontSize: 0.0,
          fontWeight: 400,
          letterSpacing: 0.0,
          wordSpacing: 0.0,
          lineHeight: 1.0,
          fontVariations: [],
        );
        textStyle.decode(reader);
        return textStyle;
      default:
        throw Exception('Unsupported type: ${type.type}');
    }
  }

  void _encodeArgument(ByteDataWriter writer, dynamic data, BundleMetadata metadata) {
    switch (data.runtimeType) {
      case Color:
        (data as Color).encode(writer);
        break;
      case TextStyle:
        (data as TextStyle).encode(writer);
        break;
      case int:
        writer.writeUint64(data as int);
        break;
      case bool:
        writer.writeBool(data as bool);
        break;
      case double:
        writer.writeFloat64(data as double);
        break;
      case String:
        writer.writeString(data as String);
        break;
      default:
        throw Exception('Unsupported type: ${data.runtimeType}');
    }
  }
}