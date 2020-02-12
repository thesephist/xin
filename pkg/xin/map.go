package xin

// StringValue is of type []byte, which is unhashable
// so we convert StringValues to hashable string type
// before using as map keys
type hashableStringProxy string

func (v hashableStringProxy) String() string {
	return "(<string proxy> " + string(v) + ")"
}

func (v hashableStringProxy) Equal(ov Value) bool {
	panic("hashableStringProxy should not be equality-checked")
}

// hashableNativeFormProxy is a hashable proxy for
// NativeFormValue
type hashableNativeFormProxy string

func (v hashableNativeFormProxy) String() string {
	return "(<native form proxy> " + string(v) + ")"
}

func (v hashableNativeFormProxy) Equal(ov Value) bool {
	panic("hashableNativeFormProxy should not be equality-checked")
}

func hashable(v Value) Value {
	switch val := v.(type) {
	case StringValue:
		return hashableStringProxy(val)
	case NativeFormValue:
		return hashableNativeFormProxy(val.name)
	default:
		return v
	}
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
	return NewMapValue(), nil
}

func mapGetForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 2 {
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

		return zeroValue, nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func mapSetForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 3 {
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
	if len(args) < 2 {
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

		return zeroValue, nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func mapDelForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 2 {
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

		return zeroValue, nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}

func mapSizeForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
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

func mapKeysForm(fr *Frame, args []Value, node *astNode) (Value, InterpreterError) {
	if len(args) < 1 {
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
		keys := []Value{}
		for k, _ := range *firstMap.items {
			keys = append(keys, k)
		}
		return NewVecValue(keys), nil
	}

	return nil, MismatchedArgumentsError{
		node: node,
		args: args,
	}
}
