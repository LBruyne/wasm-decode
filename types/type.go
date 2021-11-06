package types

import (
	"fmt"
	"github.com/LBruyne/wasm-decode/common"
	"io"
)

const (
	// FuncType represents the function of a section type
	FuncType byte = 0x60

	ElemTypeFuncRef = 0x70

	LimitTypeOnlyMin       = 0
	LimitTypeBothMinAndMax = 1

	GlobalTypeNotMutable = 1
	GlobalTypeMutable    = 0
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
	b, err := ReadByte(r)
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

	if b[0] != FuncType {
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

type TableType struct {
	ElemType byte
	Limit    *LimitsType
}

func readTableType(r io.Reader) (*TableType, error) {
	et, err := ReadByte(r)
	if err != nil {
		return nil, fmt.Errorf("read element type: %w", err)
	}

	// TODO WASM 1.0 defines that element type must be 0x70(function ref)
	if et != ElemTypeFuncRef {
		return nil, fmt.Errorf("read element type: not be 0x70, which is defined by WASM 1.0")
	}

	l, err := readLimitsType(r)
	if err != nil {
		return nil, fmt.Errorf("read limits type: %w", err)
	}

	return &TableType{
		ElemType: et,
		Limit:    l,
	}, nil
}

type LimitsType struct {
	Tag byte
	Min uint32
	Max uint32
}

func readLimitsType(r io.Reader) (*LimitsType, error) {
	b, err := ReadByte(r)
	if err != nil {
		return nil, fmt.Errorf("read limits type tag: %w", err)
	}

	ret := &LimitsType{
		Tag: b,
	}
	switch b {
	case LimitTypeOnlyMin:
		ret.Min, _, err = common.DecodeUint32(r)
		if err != nil {
			return nil, fmt.Errorf("read min of limit: %w", err)
		}
	case LimitTypeBothMinAndMax:
		ret.Min, _, err = common.DecodeUint32(r)
		if err != nil {
			return nil, fmt.Errorf("read min of limit: %w", err)
		}
		ret.Max, _, err = common.DecodeUint32(r)
		if err != nil {
			return nil, fmt.Errorf("read max of limit: %w", err)
		}
	default:
		return nil, fmt.Errorf("invalid byte for limit type tag: %#x != 0x00 or 0x01", b)
	}

	return ret, nil
}

type GlobalType struct {
	Value   ValueType
	Mutable bool
}

func readGlobalType(r io.Reader) (*GlobalType, error) {
	vt, err := readValueType(r)
	if err != nil {
		return nil, fmt.Errorf("read value type: %w", err)
	}

	ret := &GlobalType{
		Value: vt,
	}

	mut, err := ReadByte(r)
	if err != nil {
		return nil, fmt.Errorf("read mutablity: %w", err)
	}

	switch mut {
	case GlobalTypeMutable:
		ret.Mutable = true
	case GlobalTypeNotMutable:
		ret.Mutable = false
	default:
		return nil, fmt.Errorf("invalid byte for mutability: %#x != 0x00 or 0x01", mut)
	}

	return ret, nil
}

type MemoryType = LimitsType

func readMemoryType(r io.Reader) (*MemoryType, error) {
	return readLimitsType(r)
}

type LocalValueType struct {
	Count uint32
	Type  ValueType
}

func readLocalValueType(r io.Reader) (*LocalValueType, error) {
	c, _, err := common.DecodeUint32(r)
	if err != nil {
		return nil, fmt.Errorf("read number of locals: %w", err)
	}

	vt, err := readValueType(r)
	if err != nil {
		return nil, fmt.Errorf("read value type of locals: %w", err)
	}

	return &LocalValueType{
		Count: c,
		Type:  vt,
	}, nil
}

type CodeSegmentBody []byte
