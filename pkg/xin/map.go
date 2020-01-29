package xin

// StringValue is of type []byte, which is unhashable
// so we convert StringValues to hashable string type
// before using as map keys
type hashableStringProxyValue string

func (v hashableStringProxyValue) String() string {
	return "(<string proxy> " + string(v) + ")"
}

func (v hashableStringProxyValue) Equal(ov Value) bool {
	panic("hashableStringProxyValue should not be equality-checked")
}

func hashable(v Value) Value {
	if asStr, ok := v.(StringValue); ok {
		return hashableStringProxyValue(asStr)
	}

	return v
}

type mapItems map[Value]Value

type MapValue struct {
	// indirection used to allow MapValue to be hashable
	// and correctly equality-checked
	items *mapItems
}

func (v MapValue) String() string {
	ss := ""
	for k, val := range *v.items {
		ss += " " + k.String() + "->" + val.String()
	}
	return "(<map>" + ss + ")"
}

func (v MapValue) Equal(o Value) bool {
	if ov, ok := o.(MapValue); ok {
		return v.items == ov.items
	}

	return false
}

func NewMapValue() MapValue {
	return MapValue{
		items: &mapItems{},
	}
}

func mapForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) != 0 {
		return nil, IncorrectNumberOfArgsError{
			node:     node,
			required: 0,
			given:    len(args),
		}
	}

	return NewMapValue(), nil
}

func mapGetForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
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

	if firstMap, ok := first.(MapValue); ok {
		val, prs := (*firstMap.items)[hashable(second)]
		if prs {
			return val, nil
		}

		return IntValue(0), nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func mapSetForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
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

	if firstMap, ok := first.(MapValue); ok {
		(*firstMap.items)[hashable(second)] = third
		return firstMap, nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func mapHasForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
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

	if firstMap, ok := first.(MapValue); ok {
		_, prs := (*firstMap.items)[hashable(second)]
		if prs {
			return IntValue(1), nil
		}

		return IntValue(0), nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func mapDelForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
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

	if firstMap, ok := first.(MapValue); ok {
		_, prs := (*firstMap.items)[hashable(second)]
		if prs {
			delete(*firstMap.items, second)
			return firstMap, nil
		}

		return IntValue(0), nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func mapSizeForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
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

	if firstMap, fok := first.(MapValue); fok {
		return IntValue(len(*firstMap.items)), nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}
