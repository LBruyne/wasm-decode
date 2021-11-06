package types

import (
	"fmt"
	"github.com/LBruyne/wasm-decode/common"
	"github.com/LBruyne/wasm-decode/params"
	"io"
)

// Module represent the wasm module
type Module struct {
	Version     []byte
	MagicNumber []byte

	SecType     []*FunctionType
	SecFunction []uint32
	SecTable    []*TableType
	SecMemory   []*MemoryType
	SecGlobal   []*GlobalSegment
	SecElement  []*ElementSegment
	SecData     []*DataSegment
	SecStart    interface{}
	SecImport   []*ImportSegment
	SecExport   []*ExportSegment
	SecCode     []*CodeSegment
	SecCustom   *CustomSec
}

// Decode decodes a wasm module from io.Reader which contains full bytecodes of .wasm file
func (m *Module) Decode(r io.Reader) error {
	// magic number
	buf := make([]byte, 4)
	if n, err := io.ReadFull(r, buf); err != nil || n != 4 {
		return common.ErrInvalidMagicNumber
	}
	for i := 0; i < 4; i++ {
		if buf[i] != params.MagicNumber[i] {
			return common.ErrInvalidMagicNumber
		}
	}
	m.MagicNumber = params.MagicNumber

	// version
	if n, err := io.ReadFull(r, buf); err != nil || n != 4 {
		return err
	}
	for i := 0; i < 4; i++ {
		if buf[i] != params.Version[i] {
			return common.ErrInvalidVersion
		}
	}
	m.Version = params.Version

	// read sections
	if err := m.readSections(r); err != nil {
		return fmt.Errorf("readSections failed: %w", err)
	}
	return nil
}
