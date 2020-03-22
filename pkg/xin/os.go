package xin

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

const readBufferSize = 4096

func osWaitForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 2 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 2,
			given:    len(args),
		}
	}

	first, second := args[0], args[1]

	var duration float64
	if firstInt, fok := first.(IntValue); fok {
		duration = float64(int64(firstInt))
	} else if firstFrac, fok := first.(FracValue); fok {
		duration = float64(firstFrac)
	} else {
		return nil, MismatchedArgumentsError{
			node: node,
			args: args,
		}
	}

	vm := fr.Vm
	vm.waiter.Add(1)
	go func() {
		defer vm.waiter.Done()

		time.Sleep(time.Duration(
			int64(duration * float64(time.Second)),
		))

		vm.Lock()
		defer vm.Unlock()

		_, err := unlazyEvalFormValue(fr, &astNode{
			// dummy function invocation astNode
			// to help generate proper error messages
			isForm: true,
			leaves: []*astNode{
				&astNode{
					isForm: false,
				},
			},
		}, second)
		if err != nil {
			fmt.Println("Eval error in os::wait:", FormatError(err))
			return
		}
	}()

	return trueValue, nil
}

func osOpenForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 1,
			given:    len(args),
		}
	}

	first := args[0]

	if firstStr, ok := first.(StringValue); ok {
		flag := os.O_CREATE | os.O_RDWR
		file, err := os.OpenFile(string(firstStr), flag, 0644)
		if err != nil {
			return zeroValue, nil
		}

		fileStream := NewStream()
		reader := bufio.NewReader(file)
		closed := false

		// read
		fileStream.callbacks.source = func() (Value, InterpreterError) {
			if closed {
				return zeroValue, nil
			}

			buffer := make([]byte, readBufferSize)
			readBytes, err := reader.Read(buffer)
			if err != nil {
				return zeroValue, nil
			}

			return StringValue(buffer[:readBytes]), nil
		}
		// write
		fileStream.callbacks.sink = func(v Value) InterpreterError {
			if closed {
				return nil
			}

			if strVal, ok := v.(StringValue); ok {
				_, err := file.Write(strVal)
				if err != nil {
					return nil
				}
				return nil
			}

			return MismatchedArgumentsError{
				node: node,
				args: []Value{v},
			}
		}
		// release handle
		fileStream.callbacks.closer = func() InterpreterError {
			if !closed {
				closed = true
				file.Close()
			}
			return nil
		}
		return fileStream, nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func osDeleteForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 1,
			given:    len(args),
		}
	}

	first := args[0]
	var second Value
	if len(args) >= 2 {
		second = args[1]
	} else {
		second = Noop
	}

	var secondForm Value
	firstStr, fok := first.(StringValue)
	secondForm, sok := second.(FormValue)
	if !sok {
		secondForm, sok = second.(NativeFormValue)
	}

	if fok && sok {
		vm := fr.Vm
		vm.waiter.Add(1)
		go func() {
			defer vm.waiter.Done()

			err := os.RemoveAll(string(firstStr))
			rv := falseValue
			if err != nil {
				rv = trueValue
			}

			switch form := secondForm.(type) {
			case FormValue:
				localFrame := newFrame(form.frame)

				if len(*form.arguments) > 0 {
					localFrame.Put((*form.arguments)[0], rv)
				}

				lv := LazyValue{
					frame: localFrame,
					node:  form.definition,
				}
				_, err := unlazy(lv)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			case NativeFormValue:
				_, err := form.evaler(fr, []Value{rv}, node)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			default:
				err := InvalidFormError{
					position: node.position,
				}
				fmt.Println(err.Error())
			}
		}()

		return zeroValue, nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func validateNetworkArgs(args []Value, node *astNode) (string, string, InterpreterError) {
	if len(args) < 2 {
		return "", "", IncorrectNumberOfArgsError{
			node:     node,
			required: 2,
			given:    len(args),
		}
	}

	first, second := args[0], args[1]

	firstStr, fok := first.(StringValue)
	secondStr, sok := second.(StringValue)

	if !fok || !sok {
		return "", "", MismatchedArgumentsError{
			node: node,
			args: args,
		}
	}

	network := string(firstStr)
	addr := string(secondStr)

	if (network) != "tcp" && (network) != "udp" {
		return "", "", NetworkError{
			reason:   `Network specified to os::dial must be "tcp" or "udp"`,
			position: node.position,
		}
	}

	return network, addr, nil
}

func osDialForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	network, addr, err := validateNetworkArgs(args, node)
	if err != nil {
		return nil, err
	}

	conn, netErr := net.Dial(network, addr)
	if netErr != nil {
		return nil, NetworkError{
			reason:   netErr.Error(),
			position: node.position,
		}
	}
	connStream := NewStream()
	_ = conn

	// TODO: fill these out
	connStream.callbacks.source = func() (Value, InterpreterError) {
		return StringValue(""), nil
	}
	connStream.callbacks.sink = func(v Value) InterpreterError {
		return nil
	}
	connStream.callbacks.closer = func() InterpreterError {
		return nil
	}

	return connStream, nil
}

func osListenForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	network, addr, err := validateNetworkArgs(args, node)
	if err != nil {
		return nil, err
	}

	listener, netErr := net.Listen(network, addr)
	if netErr != nil {
		return nil, NetworkError{
			reason:   netErr.Error(),
			position: node.position,
		}
	}
	connStream := NewStream()
	_ = listener

	// TODO: fill these out
	connStream.callbacks.source = func() (Value, InterpreterError) {
		return StringValue(""), nil
	}
	connStream.callbacks.sink = func(v Value) InterpreterError {
		return nil
	}
	connStream.callbacks.closer = func() InterpreterError {
		return nil
	}

	return connStream, nil
}

func osArgsForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	argsVec := make([]Value, len(os.Args))
	for i, a := range os.Args {
		argsVec[i] = StringValue(a)
	}
	return NewVecValue(argsVec), nil
}

func osTimeForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	return FracValue(
		float64(time.Now().UnixNano()) / 1e9,
	), nil
}

func debugDumpForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	return StringValue(fr.String()), nil
}
