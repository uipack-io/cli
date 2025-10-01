typedef Uint4 = int;

class Variant {
  int _value;

  Variant(this._value);

  int get value => _value;

  Uint4 getMode(Uint4 modeIndex) {
    final effectiveModeIndex = modeIndex & 0xF;
    return (_value >> (effectiveModeIndex * 4)) & 0xF;
  }

  Variant setMode(Uint4 modeIndex, Uint4 value) {
    final effectiveModeIndex = modeIndex & 0xF;
    final effectiveValue = value & 0xF;
    final newValue = (_value & ~(0xF << (effectiveModeIndex * 4))) | 
                     (effectiveValue << (effectiveModeIndex * 4));
    return Variant(newValue);
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is Variant && other._value == _value;
  }

  @override
  int get hashCode => _value.hashCode;

  @override
  String toString() => '0x${_value.toRadixString(16).padLeft(2, '0')}';
}