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

	if firstStr, ok := first.(StringValue); ok {
		err := os.RemoveAll(string(firstStr))
		if err != nil {
			return falseValue, nil
		}

		return trueValue, nil
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
		return "", "", MismatchedArgumentsError{
			// TODO: make this a more descriptive error
			node: node,
			args: args,
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
		// TODO: make this a more descriptive error
		return nil, RuntimeError{}
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
		// TODO: make this a more descriptive error
		return nil, RuntimeError{}
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
