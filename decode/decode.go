package decode

import (
	"bytes"
	"fmt"
	"github.com/LBruyne/wasm-decode/types"
	"io"
	"io/ioutil"
)

// DecodeModule decodes a WASM module from io.Reader which contains the bytes streeam of .wasm file
func DecodeModule(r io.Reader) (mod *types.Module, err error) {
	mod = &types.Module{}
	if err := mod.Decode(r); err != nil {
		return nil, fmt.Errorf("decode module: %w", err)
	}
	return mod, nil
}

func DecodeFile(fn string) (*types.Module, error) {
	bs, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, fmt.Errorf("read file %v: %w", fn, err)
	}

	if mod, err := DecodeModule(bytes.NewBuffer(bs)); err != nil {
		return nil, fmt.Errorf("decode bytes: %w", err)
	} else {
		return mod, nil
	}
}
