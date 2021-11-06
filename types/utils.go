package types

import (
	"encoding/binary"
	"fmt"
	"github.com/LBruyne/wasm-decode/common"
	"io"
	"math"
)

// ReadString try to read a string from io.Reader
func ReadString(r io.Reader) (string, error) {
	vs, _, err := common.DecodeUint32(r)
	if err != nil {
		return "", fmt.Errorf("read size of string: %w", err)
	}

	buf := make([]byte, vs)
	if _, err := io.ReadFull(r, buf); err != nil {
		return "", fmt.Errorf("read bytes of string: %w", err)
	}

	return string(buf), nil
}

// ReadByte read one byte from io.Reader
func ReadByte(r io.Reader) (byte, error) {
	p := make([]byte, 1)
	_, err := r.Read(p)
	return p[0], err
}

// IEEE 754
func ReadFloat32(r io.Reader) (float32, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return 0, err
	}
	raw := binary.LittleEndian.Uint32(buf)
	return math.Float32frombits(raw), nil
}

// IEEE 754
func ReadFloat64(r io.Reader) (float64, error) {
	buf := make([]byte, 8)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return 0, err
	}
	raw := binary.LittleEndian.Uint64(buf)
	return math.Float64frombits(raw), nil
}

// getNumLocals count the total number of locals in one code segment
func getNumLocals(locals []*LocalValueType) (numLocal uint32) {
	for _, lt := range locals {
		numLocal += lt.Count
	}
	return numLocal
}
