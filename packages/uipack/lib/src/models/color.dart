import '../encoding/binary_encoding.dart';

enum ColorSpace {
  srgb(0),
  extendedSrgb(1),
  displayP3(2);

  const ColorSpace(this.value);
  final int value;

  static ColorSpace fromValue(int value) {
    switch (value) {
      case 0: return ColorSpace.srgb;
      case 1: return ColorSpace.extendedSrgb;
      case 2: return ColorSpace.displayP3;
      default: throw ArgumentError('Invalid ColorSpace value: $value');
    }
  }
}

class Color {
  double red;
  double green;
  double blue;
  double alpha;
  ColorSpace colorSpace;

  Color({
    required this.red,
    required this.green,
    required this.blue,
    required this.alpha,
    required this.colorSpace,
  });

  static final Color black = Color(
    red: 0.0,
    green: 0.0,
    blue: 0.0,
    alpha: 1.0,
    colorSpace: ColorSpace.displayP3,
  );

  void parseHexString(String value) {
    String hex = value.startsWith('#') ? value.substring(1) : value;
    
    switch (hex.length) {
      case 8:
        final intValue = int.parse(hex, radix: 16);
        alpha = (intValue & 0xFF) / 255.0;
        red = ((intValue >> 24) & 0xFF) / 255.0;
        green = ((intValue >> 16) & 0xFF) / 255.0;
        blue = ((intValue >> 8) & 0xFF) / 255.0;
        break;
      case 6:
        final intValue = int.parse(hex, radix: 16);
        alpha = 1.0;
        red = ((intValue >> 16) & 0xFF) / 255.0;
        green = ((intValue >> 8) & 0xFF) / 255.0;
        blue = (intValue & 0xFF) / 255.0;
        break;
      default:
        throw ArgumentError('Invalid color format: $value');
    }
  }

  String toHexString() {
    int toBytes(double value) => (value * 255.0).round() & 0xFF;
    return '${toBytes(alpha).toRadixString(16).padLeft(2, '0')}'
           '${toBytes(red).toRadixString(16).padLeft(2, '0')}'
           '${toBytes(green).toRadixString(16).padLeft(2, '0')}'
           '${toBytes(blue).toRadixString(16).padLeft(2, '0')}';
  }

  void encode(ByteDataWriter writer) {
    writer.writeUint8(colorSpace.value);
    writer.writeFloat64(red);
    writer.writeFloat64(green);
    writer.writeFloat64(blue);
    writer.writeFloat64(alpha);
  }

  void decode(ByteDataReader reader) {
    colorSpace = ColorSpace.fromValue(reader.readUint8());
    red = reader.readFloat64();
    green = reader.readFloat64();
    blue = reader.readFloat64();
    alpha = reader.readFloat64();
  }

  @override
  String toString() => '#${toHexString()}';
}