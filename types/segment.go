package types

import (
	"fmt"
	"github.com/LBruyne/wasm-decode/common"
	"github.com/LBruyne/wasm-decode/operator"
	"io"
	"io/ioutil"
)

const (
	ImportTypeFunc   = 0
	ImportTypeTable  = 1
	ImportTypeMem    = 2
	ImportTypeGlobal = 3

	ExportTypeFunc   = 0
	ExportTypeTable  = 1
	ExportTypeMem    = 2
	ExportTypeGlobal = 3
)

type ImportDescription struct {
	Kind byte // represent what type is imported
	// possible value 0,1,2,3

	TypeIndex  uint32
	TableType  *TableType
	MemType    *MemoryType
	GlobalType *GlobalType
}

type ImportSegment struct {
	Name, Module string
	Desc         *ImportDescription
}

func readImportSegment(r io.Reader) (*ImportSegment, error) {
	mn, err := ReadString(r)
	if err != nil {
		return nil, fmt.Errorf("read module name of imported component: %w", err)
	}

	n, err := ReadString(r)
	if err != nil {
		return nil, fmt.Errorf("read name of imported component: %w", err)
	}

	desc, err := readImportDescription(r)
	if err != nil {
		return nil, fmt.Errorf("read description of imported component: %w", err)
	}

	return &ImportSegment{
		Module: mn,
		Name:   n,
		Desc:   desc,
	}, nil
}

func readImportDescription(r io.Reader) (*ImportDescription, error) {
	k, err := ReadByte(r)
	if err != nil {
		return nil, fmt.Errorf("read kind of import description: %w", err)
	}

	ret := &ImportDescription{
		Kind: k,
	}

	switch k {
	case ImportTypeFunc:
		ret.TypeIndex, _, err = common.DecodeUint32(r)
		if err != nil {
			return nil, fmt.Errorf("read type index: %w", err)
		}
	case ImportTypeTable:
		ret.TableType, err = readTableType(r)
		if err != nil {
			return nil, fmt.Errorf("read table type: %w", err)
		}
	case ImportTypeMem:
		ret.MemType, err = readMemoryType(r)
		if err != nil {
			return nil, fmt.Errorf("read memory type: %w", err)
		}
	case ImportTypeGlobal:
		ret.GlobalType, err = readGlobalType(r)
		if err != nil {
			return nil, fmt.Errorf("read global type: %w", err)
		}
	default:
		return nil, fmt.Errorf("invalid kind of import description: %v", k)
	}

	return ret, nil
}

type GlobalSegment struct {
	Type *GlobalType
	Init *InitExpression
}

func readGlobalSegment(r io.Reader) (*GlobalSegment, error) {
	gt, err := readGlobalType(r)
	if err != nil {
		return nil, fmt.Errorf("read global type: %w", err)
	}

	init, err := readInitExpression(r)
	if err != nil {
		return nil, fmt.Errorf("read expression: %w", err)
	}

	return &GlobalSegment{
		Type: gt,
		Init: init,
	}, nil
}

type DataSegment struct {
	MemIdx uint32
	Offset *OffsetExpression
	Init   []byte
}

func readDataSegment(r io.Reader) (*DataSegment, error) {
	mi, _, err := common.DecodeUint32(r)
	if err != nil {
		return nil, fmt.Errorf("get memory index: %w", err)
	}

	// TODO WASM 1.0 defines that memory index must be 0
	if mi != 0 {
		return nil, fmt.Errorf("invalid memory index: %d, must be 0 which is defined by WASM 1.0", mi)
	}

	expr, err := readOffsetExpression(r)
	if err != nil {
		return nil, fmt.Errorf("read expr for offset: %w", err)
	}

	if expr.OpCode != operator.OpCodeI32Const {
		return nil, fmt.Errorf("offset expression must be i32.const but get %#x", expr.OpCode)
	}

	vs, _, err := common.DecodeUint32(r)
	if err != nil {
		return nil, fmt.Errorf("get size of vector: %w", err)
	}

	init := make([]byte, vs)
	if _, err := io.ReadFull(r, init); err != nil {
		return nil, fmt.Errorf("read bytes for init: %w", err)
	}

	return &DataSegment{
		MemIdx: mi,
		Offset: expr,
		Init:   init,
	}, nil
}

