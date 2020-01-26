package xin

import (
	"fmt"
)

type formEvaler func(*Frame, []Value) (Value, error)

type DefaultFormValue struct {
	name string
	eval formEvaler
}

func (v DefaultFormValue) String() string {
	return fmt.Sprintf("Default form %s", v.name)
}

func loadAllDefaultForms(fr *Frame) {
	builtins := map[string]formEvaler{
		"+": addForm,
		"-": subtractForm,
		"*": multiplyForm,
		"/": divideForm,

		"&": andForm,
		"|": orForm,
		"^": xorForm,

		"vec":      vecForm,
		"vec-get":  vecGetForm,
		"vec-set!": vecSetForm,

		"map":      mapForm,
		"map-get":  mapGetForm,
		"map-set!": mapSetForm,
		"map-del!": mapDelForm,
	}

	for name, evaler := range builtins {
		loadDefaultForm(fr, name, evaler)
	}
}

func loadDefaultForm(fr *Frame, name string, evaler formEvaler) {
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

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return cleanFirst + cleanSecond, nil
		}
	case FracValue:
		if cleanSecond, ok := second.(FracValue); ok {
			return cleanFirst + cleanSecond, nil
		}
	case StringValue:
		if cleanSecond, ok := second.(StringValue); ok {
			return cleanFirst + cleanSecond, nil
		}
	case VecValue:
		if cleanSecond, ok := second.(VecValue); ok {
			return VecValue(append(cleanFirst, cleanSecond...)), nil
		}
	}

	return nil, MismatchedArgumentsError{}
}

func subtractForm(fr *Frame, args []Value) (Value, error) {
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

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return cleanFirst - cleanSecond, nil
		}
	case FracValue:
		if cleanSecond, ok := second.(FracValue); ok {
			return cleanFirst - cleanSecond, nil
		}
	}

	return nil, MismatchedArgumentsError{}
}

func multiplyForm(fr *Frame, args []Value) (Value, error) {
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

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return cleanFirst * cleanSecond, nil
		}
	case FracValue:
		if cleanSecond, ok := second.(FracValue); ok {
			return cleanFirst * cleanSecond, nil
		}
	case VecValue:
		if cleanSecond, ok := second.(IntValue); ok {
			result := make([]Value, 0)
			copy(result, cleanFirst)
			max := int(cleanSecond)
			for i := 0; i < max; i++ {
				result = append(result, cleanFirst...)
			}
			return VecValue(result), nil
		}
	}

	return nil, MismatchedArgumentsError{}
}

func divideForm(fr *Frame, args []Value) (Value, error) {
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

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return cleanFirst / cleanSecond, nil
		}
	case FracValue:
		if cleanSecond, ok := second.(FracValue); ok {
			return cleanFirst / cleanSecond, nil
		}
	}

	return nil, MismatchedArgumentsError{}
}

func andForm(fr *Frame, args []Value) (Value, error) {
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

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return cleanFirst & cleanSecond, nil
		}
	}

	return nil, MismatchedArgumentsError{}
}

func orForm(fr *Frame, args []Value) (Value, error) {
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

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return cleanFirst | cleanSecond, nil
		}
	}

	return nil, MismatchedArgumentsError{}
}

func xorForm(fr *Frame, args []Value) (Value, error) {
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

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return cleanFirst ^ cleanSecond, nil
		}
	}

	return nil, MismatchedArgumentsError{}
}
