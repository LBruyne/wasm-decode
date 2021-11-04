package decode

import (
	"fmt"
	"github.com/LBruyne/wasm-decode/types"
	"io"
)

// DecodeModule decodes a WASM module from io.Reader which contains the bytes streeam of .wasm file
func DecodeModule(r io.Reader) (mod *types.Module, err error) {
	mod = &types.Module{}
	if err := mod.Decode(r); err != nil {
		return nil, fmt.Errorf("decode module: %w", err)
	}
	return mod, nil
}
