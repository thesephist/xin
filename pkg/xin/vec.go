package xin

// vecUnderlying provides a layer of indirection
// we need to allow vecs to be mutable in-place
// because Go slices are not in-place mutable.
// This also allows VecValue to be hashable for use
// in a MapValue as a key.
type vecUnderlying struct {
	items []Value
}

type VecValue struct {
	underlying *vecUnderlying
}

func NewVecValue(items []Value) VecValue {
	return VecValue{
		underlying: &vecUnderlying{items},
	}
}

func (v VecValue) String() string {
	ss := ""
	for _, item := range v.underlying.items {
		ss += " " + item.String()
	}
	return "(<vec>" + ss + ")"
}

func (v VecValue) Equal(o Value) bool {
	if ov, ok := o.(VecValue); ok {
		return v.underlying == ov.underlying
	}

	return false
}

func vecForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	vecValues := make([]Value, len(args))
	for i, a := range args {
		val, err := unlazy(a)
		if err != nil {
			return nil, err
		}
		vecValues[i] = val
	}
	return NewVecValue(vecValues), nil
}

func vecGetForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) != 2 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
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

	firstVec, fok := first.(VecValue)
	secondInt, sok := second.(IntValue)
	if fok && sok {
		if int(secondInt) < len(firstVec.underlying.items) {
			return firstVec.underlying.items[secondInt], nil
		}

		return zeroValue, nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func vecSetForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) != 3 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
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

	firstVec, fok := first.(VecValue)
	secondInt, sok := second.(IntValue)
	if fok && sok {
		if int(secondInt) < len(firstVec.underlying.items) {
			firstVec.underlying.items[secondInt] = third
			return firstVec, nil
		}

		return zeroValue, nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func vecAddForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) != 2 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
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

	if firstVec, fok := first.(VecValue); fok {
		firstVec.underlying.items = append(firstVec.underlying.items, second)
		return firstVec, nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func vecSizeForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) != 1 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 1,
			given:    len(args),
		}
	}

	first, err := unlazy(args[0])
	if err != nil {
		return nil, err
	}

	if firstVec, fok := first.(VecValue); fok {
		return IntValue(len(firstVec.underlying.items)), nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func vecSliceForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) != 3 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
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

	firstVec, fok := first.(VecValue)
	secondInt, sok := second.(IntValue)
	thirdInt, tok := third.(IntValue)
	if fok && sok && tok {
		max := len(firstVec.underlying.items)
		inRange := func(iv IntValue) bool {
			return int(iv) >= 0 && int(iv) <= max
		}

		if int(secondInt) > max {
			secondInt = IntValue(max)
		}
		if int(thirdInt) > max {
			thirdInt = IntValue(max)
		}

		if inRange(secondInt) && inRange(thirdInt) {
			base := make([]Value, 0, thirdInt-secondInt)
			items := append(base, firstVec.underlying.items[secondInt:thirdInt]...)
			return NewVecValue(items), nil
		}

		return zeroValue, nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}
