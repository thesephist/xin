package xin

import (
	"bytes"
	"strings"
)

type StringValue []byte

func (v StringValue) String() string {
	return string(v)
}

func (v StringValue) Repr() string {
	return "'" + strings.ReplaceAll(string(v), "'", "\\'") + "'"
}

func (v StringValue) Equal(o Value) bool {
	if ov, ok := o.(StringValue); ok {
		return bytes.Compare(v, ov) == 0
	}

	return false
}

func strGetForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 2 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 2,
			given:    len(args),
		}
	}

	first, second := args[0], args[1]

	firstStr, fok := first.(StringValue)
	secondInt, sok := second.(IntValue)
	if fok && sok {
		if int(secondInt) >= 0 && int(secondInt) < len(firstStr) {
			return firstStr[secondInt : secondInt+1], nil
		}

		return zeroValue, nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func strSizeForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 1,
			given:    len(args),
		}
	}

	first := args[0]

	if firstString, ok := first.(StringValue); ok {
		return IntValue(len(firstString)), nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func strSliceForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 3 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 3,
			given:    len(args),
		}
	}

	first, second, third := args[0], args[1], args[2]

	firstStr, fok := first.(StringValue)
	secondInt, sok := second.(IntValue)
	thirdInt, tok := third.(IntValue)
	if fok && sok && tok {
		max := len(firstStr)
		inRange := func(iv IntValue) bool {
			return int(iv) >= 0 && int(iv) <= max
		}

		if int(secondInt) < 0 {
			secondInt = 0
		}
		if int(thirdInt) < 0 {
			thirdInt = 0
		}

		if int(secondInt) > max {
			secondInt = IntValue(max)
		}
		if int(thirdInt) > max {
			thirdInt = IntValue(max)
		}

		if inRange(secondInt) && inRange(thirdInt) && secondInt <= thirdInt {
			byteSlice := firstStr[secondInt:thirdInt]
			destSlice := make([]byte, len(byteSlice))
			copy(destSlice, byteSlice)
			return StringValue(destSlice), nil
		}

		return zeroValue, nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func strEncForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 1,
			given:    len(args),
		}
	}

	first := args[0]

	if firstStr, ok := first.(StringValue); ok {
		if len(firstStr) < 1 {
			return zeroValue, nil
		}

		return IntValue(firstStr[0]), nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func strDecForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 1,
			given:    len(args),
		}
	}

	first := args[0]

	if firstInt, ok := first.(IntValue); ok {
		if firstInt < 0 || firstInt > 255 {
			return zeroValue, nil
		}

		return StringValue([]byte{byte(firstInt)}), nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}
