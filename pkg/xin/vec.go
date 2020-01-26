package xin

import (
	"strings"
)

type VecValue []Value

func (v VecValue) String() string {
	ss := make([]string, len(v))
	for i, item := range v {
		ss[i] = item.String()
	}
	return "(<vec> " + strings.Join(ss, " ") + ")"
}

func (v VecValue) Equal(o Value) bool {
	if ov, ok := o.(VecValue); ok {
		if len(v) != len(ov) {
			return false
		}

		for i, x := range v {
			if x != ov[i] {
				return false
			}
		}

		return true
	}

	return false
}

func vecForm(fr *Frame, args []Value) (Value, error) {
	vecValues := make([]Value, len(args))
	for i, a := range args {
		val, err := unlazy(a)
		if err != nil {
			return nil, err
		}
		vecValues[i] = val
	}
	return VecValue(vecValues), nil
}

func vecGetForm(fr *Frame, args []Value) (Value, error) {
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

	isFirstVec, fok := first.(VecValue)
	isSecondInt, sok := second.(IntValue)
	if fok && sok {
		if int(isSecondInt) < len(isFirstVec) {
			return isFirstVec[isSecondInt], nil
		}

		return VecValue{}, nil
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func vecSetForm(fr *Frame, args []Value) (Value, error) {
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

	isFirstVec, fok := first.(VecValue)
	isSecondInt, sok := second.(IntValue)
	if fok && sok {
		if int(isSecondInt) < len(isFirstVec) {
			isFirstVec[isSecondInt] = third
			return third, nil
		}

		return VecValue{}, nil
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func vecSizeForm(fr *Frame, args []Value) (Value, error) {
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

	if isFirstVec, fok := first.(VecValue); fok {
		return IntValue(len(isFirstVec)), nil
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}
