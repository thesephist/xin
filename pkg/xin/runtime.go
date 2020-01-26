package xin

import (
	"fmt"
)

type DefaultFormValue struct {
	name string
	eval func(*Frame, []Value) (Value, error)
}

func (v DefaultFormValue) String() string {
	return fmt.Sprintf("Default form ")
}

func loadAllDefaultForms(fr *Frame) {
	loadDefaultForm(fr, "+", addForm)
}

func loadDefaultForm(fr *Frame, name string, evaler func(*Frame, []Value) (Value, error)) {
	fr.Put(name, DefaultFormValue{
		name: name,
		eval: evaler,
	})
}

type IncorrectNumberOfArgsError struct {
	required int
	given    int
}

func (e IncorrectNumberOfArgsError) Error() string {
	return fmt.Sprintf("Incorrect number of arguments: requires %d but got %d",
		e.required, e.given)
}

type MismatchedArgumentsError struct {
	args []Value
}

func (e MismatchedArgumentsError) Error() string {
	return fmt.Sprintf("Mismatched arguments")
}

func addForm(fr *Frame, args []Value) (Value, error) {
	if len(args) != 2 {
		return nil, IncorrectNumberOfArgsError{
			required: 2,
			given:    len(args),
		}
	}

	first, err := unlazy(fr, args[0])
	if err != nil {
		return nil, err
	}
	second, err := unlazy(fr, args[1])
	if err != nil {
		return nil, err
	}

	// TODO: write implementation
	if iFirst, iOk := first.(IntValue); iOk {
		if iSecond, jOk := second.(IntValue); jOk {
			return IntValue(int64(iFirst) + int64(iSecond)), nil
		}
	}

	return nil, MismatchedArgumentsError{}
}
