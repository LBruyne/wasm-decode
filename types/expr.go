package types

import (
	"bytes"
	"fmt"
	"github.com/LBruyne/wasm-decode/common"
	"github.com/LBruyne/wasm-decode/operator"
	"io"
)

// ConstExpression const expression defines the OpCode must be xx.const instruction and data is the immediate
type ConstExpression struct {
	OpCode operator.OpCode
	Data   []byte
}

type OffsetExpression = ConstExpression

func readOffsetExpression(r io.Reader) (*OffsetExpression, error) {
	return readConstExpression(r)
}

type InitExpression = ConstExpression

func readInitExpression(r io.Reader) (*InitExpression, error) {
	return readConstExpression(r)
}

func readConstExpression(r io.Reader) (*ConstExpression, error) {
	b := make([]byte, 1)
	_, err := io.ReadFull(r, b)
	if err != nil {
		return nil, fmt.Errorf("read OpCode: %w", err)
	}

	buf := new(bytes.Buffer)
	teeR := io.TeeReader(r, buf)

	OpCode := operator.OpCode(b[0])
	switch OpCode {
	case operator.OpCodeI32Const:
		_, _, err = common.DecodeInt32(teeR)
	case operator.OpCodeI64Const:
		_, _, err = common.DecodeInt64(teeR)
	case operator.OpCodeF32Const:
		_, err = ReadFloat32(teeR)
	case operator.OpCodeF64Const:
		_, err = ReadFloat64(teeR)
	case operator.OpCodeGlobalGet:
		_, _, err = common.DecodeUint32(teeR)
	default:
		return nil, fmt.Errorf("invalid byte for opt code: %#x", b[0])
	}

	if err != nil {
		return nil, fmt.Errorf("read value: %w", err)
	}

	if _, err := io.ReadFull(r, b); err != nil {
		return nil, fmt.Errorf("look for end OpCode: %w", err)
	}

	if b[0] != byte(operator.OpCodeEnd) {
		return nil, fmt.Errorf("constant expression has not terminated")
	}

	return &ConstExpression{
		OpCode: OpCode,
		Data:   buf.Bytes(),
	}, nil
}
