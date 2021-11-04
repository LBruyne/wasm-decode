package decode

import (
	"fmt"
	"github.com/LBruyne/wasm-decode/common"
	"github.com/LBruyne/wasm-decode/params"
	"io"
)

// DecodeModule decodes a WASM module from io.Reader which contains the bytes streeam of .wasm file
func DecodeModule(r io.Reader) (mod *Module, err error) {
	ret := &Module{

	}

	// magic number
	buf := make([]byte, 4)
	if n, err := io.ReadFull(r, buf); err != nil || n != 4 {
		return nil, common.ErrInvalidMagicNumber
	}
	for i := 0; i < 4; i++ {
		if buf[i] != params.MagicNumber[i] {
			return nil, common.ErrInvalidMagicNumber
		}
	}
	ret.MagicNumber = params.MagicNumber

	// version
	if n, err := io.ReadFull(r, buf); err != nil || n != 4 {
		return nil, err
	}
	for i := 0; i < 4; i++ {
		if buf[i] != params.Version[i] {
			return nil, common.ErrInvalidVersion
		}
	}
	ret.Version = params.Version

	// read sections
	if err := ret.readSections(r); err != nil {
		return nil, fmt.Errorf("readSections failed: %w", err)
	}

	return ret, nil
}
