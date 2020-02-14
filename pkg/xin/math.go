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

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}

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

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}

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

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}

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

func mathLnForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
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
