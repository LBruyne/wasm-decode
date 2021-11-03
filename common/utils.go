package common

import "io"

// ReadByte read one byte from io.Reader
func ReadByte(r io.Reader) (byte, error) {
	p := make([]byte, 1)
	_, err := r.Read(p)
	return p[0], err
}

