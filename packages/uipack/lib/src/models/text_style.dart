import '../encoding/binary_encoding.dart';

class FontVariation {
  String axis;
  double value;

  FontVariation({required this.axis, required this.value});

  void encode(ByteDataWriter writer) {
    writer.writeString(axis);
    writer.writeFloat64(value);
  }

  void decode(ByteDataReader reader) {
    axis = reader.readString();
    value = reader.readFloat64();
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is FontVariation && 
           other.axis == axis && 
           other.value == value;
  }

  @override
  int get hashCode => axis.hashCode ^ value.hashCode;
}

class TextStyle {
  String fontFamily;
  double fontSize;
  int fontWeight;
  double letterSpacing;
  double wordSpacing;
  double lineHeight;
  List<FontVariation> fontVariations;

  TextStyle({
    required this.fontFamily,
    required this.fontSize,
    required this.fontWeight,
    required this.letterSpacing,
    required this.wordSpacing,
    required this.lineHeight,
    required this.fontVariations,
  });

  void encode(ByteDataWriter writer) {
    writer.writeString(fontFamily);
    writer.writeFloat64(fontSize);
    writer.writeUint8(fontWeight);
    writer.writeFloat64(letterSpacing);
    writer.writeFloat64(wordSpacing);
    writer.writeFloat64(lineHeight);
    
    writer.writeUint32(fontVariations.length);
    for (final variation in fontVariations) {
      variation.encode(writer);
    }
  }

  void decode(ByteDataReader reader) {
    fontFamily = reader.readString();
    fontSize = reader.readFloat64();
    fontWeight = reader.readUint8();
    letterSpacing = reader.readFloat64();
    wordSpacing = reader.readFloat64();
    lineHeight = reader.readFloat64();
    
    final variationsCount = reader.readUint32();
    fontVariations = List.generate(variationsCount, (index) {
      final variation = FontVariation(axis: '', value: 0.0);
      variation.decode(reader);
      return variation;
    });
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is TextStyle &&
           other.fontFamily == fontFamily &&
           other.fontSize == fontSize &&
           other.fontWeight == fontWeight &&
           other.letterSpacing == letterSpacing &&
           other.wordSpacing == wordSpacing &&
           other.lineHeight == lineHeight &&
           _listEquals(other.fontVariations, fontVariations);
  }

  bool _listEquals<T>(List<T> a, List<T> b) {
    if (a.length != b.length) return false;
    for (int i = 0; i < a.length; i++) {
      if (a[i] != b[i]) return false;
    }
    return true;
  }

  @override
  int get hashCode => Object.hash(
    fontFamily,
    fontSize,
    fontWeight,
    letterSpacing,
    wordSpacing,
    lineHeight,
    Object.hashAll(fontVariations),
  );

  @override
  String toString() => 'TextStyle(family: $fontFamily, size: $fontSize, weight: $fontWeight)';
}