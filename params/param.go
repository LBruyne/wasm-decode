package params

var (
	NoStartFunction = -1

	// MagicNumber number of WASM 1.0
	MagicNumber = []byte{0x00, 0x61, 0x73, 0x6D} // `0asm`

	// Version number of WASM 1.0
	Version = []byte{0x01, 0x00, 0x00, 0x00}
)
