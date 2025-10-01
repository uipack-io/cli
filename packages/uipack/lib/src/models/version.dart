import '../encoding/binary_encoding.dart';

class Version {
  int major;
  int minor;

  Version({required this.major, required this.minor});

  void encode(ByteDataWriter writer) {
    writer.writeUint16(major);
    writer.writeUint16(minor);
  }

  void decode(ByteDataReader reader) {
    major = reader.readUint16();
    minor = reader.readUint16();
  }

  @override
  String toString() => '$major.$minor';
}

