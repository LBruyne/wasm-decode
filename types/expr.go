package types

import (
	"bytes"
	"fmt"
	"github.com/LBruyne/wasm-decode/common"
	"github.com/LBruyne/wasm-decode/operator"
	"io"
)

// ConstExpression const expression defines the optcode must be xx.const instruction and data is the immediate
type ConstExpression struct {
	OptCode operator.OptCode
	Data    []byte
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
		return nil, fmt.Errorf("read optcode: %w", err)
	}

	buf := new(bytes.Buffer)
	teeR := io.TeeReader(r, buf)

	optCode := operator.OptCode(b[0])
	switch optCode {
	case operator.OptCodeI32Const:
		_, _, err = common.DecodeInt32(teeR)
	case operator.OptCodeI64Const:
		_, _, err = common.DecodeInt64(teeR)
	case operator.OptCodeF32Const:
		_, err = ReadFloat32(teeR)
	case operator.OptCodeF64Const:
		_, err = ReadFloat64(teeR)
	case operator.OptCodeGlobalGet:
		_, _, err = common.DecodeUint32(teeR)
	default:
		return nil, fmt.Errorf("invalid byte for opt code: %#x", b[0])
	}

	if err != nil {
		return nil, fmt.Errorf("read value: %w", err)
	}

	if _, err := io.ReadFull(r, b); err != nil {
		return nil, fmt.Errorf("look for end optcode: %w", err)
	}

	if b[0] != byte(operator.OptCodeEnd) {
		return nil, fmt.Errorf("constant expression has not terminated")
	}

	return &ConstExpression{
		OptCode: optCode,
		Data:    buf.Bytes(),
	}, nil
}
