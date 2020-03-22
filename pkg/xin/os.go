package xin

import (
	"bufio"
	"fmt"
	"io"
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

func newRWStream(rw io.ReadWriteCloser) StreamValue {
	rwStream := NewStream()
	reader := bufio.NewReader(rw)
	closed := false

	rwStream.callbacks.source = func() (Value, InterpreterError) {
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

	rwStream.callbacks.sink = func(v Value, node *astNode) InterpreterError {
		if closed {
			return nil
		}

		if strVal, ok := v.(StringValue); ok {
			_, err := rw.Write(strVal)
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

	rwStream.callbacks.closer = func() InterpreterError {
		if !closed {
			closed = true
			rw.Close()
		}
		return nil
	}

	return rwStream
}

func osStatForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 1,
			given:    len(args),
		}
	}

	first := args[0]

	if firstStr, ok := first.(StringValue); ok {
		fileStat, err := os.Stat(string(firstStr))
		if err != nil {
			return zeroValue, nil
		}

		statMap := NewMapValue()

		if fileStat.IsDir() {
			statMap.set(StringValue("dir"), trueValue)
		} else {
			statMap.set(StringValue("dir"), falseValue)
		}
		statMap.set(StringValue("name"), StringValue(fileStat.Name()))
		statMap.set(StringValue("size"), IntValue(fileStat.Size()))
		statMap.set(StringValue("mod"), IntValue(fileStat.ModTime().Unix()))

		return statMap, nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
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

		return newRWStream(file), nil
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

			vm.Lock()
			defer vm.Unlock()

			_, intErr := unlazyEvalFormWithArgs(fr, secondForm, []Value{rv}, node)
			if intErr != nil {
				fmt.Println(FormatError(intErr))
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
	if len(args) < 2 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 2,
			given:    len(args),
		}
	}

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

	return newRWStream(conn), nil
}

func osListenForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 3 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 3,
			given:    len(args),
		}
	}

	network, addr, err := validateNetworkArgs(args, node)
	if err != nil {
		return nil, err
	}

	handler := args[2]
	signal := make(chan bool, 1)

	listener, netErr := net.Listen(network, addr)
	if netErr != nil {
		return nil, NetworkError{
			reason:   netErr.Error(),
			position: node.position,
		}
	}

	vm := fr.Vm
	vm.waiter.Add(1)
	go func(l net.Listener) {
		defer vm.waiter.Done()

		for {
			conn, err := l.Accept()
			if err != nil {
				select {
				case <-signal:
					return
				default:
					fmt.Println(err.Error())
				}
			}

			go func(c net.Conn) {
				vm.Lock()
				defer vm.Unlock()

				_, err := unlazyEvalFormWithArgs(fr, handler, []Value{newRWStream(c)}, node)
				if err != nil {
					fmt.Println(FormatError(err))
				}
			}(conn)
		}
	}(listener)

	return NativeFormValue{
		name: "os::listen::close",
		evaler: func(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
			signal <- true
			listener.Close()
			return trueValue, nil
		},
	}, nil
}

func osLogForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 1,
			given:    len(args),
		}
	}

	first := args[0]
	fmt.Println(first.String())

	return first, nil
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
