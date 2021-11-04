package decode

import "io"

type ValueType struct {
	Type     string
	Bytecode byte
}

var (
	ValueTypeI32 ValueType = ValueType{
		Type: "i32",
		Bytecode: 0x7f,
	}
	ValueTypeI64 ValueType = ValueType{
		Type: "i64",
		Bytecode: 0x7e,
	}
	ValueTypeF32 ValueType = ValueType{
		Type: "f32",
		Bytecode: 0x7d,
	}
	ValueTypeF64 ValueType = ValueType{
		Type: "f64",
		Bytecode: 0x7c,
	}
)

type FunctionType struct {
	InputType, ReturnType []ValueType
}

func readFunctionType(r io.Reader) (*FunctionType, error) {
	return nil, nil
}