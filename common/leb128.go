package common

import (
	"errors"
	"fmt"
	"io"
)

var (
	ErrReadByte = errors.New("readByte failed")
)

var sevenBits = [...]byte{
	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
	0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
	0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f,
	0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f,
	0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f,
	0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x7b, 0x7c, 0x7d, 0x7e, 0x7f,
}

// EncodeInt64 encode int64 to bytes
func EncodeInt64(num int64) (b []byte) {
	if num >= 0 && num <= 0x3f {
		return append(b, sevenBits[num])
	} else if num < 0 && num >= ^0x3f {
		return append(b, sevenBits[0x80+num])
	}

	for {
		c := uint8(num & 0x7f)
		s := uint8(num & 0x40)
		num >>= 7

		if (num != -1 || s == 0) && (num != 0 || s != 0) {
			c |= 0x80
		}

		b = append(b, c)

		if c&0x80 == 0 {
			break
		}
	}
	return b
}

// EncodeInt32 encode int32 to bytes
func EncodeInt32(num int32) []byte {
	return EncodeInt64(int64(num))
}

// EncodeUint64 encode uint64 to bytes
func EncodeUint64(num uint64) []byte {
	var b []byte
	for {
		c := uint8(num & 0x7f)
		num >>= 7
		if num != 0 {
			c |= 0x80
		}
		b = append(b, c)
		if c&0x80 == 0 {
			break
		}
	}
	return b
}

// EncodeUint32 encode uint32 to bytes
func EncodeUint32(num uint32) []byte {
	return EncodeUint64(uint64(num))
}

// DecodeUint32 decode the bytes to the uint32
func DecodeUint32(r io.Reader) (ret uint32, num uint64, err error) {
	const (
		uint32Mask  uint32 = 1 << 7
		uint32Mask2        = ^uint32Mask
	)

	//0,7,16,21,28,35
	for shift := 0; shift < 35; shift += 7 {
		b, err := readByteAsUint32(r)
		if err != nil {
			return 0, 0, fmt.Errorf("readByte failed: %w", err)
		}
		num++
		ret |= (b & uint32Mask2) << shift
		if b&uint32Mask == 0 {
			break
		}
	}
	return
}

func DecodeUint64(r io.Reader) (ret uint64, num uint64, err error) {
	const (
		uint64Mask  uint64 = 1 << 7
		uint64Mask2        = ^uint64Mask
	)
	for shift := 0; shift < 64; shift += 7 {
		b, err := readByteAsUint64(r)
		if err != nil {
			return 0, 0, fmt.Errorf("readByte failed: %w", err)
		}
		num++
		ret |= (b & uint64Mask2) << shift
		if b&uint64Mask == 0 {
			break
		}
	}
	return
}

func DecodeInt32(r io.Reader) (ret int32, num uint64, err error) {
	const (
		int32Mask  int32 = 1 << 7
		int32Mask2       = ^int32Mask
		int32Mask3       = 1 << 6
		int32Mask4       = ^0
	)
	var shift int
	var b int32
	for shift < 35 {
		b, err = readByteAsInt32(r)
		if err != nil {
			return 0, 0, fmt.Errorf("readByte failed: %w", err)
		}
		num++
		ret |= (b & int32Mask2) << shift
		shift += 7
		if b&int32Mask == 0 {
			break
		}
	}

	if shift < 32 && (b&int32Mask3) == int32Mask3 {
		ret |= int32Mask4 << shift
	}
	return
}

func DecodeInt64(r io.Reader) (ret int64, num uint64, err error) {
	const (
		int64Mask  int64 = 1 << 7
		int64Mask2       = ^int64Mask
		int64Mask3       = 1 << 6
		int64Mask4       = ^0
	)
	var shift int
	var b int64
	for shift < 64 {
		b, err = readByteAsInt64(r)
		if err != nil {
			return 0, 0, fmt.Errorf("readByte failed: %w", err)
		}
		num++
		ret |= (b & int64Mask2) << shift
		shift += 7
		if b&int64Mask == 0 {
			break
		}
	}

	if shift < 64 && (b&int64Mask3) == int64Mask3 {
		ret |= int64Mask4 << shift
	}
	return
}

func readByteAsUint32(r io.Reader) (uint32, error) {
	readByte, err := readByte(r)
	return uint32(readByte), err
}

func readByteAsInt32(r io.Reader) (int32, error) {
	b := make([]byte, 1)
	_, err := io.ReadFull(r, b)
	return int32(b[0]), err
}

func readByteAsUint64(r io.Reader) (uint64, error) {
	readByte, err := readByte(r)
	return uint64(readByte), err
}

func readByteAsInt64(r io.Reader) (int64, error) {
	readByte, err := readByte(r)
	return int64(readByte), err
}

// readByte read one byte from io.Reader
func readByte(r io.Reader) (byte, error) {
	p := make([]byte, 1)
	_, err := r.Read(p)
	return p[0], err
}
