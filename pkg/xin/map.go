package xin

func mapForm(fr *Frame, args []Value) (Value, error) {
	if len(args) != 0 {
		return nil, IncorrectNumberOfArgsError{
			required: 0,
			given:    len(args),
		}
	}

	return make(MapValue), nil
}

func mapGetForm(fr *Frame, args []Value) (Value, error) {
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

	if isFirstMap, ok := first.(MapValue); ok {
		val, prs := isFirstMap[second]
		if prs {
			return val, nil
		}

		return VecValue{}, nil
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func mapSetForm(fr *Frame, args []Value) (Value, error) {
	if len(args) != 3 {
		return nil, IncorrectNumberOfArgsError{
			required: 3,
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
	third, err := unlazy(args[2])
	if err != nil {
		return nil, err
	}

	if isFirstMap, ok := first.(MapValue); ok {
		isFirstMap[second] = third

		return third, nil
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func mapDelForm(fr *Frame, args []Value) (Value, error) {
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

	if isFirstMap, ok := first.(MapValue); ok {
		val, prs := isFirstMap[second]
		if prs {
			delete(isFirstMap, second)
			return val, nil
		}

		return VecValue{}, nil
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}

func mapSizeForm(fr *Frame, args []Value) (Value, error) {
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

	if isFirstMap, fok := first.(MapValue); fok {
		return IntValue(len(isFirstMap)), nil
	}

	return nil, MismatchedArgumentsError{
		args: args,
	}
}
