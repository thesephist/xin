package xin

import (
	"math/rand"
)

func mathRandForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	return FracValue(rand.Float64()), nil
}
