package xin

type StreamValue chan Value

func (v StreamValue) String() string {
	return "(<stream>)"
}

func (v StreamValue) Equal(o Value) bool {
	return v == o
}

func streamForm(fr *Frame, args []Value) (Value, error) {
	if len(args) != 0 {
		return nil, IncorrectNumberOfArgsError{
			required: 0,
			given:    len(args),
		}
	}

	return make(StreamValue), nil
}

func sourceForm(fr *Frame, args []Value) (Value, error) {
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

	if isFirstStream, ok := first.(StreamValue); ok {
		val := <-isFirstStream
		return val, nil
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func sinkForm(fr *Frame, args []Value) (Value, error) {
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

	if isFirstStream, ok := first.(StreamValue); ok {
		isFirstStream <- second
		return second, nil
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}
