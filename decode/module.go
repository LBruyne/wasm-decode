package decode

import (
	"errors"
	"fmt"
	"github.com/LBruyne/wasm-decode/common"
	"io"
)

// Module represent the wasm module
type Module struct {
	Version		 []byte
	MagicNumber  []byte

	SecType      []*FunctionType
	SecFunction  []uint32
	SecTable     []*TableType
	SecMemory    []*common.MemoryType
	SecGlobals   []*GlobalSegment
	SecElements  []*ElementSegment
	SecData      []*DataSegment
	SecStart     uint32
	SecImport    []*ImportSegment
	SecExport    map[string]*common.ExportSegment
	SecCodes   []*common.CodeSegment

	Logger protocol.Logger
}

