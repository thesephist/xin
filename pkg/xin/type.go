package xin

import (
	"strconv"
)

func stringForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
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

	return StringValue(first.String()), nil
}

func intForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
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

	switch val := first.(type) {
	case IntValue:
		return val, nil
	case FracValue:
		return IntValue(float64(val)), nil
	case StringValue:
		intVal, err := strconv.ParseInt(string(val), 10, 64)
		if err != nil {
			return IntValue(0), nil
		}
		return IntValue(intVal), nil
	default:
		return IntValue(0), nil
	}
}

func fracForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
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
