import '../encoding/binary_encoding.dart';
import 'version.dart';
import 'variant.dart';

enum MainValueType {
  deprecated(0),
  color(1),
  textStyle(2),
  linearGradient(3),
  radialGradient(4),
  label(5),
  string(6),
  integer(7),
  boolean(8),
  float(9),
  offset(10),
  radius(11),
  borderRadius(12),
  custom(13);

  const MainValueType(this.value);
  final int value;

  static MainValueType fromValue(int value) {
    return values.firstWhere(
      (type) => type.value == value,
      orElse: () => MainValueType.deprecated,
    );
  }
}

class ValueType {
  MainValueType type;
  int customType;

  ValueType({required this.type, this.customType = 0});

  void encode(ByteDataWriter writer) {
    writer.writeUint8(type.value);
    if (type == MainValueType.custom) {
      writer.writeUint32(customType);
    }
  }

  void decode(ByteDataReader reader) {
    type = MainValueType.fromValue(reader.readUint8());
    if (type == MainValueType.custom) {
      customType = reader.readUint32();
    }
  }
}

class ModeVariantMetadata {
  Uint4 identifier;
  String name;

  ModeVariantMetadata({required this.identifier, required this.name});

  void encode(ByteDataWriter writer) {
    writer.writeString(name);
  }

  void decode(ByteDataReader reader) {
    name = reader.readString();
  }
}

class ModeMetadata {
  Uint4 identifier;
  String name;
  List<ModeVariantMetadata> variants;

  ModeMetadata({
    required this.identifier,
    required this.name,
    required this.variants,
  });

  void encode(ByteDataWriter writer) {
    writer.writeString(name);
    writer.writeUint8(variants.length);
    for (final variant in variants) {
      variant.encode(writer);
    }
  }

  void decode(ByteDataReader reader) {
    name = reader.readString();
    final variantsCount = reader.readUint8();
    variants = List.generate(variantsCount, (i) {
      final variant = ModeVariantMetadata(identifier: i, name: '');
      variant.decode(reader);
      return variant;
    });
  }
}

class VariableMetadata {
  int identifier;
  ValueType type;
  String name;

  VariableMetadata({
    required this.identifier,
    required this.type,
    required this.name,
  });

  void encode(ByteDataWriter writer) {
    type.encode(writer);
    writer.writeString(name);
  }

  void decode(ByteDataReader reader) {
    type = ValueType(type: MainValueType.deprecated);
    type.decode(reader);
    name = reader.readString();
  }
}

class TypeDefinition {
  int identifier;
  String name;

  TypeDefinition({required this.identifier, required this.name});
}

class BundleMetadata {
  Version version;
  String name;
  List<ModeMetadata> modes;
  List<VariableMetadata> variables;
  List<TypeDefinition> types;

  BundleMetadata({
    required this.version,
    required this.name,
    required this.modes,
    required this.variables,
    required this.types,
  });

  TypeDefinition? findTypeDefinition(int identifier) {
    for (final type in types) {
      if (type.identifier == identifier) {
        return type;
      }
    }
    return null;
  }

  List<List<ModeVariantMetadata>> generateModeCombinations() {
    final result = <List<ModeVariantMetadata>>[];
    _combineModesHelper(modes, 0, <ModeVariantMetadata>[], result);
    return result;
  }

  void _combineModesHelper(
    List<ModeMetadata> modes,
    int index,
    List<ModeVariantMetadata> current,
    List<List<ModeVariantMetadata>> result,
  ) {
    if (index == modes.length) {
      result.add(List.from(current));
      return;
    }

    for (final variant in modes[index].variants) {
      current.add(variant);
      _combineModesHelper(modes, index + 1, current, result);
      current.removeLast();
    }
  }

  void encode(ByteDataWriter writer) {
    writer.writeUint16(protocolVersion);
    version.encode(writer);
    writer.writeString(name);
    writer.writeUint8(modes.length);
    for (final mode in modes) {
      mode.encode(writer);
    }
    writer.writeUint64(variables.length);
    for (final variable in variables) {
      variable.encode(writer);
    }
  }

  void decode(ByteDataReader reader) {
    final protocol = reader.readUint16();
    if (protocol != protocolVersion) {
      throw Exception('Unsupported protocol version: $protocol');
    }

    version = Version(major: 0, minor: 0);
    version.decode(reader);
    name = reader.readString();

    final modesCount = reader.readUint8();
    modes = List.generate(modesCount, (i) {
      final mode = ModeMetadata(identifier: i, name: '', variants: []);
      mode.decode(reader);
      return mode;
    });

    final variablesCount = reader.readUint64();
    variables = List.generate(variablesCount, (i) {
      final variable = VariableMetadata(
        identifier: i,
        type: ValueType(type: MainValueType.deprecated),
        name: '',
      );
      variable.decode(reader);
      return variable;
    });
  }
}

