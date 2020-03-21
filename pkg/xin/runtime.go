package xin

import (
	"bufio"
	"io"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

type formEvaler func(*Frame, []Value, *astNode) (Value, InterpreterError)

type NativeFormValue struct {
	name   string
	evaler formEvaler
}

func (v NativeFormValue) String() string {
	return "(<native form> " + v.name + ")"
}

func (v NativeFormValue) Repr() string {
	return v.String()
}

func (v NativeFormValue) Equal(o Value) bool {
	if ov, ok := o.(NativeFormValue); ok {
		return v.name == ov.name
	}

	return false
}

func loadAllDefaultValues(vm *Vm) {
	fr := vm.Frame

	stdoutStream := NewStream()
	stdoutStream.callbacks.sink = func(v Value) InterpreterError {
		os.Stdout.Write([]byte(v.String()))
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

func loadAllNativeForms(vm *Vm) {
	// seed PRNG for math::rand
	rand.Seed(time.Now().UTC().UnixNano())

	vm.evalers = map[string]formEvaler{
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
		"str::set!":  strSetForm,
		"str::add!":  strAddForm,
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
		"map::has?": mapHasForm,
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

		"math::sin":    mathSinForm,
		"math::cos":    mathCosForm,
		"math::tan":    mathTanForm,
		"math::asin":   mathAsinForm,
		"math::acos":   mathAcosForm,
		"math::atan":   mathAtanForm,
		"math::ln":     mathLnForm,
		"math::rand":   mathRandForm,
		"crypto::rand": cryptoRandForm,

		"os::wait":   osWaitForm,
		"os::open":   osOpenForm,
		"os::delete": osDeleteForm,
		"os::dial":   osDialForm,
		"os::listen": osListenForm,
		"os::args":   osArgsForm,
		"os::time":   osTimeForm,

		"debug::dump": debugDumpForm,
	}

	fr := vm.Frame
	for name, evaler := range vm.evalers {
		fr.Put(name, NativeFormValue{
			name:   name,
			evaler: evaler,
		})
	}
}

func addForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 2 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 2,
			given:    len(args),
		}
	}

	first, second := args[0], args[1]

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
			return StringValue(append(append(base, cleanFirst...), cleanSecond...)), nil
		}
	case VecValue:
		if cleanSecond, ok := second.(VecValue); ok {
			base := make([]Value, 0, len(cleanFirst.underlying.items)+len(cleanSecond.underlying.items))
			return NewVecValue(append(append(base, cleanFirst.underlying.items...), cleanSecond.underlying.items...)), nil
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

	first, second := args[0], args[1]

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

	first, second := args[0], args[1]

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

	first, second := args[0], args[1]

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
			if cleanSecond == zeroValue {
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

	first, second := args[0], args[1]

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
			if cleanSecond == zeroValue {
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

	first, second := args[0], args[1]

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

	first := args[0]

	if firstInt, ok := first.(IntValue); ok {
		if firstInt.Equal(zeroValue) {
			return trueValue, nil
		} else {
			return falseValue, nil
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

	first, second := args[0], args[1]

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

	first, second := args[0], args[1]

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

	first, second := args[0], args[1]

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

	first, second := args[0], args[1]

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
			if cleanFirst > cleanSecond {
				return trueValue, nil
			} else {
				return falseValue, nil
			}
		}
	case FracValue:
		if cleanSecond, ok := second.(FracValue); ok {
			if cleanFirst > cleanSecond {
				return trueValue, nil
			} else {
				return falseValue, nil
			}
		}
	case StringValue:
		if cleanSecond, ok := second.(StringValue); ok {
			cmp := strings.Compare(string(cleanFirst), string(cleanSecond))
			if cmp == 1 {
				return trueValue, nil
			} else {
				return falseValue, nil
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

	first, second := args[0], args[1]

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
			if cleanFirst < cleanSecond {
				return trueValue, nil
			} else {
				return falseValue, nil
			}
		}
	case FracValue:
		if cleanSecond, ok := second.(FracValue); ok {
			if cleanFirst < cleanSecond {
				return trueValue, nil
			} else {
				return falseValue, nil
			}
		}
	case StringValue:
		if cleanSecond, ok := second.(StringValue); ok {
			cmp := strings.Compare(string(cleanFirst), string(cleanSecond))
			if cmp == -1 {
				return trueValue, nil
			} else {
				return falseValue, nil
			}
		}
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}
