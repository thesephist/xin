package xin

func vecForm(fr *Frame, args []Value) (Value, error) {
	vecValues := make([]Value, len(args))
	for i, a := range args {
		val, err := unlazy(fr, a)
		if err != nil {
			return nil, err
		}
		vecValues[i] = val
	}
	return VecValue(vecValues), nil
}

func vecGetForm(fr *Frame, args []Value) (Value, error) {
	return nil, MismatchedArgumentsError{}
}

func vecSetForm(fr *Frame, args []Value) (Value, error) {
	return nil, MismatchedArgumentsError{}
}
