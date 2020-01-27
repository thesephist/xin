package xin

import (
	"strings"
)

type MapValue map[Value]Value

func (v MapValue) String() string {
	i := 0
	ss := make([]string, len(v))
	for k, val := range v {
		ss[i] = k.String() + "->" + val.String()
		i++
	}
	return "(<map> " + strings.Join(ss, " ") + ")"
}

func (v MapValue) Equal(o Value) bool {
	if ov, ok := o.(MapValue); ok {
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

func mapForm(fr *Frame, args []Value) (Value, InterpreterError) {
	if len(args) != 0 {
		return nil, IncorrectNumberOfArgsError{
			required: 0,
			given:    len(args),
		}
	}

	return make(MapValue), nil
}

func mapGetForm(fr *Frame, args []Value) (Value, InterpreterError) {
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

	if firstMap, ok := first.(MapValue); ok {
		val, prs := firstMap[second]
		if prs {
			return val, nil
		}

		return VecValue{}, nil
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func mapSetForm(fr *Frame, args []Value) (Value, InterpreterError) {
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

	if firstMap, ok := first.(MapValue); ok {
		firstMap[second] = third

		return third, nil
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func mapDelForm(fr *Frame, args []Value) (Value, InterpreterError) {
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

	if firstMap, ok := first.(MapValue); ok {
		val, prs := firstMap[second]
		if prs {
			delete(firstMap, second)
			return val, nil
		}

		return VecValue{}, nil
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func mapSizeForm(fr *Frame, args []Value) (Value, InterpreterError) {
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

	if firstMap, fok := first.(MapValue); fok {
		return IntValue(len(firstMap)), nil
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}
