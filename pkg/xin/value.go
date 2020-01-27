package xin

import (
	"fmt"
	"strconv"
	"strings"
)

type Value interface {
	String() string
	Equal(Value) bool
}

type IntValue int64

func (v IntValue) String() string {
	return strconv.FormatInt(int64(v), 10)
}

func (v IntValue) Equal(o Value) bool {
	if ov, ok := o.(IntValue); ok {
		return v == ov
	}

	return false
}

type FracValue float64

func (v FracValue) String() string {
	return fmt.Sprintf("%.8f", float64(v))
}

func (v FracValue) Equal(o Value) bool {
	if ov, ok := o.(FracValue); ok {
		return v == ov
	}

	return false
}

type FormValue struct {
	frame      *Frame
	arguments  []string
	definition *astNode
}

func (v FormValue) String() string {
	return "(<form> " + strings.Join(v.arguments, " ") + ") " + v.definition.String()
}

func (v FormValue) Equal(o Value) bool {
	if ov, ok := o.(FormValue); ok {
		return v.definition == ov.definition
	}

	return false
}

type LazyValue struct {
	frame *Frame
	node  *astNode
}

func (v LazyValue) String() string {
	return "(<lazy> " + v.node.String() + ")"
}

func (v LazyValue) Equal(o Value) bool {
	// should never run
	panic("<lazy> value should never be compared!")
}

func unlazy(v Value) (Value, InterpreterError) {
	var lzv LazyValue
	var isLazy bool
	var err InterpreterError

	lzv, isLazy = v.(LazyValue)
	for isLazy {
		v, err = eval(lzv.frame, lzv.node)
		if err != nil {
			return nil, err
		}

		lzv, isLazy = v.(LazyValue)
	}

	return v, nil
}

func unlazyEval(fr *Frame, node *astNode) (Value, InterpreterError) {
	asLazy, err := eval(fr, node)
	if err != nil {
		return nil, err
	}
	return unlazy(asLazy)
}