type ElementSegment struct {
	TableIdx uint32
	Offset   *OffsetExpression
	Init     []uint32 // function index
}

func readElementSegment(r io.Reader) (*ElementSegment, error) {
	ti, _, err := common.DecodeUint32(r)
	if err != nil {
		return nil, fmt.Errorf("get table index: %w", err)
	}

	expr, err := readOffsetExpression(r)
	if err != nil {
		return nil, fmt.Errorf("read expr for offset: %w", err)
	}

	if expr.OpCode != operator.OpCodeI32Const {
		return nil, fmt.Errorf("offset expression must be i32.const but get %#x", expr.OpCode)
	}

	vs, _, err := common.DecodeUint32(r)
	if err != nil {
		return nil, fmt.Errorf("get size of vector: %w", err)
	}

	init := make([]uint32, vs)
	for i := range init {
		fi, _, err := common.DecodeUint32(r)
		if err != nil {
			return nil, fmt.Errorf("read %v-th function index: %w", i, err)
		}
		init[i] = fi
	}

	return &ElementSegment{
		TableIdx: ti,
		Offset:   expr,
		Init:     init,
	}, nil
}

type ExportSegment struct {
	Name string
	Desc *ExportDescription
}

type ExportDescription struct {
	Kind  byte
	Index uint32
}

func readExportSegment(r io.Reader) (*ExportSegment, error) {
	name, err := ReadString(r)
	if err != nil {
		return nil, fmt.Errorf("read name of export module: %w", err)
	}

	desc, err := readExportDescription(r)
	if err != nil {
		return nil, fmt.Errorf("read export description: %w", err)
	}

	return &ExportSegment{
		Name: name,
		Desc: desc,
	}, nil
}

func readExportDescription(r io.Reader) (*ExportDescription, error) {
	k, err := ReadByte(r)
	if err != nil {
		return nil, fmt.Errorf("read kind of export description: %w", err)
	}

	// k is the Kind of export type
	// valid values are 0, 1, 2, 3
	if k >= 0x04 {
		return nil, fmt.Errorf("invalid byte for export description: %#x", k)
	}

	id, _, err := common.DecodeUint32(r)
	if err != nil {
		return nil, fmt.Errorf("read idx: %w", err)
	}

	return &ExportDescription{
		Kind:  k,
		Index: id,
	}, nil
}

type CodeSegment struct {
	Locals    []*LocalValueType
	NumLocals uint32
	Body      CodeSegmentBody
}

func readCodeSegment(r io.Reader) (*CodeSegment, error) {
	ss, _, err := common.DecodeUint32(r)
	if err != nil {
		return nil, fmt.Errorf("get the size of code segment: %w", err)
	}

	r = io.LimitReader(r, int64(ss))

	// parse locals
	ls, _, err := common.DecodeUint32(r)
	if err != nil {
		return nil, fmt.Errorf("get the size locals: %w", err)
	}

	locals := make([]*LocalValueType, ls)
	for i := range locals {
		l, err := readLocalValueType(r)
		if err != nil {
			return nil, fmt.Errorf("read %v-th local value: %w", i, err)
		}
		locals[i] = l
	}

	// parse code body
	cb, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read code body: %w", err)
	}
	if operator.OpCode(cb[len(cb)-1]) != operator.OpCodeEnd {
		return nil, fmt.Errorf("read code body: invalid end OpCode")
	}

	return &CodeSegment{
		Body:      cb,
		Locals:    locals,
		NumLocals: getNumLocals(locals),
	}, nil
}
