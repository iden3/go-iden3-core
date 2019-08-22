package db

import "bytes"

func concat(vs ...[]byte) []byte {
	var b bytes.Buffer
	for _, v := range vs {
		b.Write(v)
	}
	return b.Bytes()
}

func clone(b0 []byte) []byte {
	b1 := make([]byte, len(b0))
	copy(b1, b0)
	return b1
}
