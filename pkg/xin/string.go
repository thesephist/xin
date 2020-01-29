package xin

import (
	"bytes"
)

type StringValue []byte

func (v StringValue) String() string {
	return string(v)
}

func (v StringValue) Equal(o Value) bool {
	if ov, ok := o.(StringValue); ok {
		return bytes.Compare(v, ov) == 0
	}

	return false
}

func strGetForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) != 2 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 2,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}
	second, err := unlazy(args[1])
	if err != nil {
		return nil, err
	}

	firstStr, fok := first.(StringValue)
	secondInt, sok := second.(IntValue)
	if fok && sok {
		if int(secondInt) < len(firstStr) {
			return firstStr[secondInt : secondInt+1], nil
		}

		return IntValue(0), nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func strSizeForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) != 1 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 1,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}

	if firstString, ok := first.(StringValue); ok {
		return IntValue(len(firstString)), nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func strSliceForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) != 3 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 3,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}
	second, err := unlazy(args[1])
	if err != nil {
		return nil, err
	}
	third, err := unlazy(args[2])
	if err != nil {
		return nil, err
	}

	firstStr, fok := first.(StringValue)
	secondInt, sok := second.(IntValue)
	thirdInt, tok := third.(IntValue)
	if fok && sok && tok {
		max := len(firstStr)
		inRange := func(iv IntValue) bool {
			return int(iv) >= 0 && int(iv) <= max
		}

		if int(secondInt) > max {
			secondInt = IntValue(max)
		}
		if int(thirdInt) > max {
			thirdInt = IntValue(max)
		}

		if inRange(secondInt) && inRange(thirdInt) {
			byteSlice := firstStr[secondInt:thirdInt]
			destSlice := make([]byte, len(byteSlice))
			copy(destSlice, byteSlice)
			return StringValue(destSlice), nil
		}

		return IntValue(0), nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}
