package xin

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type formEvaler func(*Frame, []Value) (Value, InterpreterError)

type DefaultFormValue struct {
	name   string
	evaler formEvaler
}

func (v DefaultFormValue) eval(fr *Frame, args []Value) (Value, InterpreterError) {
	return v.evaler(fr, args)
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

func loadAllDefaultValues(vm *Vm) {
	fr := vm.Frame

	stdoutStream := NewStream()
	stdoutStream.callbacks.sink = func(v Value) InterpreterError {
		fmt.Printf(v.String())
		return nil
	}
	fr.Put("os::stdout", stdoutStream)

	stdinStream := NewStream()
	stdinStream.callbacks.source = func() (Value, InterpreterError) {
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
	}
	fr.Put("os::stdin", stdinStream)
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

		"string": stringForm,
		"int":    intForm,
		"frac":   fracForm,

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

		"os::dump": osDumpForm,
		"os::wait": osWaitForm,
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

func addForm(fr *Frame, args []Value) (Value, InterpreterError) {
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

func subtractForm(fr *Frame, args []Value) (Value, InterpreterError) {
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

func multiplyForm(fr *Frame, args []Value) (Value, InterpreterError) {
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

func divideForm(fr *Frame, args []Value) (Value, InterpreterError) {
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

func andForm(fr *Frame, args []Value) (Value, InterpreterError) {
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

func orForm(fr *Frame, args []Value) (Value, InterpreterError) {
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

func xorForm(fr *Frame, args []Value) (Value, InterpreterError) {
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

func greaterForm(fr *Frame, args []Value) (Value, InterpreterError) {
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
func lessForm(fr *Frame, args []Value) (Value, InterpreterError) {
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

func equalForm(fr *Frame, args []Value) (Value, InterpreterError) {
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
