package xin

import (
	"math"
	"math/rand"
)

func mathRandForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	return FracValue(rand.Float64()), nil
}

func mathSinForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 1,
			given:    len(args),
		}
	}

	first := args[0]

	switch cleanFirst := first.(type) {
	case IntValue:
		return FracValue(math.Sin(float64(cleanFirst))), nil
	case FracValue:
		return FracValue(math.Sin(float64(cleanFirst))), nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func mathCosForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 1,
			given:    len(args),
		}
	}

	first := args[0]

	switch cleanFirst := first.(type) {
	case IntValue:
		return FracValue(math.Cos(float64(cleanFirst))), nil
	case FracValue:
		return FracValue(math.Cos(float64(cleanFirst))), nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func mathTanForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 1,
			given:    len(args),
		}
	}

	first := args[0]

	switch cleanFirst := first.(type) {
	case IntValue:
		return FracValue(math.Tan(float64(cleanFirst))), nil
	case FracValue:
		return FracValue(math.Tan(float64(cleanFirst))), nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func mathAsinForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 1,
			given:    len(args),
		}
	}

	first := args[0]

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanFirst > 1 || cleanFirst < -1 {
			return zeroValue, nil
		}

		return FracValue(math.Asin(float64(cleanFirst))), nil
	case FracValue:
		if cleanFirst > 1 || cleanFirst < -1 {
			return zeroValue, nil
		}

		return FracValue(math.Asin(float64(cleanFirst))), nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func mathAcosForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 1,
			given:    len(args),
		}
	}

	first := args[0]

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanFirst > 1 || cleanFirst < -1 {
			return zeroValue, nil
		}

		return FracValue(math.Acos(float64(cleanFirst))), nil
	case FracValue:
		if cleanFirst > 1 || cleanFirst < -1 {
			return zeroValue, nil
		}

		return FracValue(math.Acos(float64(cleanFirst))), nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func mathAtanForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 1,
			given:    len(args),
		}
	}

	first := args[0]

	switch cleanFirst := first.(type) {
	case IntValue:
		return FracValue(math.Atan(float64(cleanFirst))), nil
	case FracValue:
		return FracValue(math.Atan(float64(cleanFirst))), nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func mathLnForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 1,
			given:    len(args),
		}
	}

	first := args[0]

	switch cleanFirst := first.(type) {
	case IntValue:
		return FracValue(math.Log(float64(cleanFirst))), nil
	case FracValue:
		return FracValue(math.Log(float64(cleanFirst))), nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}
