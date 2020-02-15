package xin

import (
	"bufio"
	"fmt"
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

	return IntValue(1), nil
}

func osReadForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 1,
			given:    len(args),
		}
	}

	first := args[0]

	if firstStr, ok := first.(StringValue); ok {
		file, err := os.Open(string(firstStr))
		if err != nil {
			return zeroValue, nil
		}

		fileStream := NewStream()
		reader := bufio.NewReader(file)
		closed := false

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

func osWriteForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 1,
			given:    len(args),
		}
	}

	first := args[0]

	if firstStr, ok := first.(StringValue); ok {
		flag := os.O_CREATE | os.O_WRONLY
		file, err := os.OpenFile(string(firstStr), flag, 0644)
		if err != nil {
			return zeroValue, nil
		}

		fileStream := NewStream()
		closed := false

		fileStream.callbacks.sink = func(v Value) InterpreterError {
			if closed {
				// TODO: maybe we should throw on write after close?
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
			return zeroValue, nil
		}

		return IntValue(1), nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func osArgsForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	argsVec := make([]Value, len(os.Args))
	for i, a := range os.Args {
		argsVec[i] = StringValue(a)
	}
	return NewVecValue(argsVec), nil
}

func debugDumpForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	return StringValue(fr.String()), nil
}
