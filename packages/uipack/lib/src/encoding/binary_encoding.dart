import 'dart:typed_data';
import 'dart:convert';

const int protocolVersion = 1;

class ByteDataWriter {
  final List<int> _buffer = [];

  List<int> get buffer => _buffer;

  void writeUint8(int value) {
    _buffer.add(value & 0xFF);
  }

  void writeUint16(int value) {
    final bytes = ByteData(2);
    bytes.setUint16(0, value, Endian.big);
    _buffer.addAll(bytes.buffer.asUint8List());
  }

  void writeUint32(int value) {
    final bytes = ByteData(4);
    bytes.setUint32(0, value, Endian.big);
    _buffer.addAll(bytes.buffer.asUint8List());
  }

  void writeUint64(int value) {
    final bytes = ByteData(8);
    bytes.setUint64(0, value, Endian.big);
    _buffer.addAll(bytes.buffer.asUint8List());
  }

  void writeFloat64(double value) {
    final bytes = ByteData(8);
    bytes.setFloat64(0, value, Endian.little);
    _buffer.addAll(bytes.buffer.asUint8List());
  }

  void writeString(String value) {
    final bytes = utf8.encode(value);
    writeUint32(bytes.length);
    _buffer.addAll(bytes);
  }

  void writeBool(bool value) {
    writeUint8(value ? 1 : 0);
  }
}

class ByteDataReader {
  final Uint8List _data;
  int _offset = 0;

  ByteDataReader(this._data);

  int readUint8() {
    if (_offset >= _data.length) throw Exception('Buffer underflow');
    return _data[_offset++];
  }

  int readUint16() {
    if (_offset + 2 > _data.length) throw Exception('Buffer underflow');
    final bytes = ByteData.sublistView(_data, _offset, _offset + 2);
    _offset += 2;
    return bytes.getUint16(0, Endian.big);
  }

  int readUint32() {
    if (_offset + 4 > _data.length) throw Exception('Buffer underflow');
    final bytes = ByteData.sublistView(_data, _offset, _offset + 4);
    _offset += 4;
    return bytes.getUint32(0, Endian.big);
  }

  int readUint64() {
    if (_offset + 8 > _data.length) throw Exception('Buffer underflow');
    final bytes = ByteData.sublistView(_data, _offset, _offset + 8);
    _offset += 8;
    return bytes.getUint64(0, Endian.big);
  }

  double readFloat64() {
    if (_offset + 8 > _data.length) throw Exception('Buffer underflow');
    final bytes = ByteData.sublistView(_data, _offset, _offset + 8);
    _offset += 8;
    return bytes.getFloat64(0, Endian.little);
  }

  String readString() {
    final length = readUint32();
    if (_offset + length > _data.length) throw Exception('Buffer underflow');
    final bytes = _data.sublist(_offset, _offset + length);
    _offset += length;
    return utf8.decode(bytes);
  }

  bool readBool() {
    return readUint8() == 1;
  }
}

