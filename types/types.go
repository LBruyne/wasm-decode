package types

import (
	"fmt"
	"github.com/LBruyne/wasm-decode/common"
	"io"
)

var (
	ValueTypeI32 = ValueType{
		Type:     "i32",
		Bytecode: 0x7f,
	}
	ValueTypeI64 = ValueType{
		Type:     "i64",
		Bytecode: 0x7e,
	}
	ValueTypeF32 = ValueType{
		Type:     "f32",
		Bytecode: 0x7d,
	}
	ValueTypeF64 = ValueType{
		Type:     "f64",
		Bytecode: 0x7c,
	}
)

type ValueType struct {
	Type     string
	Bytecode byte
}

func getValueType(bc byte) (ValueType, error) {
	switch bc {
	case ValueTypeF64.Bytecode:
		return ValueTypeF64, nil
	case ValueTypeF32.Bytecode:
		return ValueTypeF32, nil
	case ValueTypeI64.Bytecode:
		return ValueTypeI64, nil
	case ValueTypeI32.Bytecode:
		return ValueTypeI32, nil
	default:
		return ValueType{}, fmt.Errorf("invalid value type: %v", bc)
	}
}

// readValueTypes read s ValueTypes from r
func readValueTypes(r io.Reader, s uint32) ([]ValueType, error) {
	ret := make([]ValueType, s)
	for i := range ret {
		vt, err := readValueType(r)
		if err != nil {
			return nil, fmt.Errorf("read %v-th value type: %w", i, err)
		}
		ret = append(ret, vt)
	}
	return ret, nil
}

// readValueType read a ValueType from r
func readValueType(r io.Reader) (ValueType, error) {
	b, err := common.ReadByte(r)
	if err != nil {
		return ValueType{}, err
	}

	vt, err := getValueType(b)
	if err != nil {
		return ValueType{}, err
	}
	return vt, nil
}

type FunctionType struct {
	InputType, ReturnType []ValueType
}

// readValueTypes read a FunctionType from r
func readFunctionType(r io.Reader) (*FunctionType, error) {
	// first read a byte `0x60`
	b := make([]byte, 1)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, fmt.Errorf("read leading byte: %w", err)
	}

	if b[0] != common.FuncType {
		return nil, fmt.Errorf("%w: %#x != 0x60", common.ErrInvalidByte, b[0])
	}

	// read inputs
	is, _, err := common.DecodeUint32(r)
	if err != nil {
		return nil, fmt.Errorf("get the size of input value types: %w", err)
	}

	in, err := readValueTypes(r, is)
	if err != nil {
		return nil, fmt.Errorf("read value types of inputs: %w", err)
	}

	// read outputs
	os, _, err := common.DecodeUint32(r)
	if err != nil {
		return nil, fmt.Errorf("get the size of output value types: %w", err)
	}

	out, err := readValueTypes(r, os)
	if err != nil {
		return nil, fmt.Errorf("read value types of outputs: %w", err)
	}

	return &FunctionType{
		InputType:  in,
		ReturnType: out,
	}, nil
}

type TableType

type GlobalSegment

type ElementSegment

type MemoryType

type DataSegment

type ImportSegment

type ExportSegment

type CodeSegment
