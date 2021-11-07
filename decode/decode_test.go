package decode

import (
	"bytes"
	"github.com/LBruyne/wasm-decode/common"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

var (
	fileName = "../examples/wasm/test.wasm"
)

func TestDecodeFile(t *testing.T) {
	mod, err := DecodeFile(fileName)
	assert.Nil(t, err)
	assert.NotNil(t, mod)
}

func TestDecodeModule(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		buf, err := ioutil.ReadFile(fileName)
		assert.Nil(t, err)

		mod, err := DecodeModule(bytes.NewBuffer(buf))
		assert.Nil(t, err)
		assert.NotNil(t, mod)
	})

	t.Run("invalid_magic_number", func(t *testing.T) {
		mod, err := DecodeModule(bytes.NewBuffer([]byte{}))
		assert.Nil(t, mod)
		assert.True(t, strings.Contains(err.Error(), common.ErrInvalidMagicNumber.Error()))

		mod, err = DecodeModule(bytes.NewBuffer([]byte{1, 2, 3, 4}))
		assert.Nil(t, mod)
		assert.True(t, strings.Contains(err.Error(), common.ErrInvalidMagicNumber.Error()))
	})

	t.Run("invalid_version", func(t *testing.T) {
		mod, err := DecodeModule(bytes.NewBuffer([]byte{0x00, 0x61, 0x73, 0x6D}))
		assert.Nil(t, mod)
		assert.True(t, strings.Contains(err.Error(), io.EOF.Error()))

		mod, err = DecodeModule(bytes.NewBuffer([]byte{0x00, 0x61, 0x73, 0x6D, 0x12, 0x12, 0x12, 0x12}))
		assert.Nil(t, mod)
		assert.True(t, strings.Contains(err.Error(), common.ErrInvalidVersion.Error()))
	})

	t.Run("read_section_fail", func(t *testing.T) {
		mod, err := DecodeModule(bytes.NewBuffer([]byte{0x00, 0x61, 0x73, 0x6D, 0x01, 0x00, 0x00, 0x00, 0x11, 0x00}))
		assert.Nil(t, mod)
		assert.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "readSections failed"))
	})
}
