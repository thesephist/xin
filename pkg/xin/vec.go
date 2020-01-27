package xin

import (
	"strings"
)

type vecUnderlying []Value

type VecValue struct {
	items vecUnderlying
}

func NewVecValue(items []Value) VecValue {
	return VecValue{
		items: items,
	}
}

func (v VecValue) String() string {
	ss := make([]string, len(v.items))
	for i, item := range v.items {
		ss[i] = item.String()
	}
	return "(<vec> " + strings.Join(ss, " ") + ")"
}

func (v VecValue) Equal(o Value) bool {
	if ov, ok := o.(VecValue); ok {
		if len(v.items) != len(ov.items) {
			return false
		}

		for i, x := range v.items {
			if x != ov.items[i] {
				return false
			}
		}

		return true
	}

	return false
}

func vecForm(fr *Frame, args []Value) (Value, InterpreterError) {
	vecValues := make([]Value, len(args))
	for i, a := range args {
		val, err := unlazy(a)
		if err != nil {
			return nil, err
		}
		vecValues[i] = val
	}
	return NewVecValue(vecValues), nil
}

func vecGetForm(fr *Frame, args []Value) (Value, InterpreterError) {
	if len(args) != 2 {
		return nil, IncorrectNumberOfArgsError{
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

	firstVec, fok := first.(VecValue)
	secondInt, sok := second.(IntValue)
	if fok && sok {
		if int(secondInt) < len(firstVec.items) {
			return firstVec.items[secondInt], nil
		}

		return VecValue{}, nil
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func vecSetForm(fr *Frame, args []Value) (Value, InterpreterError) {
	if len(args) != 3 {
		return nil, IncorrectNumberOfArgsError{
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

	firstVec, fok := first.(VecValue)
	secondInt, sok := second.(IntValue)
	if fok && sok {
		if int(secondInt) < len(firstVec.items) {
			firstVec.items[secondInt] = third
			return third, nil
		}

		return VecValue{}, nil
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func vecAddForm(fr *Frame, args []Value) (Value, InterpreterError) {
	if len(args) != 2 {
		return nil, IncorrectNumberOfArgsError{
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

	if firstVec, fok := first.(VecValue); fok {
		firstVec.items = append(firstVec.items, second)
		return second, nil
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func vecSizeForm(fr *Frame, args []Value) (Value, InterpreterError) {
	if len(args) != 1 {
		return nil, IncorrectNumberOfArgsError{
			required: 1,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}

	if firstVec, fok := first.(VecValue); fok {
		return IntValue(len(firstVec.items)), nil
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}
