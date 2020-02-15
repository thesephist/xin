package xin

import "fmt"

var streamId int64 = 0

type sinkCallback func(Value) InterpreterError

type sourceCallback func() (Value, InterpreterError)

type closerCallback func() InterpreterError

type streamCallbacks struct {
	sink   sinkCallback
	source sourceCallback
	closer closerCallback
}

type StreamValue struct {
	// id is used to compare and de-duplicate stream values in memory
	id        int64
	callbacks *streamCallbacks
}

func NewStream() StreamValue {
	streamId++
	return StreamValue{
		id:        streamId,
		callbacks: &streamCallbacks{},
	}
}

func (v StreamValue) isSink() bool {
	return v.callbacks.sink != nil
}

func (v StreamValue) isSource() bool {
	return v.callbacks.source != nil
}

func (v StreamValue) isClose() bool {
	return v.callbacks.closer != nil
}

func (v StreamValue) String() string {
	streamType := ""
	if v.isSink() {
		streamType += "sink "
	}
	if v.isSource() {
		streamType += "source "
	}
	if v.isClose() {
		streamType += "close "
	}
	return fmt.Sprintf("(%s<stream %d>)", streamType, v.id)
}

func (v StreamValue) Equal(o Value) bool {
	if ov, ok := o.(StreamValue); ok {
		return v.id == ov.id
	}

	return false
}

func streamForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	return NewStream(), nil
}

func streamSetSink(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 2 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 2,
			given:    len(args),
		}
	}

	first, second := args[0], args[1]

	if firstStream, ok := first.(StreamValue); ok {
		if secondForm, ok := second.(FormValue); ok {
			if len(*secondForm.arguments) != 1 {
				return nil, InvalidStreamCallbackError{
					reason: "Mismatched argument count in callback",
				}
			}

			firstStream.callbacks.sink = func(v Value) InterpreterError {
				localFrame := newFrame(fr)
				localFrame.Put((*secondForm.arguments)[0], v)

				_, err := unlazyEval(localFrame, secondForm.definition)
				return err
			}

			return secondForm, nil
		}
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func streamSetSource(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 2 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 2,
			given:    len(args),
		}
	}

	first, second := args[0], args[1]

	if firstStream, ok := first.(StreamValue); ok {
		if secondForm, ok := second.(FormValue); ok {
			if len(*secondForm.arguments) != 0 {
				return nil, InvalidStreamCallbackError{
					reason: "Mismatched argument count in callback",
				}
			}

			firstStream.callbacks.source = func() (Value, InterpreterError) {
				localFrame := newFrame(fr)
				return eval(localFrame, secondForm.definition)
			}

			return secondForm, nil
		}
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func streamSetClose(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 2 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 2,
			given:    len(args),
		}
	}

	first, second := args[0], args[1]

	if firstStream, ok := first.(StreamValue); ok {
		if secondForm, ok := second.(FormValue); ok {
			if len(*secondForm.arguments) != 0 {
				return nil, InvalidStreamCallbackError{
					reason: "Mismatched argument count in callback",
				}
			}

			firstStream.callbacks.closer = func() InterpreterError {
				localFrame := newFrame(fr)
				eval(localFrame, secondForm.definition)
				return nil
			}

			return secondForm, nil
		}
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func streamSourceForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 1,
			given:    len(args),
		}
	}

	first := args[0]

	if firstStream, ok := first.(StreamValue); ok {
		if !firstStream.isSource() {
			return nil, InvalidStreamCallbackError{
				reason: "Cannot try to source from a non-source stream",
			}
		}

		return firstStream.callbacks.source()
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func streamSinkForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 2 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 2,
			given:    len(args),
		}
	}

	first, second := args[0], args[1]

	if firstStream, ok := first.(StreamValue); ok {
		if !firstStream.isSink() {
			return nil, InvalidStreamCallbackError{
				reason: "Cannot try to sink to a non-sink stream",
			}
		}

		err := firstStream.callbacks.sink(second)
		if err != nil {
			return nil, err
		}
		return second, nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func streamCloseForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 1,
			given:    len(args),
		}
	}

	first := args[0]

	if firstStream, ok := first.(StreamValue); ok {
		if !firstStream.isClose() {
			return nil, InvalidStreamCallbackError{
				reason: "Cannot try to close to a non-close stream",
			}
		}

		err := firstStream.callbacks.closer()
		if err != nil {
			return nil, err
		}
		return firstStream, nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}
