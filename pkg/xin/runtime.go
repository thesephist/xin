package xin

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type formEvaler func(*Vm, *Frame, []Value) (Value, error)

type DefaultFormValue struct {
	name   string
	evaler formEvaler
}

func (v DefaultFormValue) eval(fr *Frame, args []Value) (Value, error) {
	return v.evaler(fr.Vm, fr, args)
}

func (v DefaultFormValue) String() string {
	return fmt.Sprintf("Default form %s", v.name)
}

func (v DefaultFormValue) Equal(o Value) bool {
	if ov, ok := o.(DefaultFormValue); ok {
		return v.name == ov.name
	}

	return false
}

type RuntimeError struct {
	reason string
}

func (e RuntimeError) Error() string {
	return "Runtime error:" + e.reason
}

func loadAllDefaultValues(vm *Vm) {
	fr := vm.Frame

	fr.Put("os::stdout", StreamValue{
		sink: func(v Value) error {
			fmt.Printf(v.String())
			return nil
		},
	})
	fr.Put("os::stdin", StreamValue{
		source: func() (Value, error) {
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err == io.EOF {
				return StringValue(""), nil
			} else if err != nil {
				return nil, RuntimeError{
					reason: "Cannot read from stdin",
				}
			}

			return StringValue(input[:len(input)-1]), nil
		},
	})
	fr.Put("os::sleep", StreamValue{
		sink: func(v Value) error {
			if duration, ok := v.(IntValue); ok {
				time.Sleep(time.Duration(
					int64(float64(duration) * float64(time.Second)),
				))
				return nil
			} else if duration, ok := v.(FracValue); ok {
				time.Sleep(time.Duration(
					int64(float64(duration) * float64(time.Second)),
				))
				return nil
			}

			return MismatchedArgumentsError{
				args: []Value{v},
			}
		},
	})
}

func loadAllDefaultForms(vm *Vm) {
	builtins := map[string]formEvaler{
		"+": addForm,
		"-": subtractForm,
		"*": multiplyForm,
		"/": divideForm,

		">": greaterForm,
		"<": lessForm,
		"=": equalForm,

		"&": andForm,
		"|": orForm,
		"^": xorForm,

		"vec":      vecForm,
		"vec-get":  vecGetForm,
		"vec-set!": vecSetForm,
		"vec-size": vecSizeForm,

		"map":      mapForm,
		"map-get":  mapGetForm,
		"map-set!": mapSetForm,
		"map-del!": mapDelForm,
		"map-size": mapSizeForm,

		"stream":             streamForm,
		"stream-set-sink!":   streamSetSink,
		"stream-set-source!": streamSetSource,
		"->":                 sourceForm,
		"<-":                 sinkForm,
	}

	fr := vm.Frame
	for name, evaler := range builtins {
		loadDefaultForm(vm, fr, name, evaler)
	}
}

func loadDefaultForm(vm *Vm, fr *Frame, name string, evaler formEvaler) {
	fr.Put(name, DefaultFormValue{
		name:   name,
		evaler: evaler,
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
	ss := make([]string, len(e.args))
	for i, n := range e.args {
		ss[i] = n.String()
	}
	return fmt.Sprintf("Mismatched arguments: %s", strings.Join(ss, " "))
}

func addForm(vm *Vm, fr *Frame, args []Value) (Value, error) {
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

	if firstInt, fok := first.(IntValue); fok {
		if _, sok := second.(FracValue); sok {
			first = FracValue(float64(firstInt))
		}
	} else if _, fok := first.(FracValue); fok {
		if secondInt, sok := second.(IntValue); sok {
			second = FracValue(float64(secondInt))
		}
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

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func subtractForm(vm *Vm, fr *Frame, args []Value) (Value, error) {
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

	if firstInt, fok := first.(IntValue); fok {
		if _, sok := second.(FracValue); sok {
			first = FracValue(float64(firstInt))
		}
	} else if _, fok := first.(FracValue); fok {
		if secondInt, sok := second.(IntValue); sok {
			second = FracValue(float64(secondInt))
		}
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

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func multiplyForm(vm *Vm, fr *Frame, args []Value) (Value, error) {
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

	if firstInt, fok := first.(IntValue); fok {
		if _, sok := second.(FracValue); sok {
			first = FracValue(float64(firstInt))
		}
	} else if _, fok := first.(FracValue); fok {
		if secondInt, sok := second.(IntValue); sok {
			second = FracValue(float64(secondInt))
		}
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

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func divideForm(vm *Vm, fr *Frame, args []Value) (Value, error) {
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

	if firstInt, fok := first.(IntValue); fok {
		if _, sok := second.(FracValue); sok {
			first = FracValue(float64(firstInt))
		}
	} else if _, fok := first.(FracValue); fok {
		if secondInt, sok := second.(IntValue); sok {
			second = FracValue(float64(secondInt))
		}
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

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func andForm(vm *Vm, fr *Frame, args []Value) (Value, error) {
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

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return cleanFirst & cleanSecond, nil
		}
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func orForm(vm *Vm, fr *Frame, args []Value) (Value, error) {
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

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return cleanFirst | cleanSecond, nil
		}
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func xorForm(vm *Vm, fr *Frame, args []Value) (Value, error) {
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

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return cleanFirst ^ cleanSecond, nil
		}
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func greaterForm(vm *Vm, fr *Frame, args []Value) (Value, error) {
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

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			if cleanFirst > cleanSecond {
				return IntValue(1), nil
			} else {
				return IntValue(0), nil
			}
		}
	case StringValue:
		if cleanSecond, ok := second.(StringValue); ok {
			cmp := strings.Compare(string(cleanFirst), string(cleanSecond))
			if cmp == 1 {
				return IntValue(1), nil
			} else {
				return IntValue(0), nil
			}
		}
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}
func lessForm(vm *Vm, fr *Frame, args []Value) (Value, error) {
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

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			if cleanFirst < cleanSecond {
				return IntValue(1), nil
			} else {
				return IntValue(0), nil
			}
		}
	case StringValue:
		if cleanSecond, ok := second.(StringValue); ok {
			cmp := strings.Compare(string(cleanFirst), string(cleanSecond))
			if cmp == -1 {
				return IntValue(1), nil
			} else {
				return IntValue(0), nil
			}
		}
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func equalForm(vm *Vm, fr *Frame, args []Value) (Value, error) {
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
		return IntValue(0), nil
	}
}
