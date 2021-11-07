package cli

import (
	"github.com/LBruyne/wasm-decode/decode"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	fileName = "../examples/wasm/test.wasm"
)

func TestDump(t *testing.T) {
	mod, err := decode.DecodeFile(fileName)
	assert.Nil(t, err)
	assert.NotNil(t, mod)

	d := NewDumper(mod)
	d.Dump()
}
