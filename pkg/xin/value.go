package xin

import (
	"fmt"
	"strconv"
	"strings"
)

// Value represents a value in the Xin runtime.
//
// It is important that Value types are:
// (1) freely copyable without losing information
// 	(i.e. copying Values around should not alter language semantics)
// (2) hashable, for use as a MapValue key. StringValue is a special
// 	exception to this case, where a proxy Value type is used instead
// 	which is hashable.
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

type argList []string

type FormValue struct {
	frame *Frame
	// this level of indirection is to allow FormValue
	// to be hashable for inclusion in a MapValue
	arguments  *argList
	definition *astNode
}

func (v FormValue) String() string {
	return "(<form> " + strings.Join(*v.arguments, " ") + ") " + v.definition.String()
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
	// hot path, a shortcut for frequent case
	if _, isLazy := v.(LazyValue); !isLazy {
		return v, nil
	}

	var err InterpreterError
	for lzv, isLazy := v.(LazyValue); isLazy; lzv, isLazy = v.(LazyValue) {
		v, err = eval(lzv.frame, lzv.node)
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

func unlazyEval(fr *Frame, node *astNode) (Value, InterpreterError) {
	val, err := eval(fr, node)
	if err != nil {
		return nil, err
	}

	return unlazy(val)
}
