package xin

var streamId int64 = 0

type sinkCallback func(Value) error

type sourceCallback func() (Value, error)

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
	return "(" + streamType + "<stream>)"
}

func (v StreamValue) Equal(o Value) bool {
	if ov, ok := o.(StreamValue); ok {
		return v.id == ov.id
	}

	return false
}

func streamForm(vm *Vm, fr *Frame, args []Value) (Value, error) {
	if len(args) != 0 {
		return nil, IncorrectNumberOfArgsError{
			required: 0,
			given:    len(args),
		}
	}

	return StreamValue{}, nil
}

type InvalidStreamCallbackError struct {
	reason string
}

func (e InvalidStreamCallbackError) Error() string {
	return "Invalid stream callback:" + e.reason
}

func streamSetSink(vm *Vm, fr *Frame, args []Value) (Value, error) {
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

			firstStream.sink = func(v Value) error {
				localFrame := newFrame(fr)
				localFrame.Put(secondForm.arguments[0], v)

				vm.Unlock()
				defer vm.Lock()

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

func streamSetSource(vm *Vm, fr *Frame, args []Value) (Value, error) {
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

			firstStream.source = func() (Value, error) {
				localFrame := newFrame(fr)

				vm.Unlock()
				defer vm.Lock()

				return eval(localFrame, secondForm.definition)
			}

			return secondForm, nil
		}
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func sourceForm(vm *Vm, fr *Frame, args []Value) (Value, error) {
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

func sinkForm(vm *Vm, fr *Frame, args []Value) (Value, error) {
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
				reason: "Cannot try to sink from a non-sink stream",
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
