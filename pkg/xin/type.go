package xin

import (
	"strconv"
)

func stringForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
		return nil, IncorrectNumberOfArgsError{
			required: 1,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}

	return StringValue(first.String()), nil
}

func intForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
		return nil, IncorrectNumberOfArgsError{
			required: 1,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}

	switch val := first.(type) {
	case IntValue:
		return val, nil
	case FracValue:
		return IntValue(float64(val)), nil
	case StringValue:
		intVal, err := strconv.ParseInt(string(val), 10, 64)
		if err != nil {
			return zeroValue, nil
		}
		return IntValue(intVal), nil
	default:
		return zeroValue, nil
	}
}

func fracForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
		return nil, IncorrectNumberOfArgsError{
			required: 1,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}

	switch val := first.(type) {
	case IntValue:
		return FracValue(val), nil
	case FracValue:
		return val, nil
	case StringValue:
		floatVal, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return FracValue(0.0), nil
		}
		return FracValue(floatVal), nil
	default:
		return FracValue(0), nil
	}
}

func equalForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 2 {
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

	if firstInt, fok := first.(IntValue); fok {
		if _, sok := second.(FracValue); sok {
			first = FracValue(float64(firstInt))
		}
	} else if _, fok := first.(FracValue); fok {
		if secondInt, sok := second.(IntValue); sok {
			second = FracValue(float64(secondInt))
		}
	}

	if first.Equal(second) {
		return IntValue(1), nil
	} else {
		return zeroValue, nil
	}
}

func typeForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
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

	switch first.(type) {
	case IntValue:
		return NativeFormValue{
			name:   "int",
			evaler: intForm,
		}, nil
	case FracValue:
		return NativeFormValue{
			name:   "frac",
			evaler: fracForm,
		}, nil
	case StringValue:
		return NativeFormValue{
			name:   "str",
			evaler: streamForm,
		}, nil
	case VecValue:
		return NativeFormValue{
			name:   "vec",
			evaler: vecForm,
		}, nil
	case MapValue:
		return NativeFormValue{
			name:   "map",
			evaler: mapForm,
		}, nil
	case StreamValue:
		return NativeFormValue{
			name:   "stream",
			evaler: streamForm,
		}, nil
	case FormValue, NativeFormValue:
		return NativeFormValue{
			name: "form",
			evaler: func(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
				panic("form constructor as type should never be called")
			},
		}, nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}
