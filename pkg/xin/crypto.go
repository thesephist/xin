package xin

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
)

func cryptoRandForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	buf64 := make([]byte, 8)
	_, err := crand.Read(buf64)
	if err != nil {
		return zeroValue, nil
	}

	n, readBytes := binary.Varint(buf64)
	if readBytes <= 0 {
		return zeroValue, nil
	}

	r := rand.New(rand.NewSource(n))
	return FracValue(r.Float64()), nil
}
