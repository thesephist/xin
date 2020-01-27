package xin

import "fmt"

var streamId int64 = 0

type sinkCallback func(Value) InterpreterError

type sourceCallback func() (Value, InterpreterError)

type StreamValue struct {
	// id is used to compare and de-duplicate stream values in memory
	id     int64
	sink   sinkCallback
	source sourceCallback
}

func (v StreamValue) isSink() bool {
	return v.sink != nil
}

func (v StreamValue) isSource() bool {
	return v.source != nil
}

func (v StreamValue) String() string {
	streamType := ""
	if v.isSink() {
		streamType += "sink "
	}
	if v.isSource() {
		streamType += "source "
	}
	return fmt.Sprintf("(%s<stream %d>)", streamType, v.id)
}

func (v StreamValue) Equal(o Value) bool {
	if ov, ok := o.(StreamValue); ok {
		return v.id == ov.id
	}

	return false
}

func streamForm(fr *Frame, args []Value) (Value, InterpreterError) {
	if len(args) != 0 {
		return nil, IncorrectNumberOfArgsError{
			required: 0,
			given:    len(args),
		}
	}

	streamId++

	return StreamValue{
		id: streamId,
	}, nil
}

func streamSetSink(fr *Frame, args []Value) (Value, InterpreterError) {
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

	if firstStream, ok := first.(StreamValue); ok {
		if secondForm, ok := second.(FormValue); ok {
			if len(secondForm.arguments) != 1 {
				return nil, InvalidStreamCallbackError{
					reason: "Mismatched argument count in callback",
				}
			}

			firstStream.sink = func(v Value) InterpreterError {
				localFrame := newFrame(fr)
				localFrame.Put(secondForm.arguments[0], v)

				_, err := eval(localFrame, secondForm.definition)
				return err
			}

			return secondForm, nil
		}
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func streamSetSource(fr *Frame, args []Value) (Value, InterpreterError) {
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

	if firstStream, ok := first.(StreamValue); ok {
		if secondForm, ok := second.(FormValue); ok {
			if len(secondForm.arguments) != 0 {
				return nil, InvalidStreamCallbackError{
					reason: "Mismatched argument count in callback",
				}
			}

			firstStream.source = func() (Value, InterpreterError) {
				localFrame := newFrame(fr)

				return eval(localFrame, secondForm.definition)
			}

			return secondForm, nil
		}
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func sourceForm(fr *Frame, args []Value) (Value, InterpreterError) {
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

	if firstStream, ok := first.(StreamValue); ok {
		if !firstStream.isSource() {
			return nil, InvalidStreamCallbackError{
				reason: "Cannot try to source from a non-source stream",
			}
		}

		return firstStream.source()
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func sinkForm(fr *Frame, args []Value) (Value, InterpreterError) {
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

	if firstStream, ok := first.(StreamValue); ok {
		if !firstStream.isSink() {
			return nil, InvalidStreamCallbackError{
				reason: "Cannot try to sink to a non-sink stream",
			}
		}

		err := firstStream.sink(second)
		if err != nil {
			return nil, err
		}
		return second, nil
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}
