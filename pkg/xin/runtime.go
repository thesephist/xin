package xin

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

type formEvaler func(*Frame, []Value, *astNode) (Value, InterpreterError)

type DefaultFormValue struct {
	name   string
	evaler formEvaler
}

func (v DefaultFormValue) eval(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	return v.evaler(fr, args, node)
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
		"%": modForm,
		"^": powForm,

		">": greaterForm,
		"<": lessForm,
		"=": equalForm,

		"!":   notForm,
		"&":   andForm,
		"|":   orForm,
		"xor": xorForm,

		"int":  intForm,
		"frac": fracForm,
		"str":  stringForm,
		"type": typeForm,

		"str::get":   strGetForm,
		"str::size":  strSizeForm,
		"str::slice": strSliceForm,
		"str::enc":   strEncForm,
		"str::dec":   strDecForm,

		"vec":        vecForm,
		"vec::get":   vecGetForm,
		"vec::set!":  vecSetForm,
		"vec::add!":  vecAddForm,
		"vec::size":  vecSizeForm,
		"vec::slice": vecSliceForm,

		"map":       mapForm,
		"map::get":  mapGetForm,
		"map::set!": mapSetForm,
		"map::has":  mapHasForm,
		"map::del!": mapDelForm,
		"map::size": mapSizeForm,
		"map::keys": mapKeysForm,

		"stream":              streamForm,
		"stream::set-sink!":   streamSetSink,
		"stream::set-source!": streamSetSource,
		"stream::set-close!":  streamSetClose,
		"->":                  streamSourceForm,
		"<-":                  streamSinkForm,
		"stream::close!":      streamCloseForm,

		"os::wait":   osWaitForm,
		"os::read":   osReadForm,
		"os::write":  osWriteForm,
		"os::delete": osDeleteForm,
		"os::args":   osArgsForm,

		"debug::dump": debugDumpForm,
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

func addForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
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
			// In this context, strings are immutable. i.e. concatenating
			// strings should produce a completely new string whose modifications
			// won't be observable by the original strings.
			base := make([]byte, 0, len(cleanFirst)+len(cleanSecond))
			base = append(base, cleanFirst...)
			return StringValue(append(base, cleanSecond...)), nil
		}
	case VecValue:
		if cleanSecond, ok := second.(VecValue); ok {
			return NewVecValue(append(cleanFirst.underlying.items, cleanSecond.underlying.items...)), nil
		}
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func subtractForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
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
		node: node,
		args: args,
	}
}

func multiplyForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
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

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return cleanFirst * cleanSecond, nil
		}
	case FracValue:
		if cleanSecond, ok := second.(FracValue); ok {
			return cleanFirst * cleanSecond, nil
		}
	case StringValue:
		if cleanSecond, ok := second.(IntValue); ok {
			max := int(cleanSecond)
			result, iter := "", string(cleanFirst)
			for i := 0; i < max; i++ {
				result += iter
			}
			return StringValue(result), nil
		}
	case VecValue:
		if cleanSecond, ok := second.(IntValue); ok {
			max := int(cleanSecond)
			result := make([]Value, 0, max*len(cleanFirst.underlying.items))
			copy(result, cleanFirst.underlying.items)
			for i := 0; i < max; i++ {
				result = append(result, cleanFirst.underlying.items...)
			}
			return NewVecValue(result), nil
		}
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func divideForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
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

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			if cleanSecond == IntValue(0) {
				return zeroValue, nil
			}

			return cleanFirst / cleanSecond, nil
		}
	case FracValue:
		if cleanSecond, ok := second.(FracValue); ok {
			if cleanSecond == FracValue(0) {
				return zeroValue, nil
			}

			return cleanFirst / cleanSecond, nil
		}
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func modForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
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

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			if cleanSecond == IntValue(0) {
				return zeroValue, nil
			}

			return cleanFirst % cleanSecond, nil
		}
	case FracValue:
		if cleanSecond, ok := second.(FracValue); ok {
			if cleanSecond == FracValue(0) {
				return zeroValue, nil
			}

			modulus := math.Mod(
				float64(cleanFirst),
				float64(cleanSecond),
			)
			return FracValue(modulus), nil
		}
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func powForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
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

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return IntValue(math.Pow(float64(cleanFirst), float64(cleanSecond))), nil
		}
	case FracValue:
		if cleanSecond, ok := second.(FracValue); ok {
			power := math.Pow(
				float64(cleanFirst),
				float64(cleanSecond),
			)
			return FracValue(power), nil
		}
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func notForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
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

	if firstInt, ok := first.(IntValue); ok {
		if firstInt.Equal(zeroValue) {
			return IntValue(1), nil
		} else {
			return zeroValue, nil
		}
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func andForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
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

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return cleanFirst & cleanSecond, nil
		}
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func orForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
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

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return cleanFirst | cleanSecond, nil
		}
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func xorForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
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

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			return cleanFirst ^ cleanSecond, nil
		}
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func greaterForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
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

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			if cleanFirst > cleanSecond {
				return IntValue(1), nil
			} else {
				return zeroValue, nil
			}
		}
	case StringValue:
		if cleanSecond, ok := second.(StringValue); ok {
			cmp := strings.Compare(string(cleanFirst), string(cleanSecond))
			if cmp == 1 {
				return IntValue(1), nil
			} else {
				return zeroValue, nil
			}
		}
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}
func lessForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
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

	switch cleanFirst := first.(type) {
	case IntValue:
		if cleanSecond, ok := second.(IntValue); ok {
			if cleanFirst < cleanSecond {
				return IntValue(1), nil
			} else {
				return zeroValue, nil
			}
		}
	case StringValue:
		if cleanSecond, ok := second.(StringValue); ok {
			cmp := strings.Compare(string(cleanFirst), string(cleanSecond))
			if cmp == -1 {
				return IntValue(1), nil
			} else {
				return zeroValue, nil
			}
		}
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}
