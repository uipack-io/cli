package uipack

const (
	PROTOCOL_VERSION uint16 = 1
)

// A four bits unsigned integer. It is actually an uint8 in memory, but actually limited to 4 bits.
type Uint4 uint8
