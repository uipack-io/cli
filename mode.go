package uipack

// A combination of sixteen different modes, with each mode being a four-bit value (0-15).
type Variant int64

// Takes the four bytes of the mode and returns the mode at the given index.
//
// Examples :
//
//	... 0001 0000 0000 0000 -> Mode1 = 0, Mode2 = 0, Mode3 = 0, Mode4 = 1, ...
//	... 0001 0010 0000 0000 -> Mode1 = 0, Mode2 = 0, Mode3 = 2, Mode4 = 1, ...
func (m Variant) GetMode(modeIndex Uint4) Uint4 {
	// We only keep the last 4 bits of the modeIndex.
	effectiveModeIndex := modeIndex & 0xF
	return Uint4(m >> (effectiveModeIndex * 4) & 0xF)
}

func (m Variant) SetMode(modeIndex Uint4, value Uint4) Variant {
	// We only keep the last 4 bits of the modeIndex.
	effectiveModeIndex := modeIndex & 0xF
	// We only keep the last 4 bits of the value.
	effectiveValue := value & 0xF
	// We set the value at the given index.
	return m & ^(0xF<<(effectiveModeIndex*4)) | Variant(effectiveValue)<<(effectiveModeIndex*4)
}
