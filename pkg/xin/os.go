package xin

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const readBufferSize = 4096

func osWaitForm(fr *Frame, args []Value) (Value, InterpreterError) {
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

	var duration float64
	if firstInt, fok := first.(IntValue); fok {
		duration = float64(int64(firstInt))
	} else if firstFrac, fok := first.(FracValue); fok {
		duration = float64(firstFrac)
	} else {
		return nil, MismatchedArgumentsError{
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

		_, err := unlazy(args[1])
		if err != nil {
			fmt.Println("Eval error:", FormatError(err))
			return
		}
	}()

	return IntValue(1), nil
}

func osReadForm(fr *Frame, args []Value) (Value, InterpreterError) {
	if len(args) != 1 {
		return nil, IncorrectNumberOfArgsError{
			required: 1,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}

	if firstStr, ok := first.(StringValue); ok {
		file, err := os.Open(string(firstStr))
		if err != nil {
			return IntValue(0), nil
		}

		fileStream := NewStream()
		reader := bufio.NewReader(file)
		closed := false

		fileStream.callbacks.source = func() (Value, InterpreterError) {
			if closed {
				return IntValue(0), nil
			}

			buffer := make([]byte, readBufferSize)
			readBytes, err := reader.Read(buffer)
			if err != nil {
				return IntValue(0), nil
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
		args: args,
	}
}

func osWriteForm(fr *Frame, args []Value) (Value, InterpreterError) {
	if len(args) != 1 {
		return nil, IncorrectNumberOfArgsError{
			required: 1,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}

	if firstStr, ok := first.(StringValue); ok {
		flag := os.O_CREATE | os.O_WRONLY
		file, err := os.OpenFile(string(firstStr), flag, 0644)
		if err != nil {
			return IntValue(0), nil
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
		args: args,
	}
}

func osDeleteForm(fr *Frame, args []Value) (Value, InterpreterError) {
	if len(args) != 1 {
		return nil, IncorrectNumberOfArgsError{
			required: 1,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}

	if firstStr, ok := first.(StringValue); ok {
		err := os.RemoveAll(string(firstStr))
		if err != nil {
			return IntValue(0), nil
		}

		return IntValue(1), nil
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func osDumpForm(fr *Frame, args []Value) (Value, InterpreterError) {
	if len(args) != 0 {
		return nil, IncorrectNumberOfArgsError{
			required: 0,
			given:    len(args),
		}
	}

	return StringValue(fr.String()), nil
}
