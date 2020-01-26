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

type StringValue string

func (v StringValue) String() string {
	return "'" + string(v) + "'"
}

func (v StringValue) Equal(o Value) bool {
	if ov, ok := o.(StringValue); ok {
		return strings.Compare(string(v), string(ov)) == 0
	}

	return false
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

type VecValue []Value

func (v VecValue) String() string {
	ss := make([]string, len(v))
	for i, item := range v {
		ss[i] = item.String()
	}
	return "(<vec> " + strings.Join(ss, " ") + ")"
}

func (v VecValue) Equal(o Value) bool {
	if ov, ok := o.(VecValue); ok {
		if len(v) != len(ov) {
			return false
		}

		for i, x := range v {
			if x != ov[i] {
				return false
			}
		}

		return true
	}

	return false
}

type MapValue map[Value]Value

func (v MapValue) String() string {
	i := 0
	ss := make([]string, len(v))
	for k, val := range v {
		ss[i] = k.String() + "->" + val.String()
		i++
	}
	return "(<map> " + strings.Join(ss, " ") + ")"
}

func (v MapValue) Equal(o Value) bool {
	if ov, ok := o.(MapValue); ok {
		if len(v) != len(ov) {
			return false
		}

		for i, x := range v {
			if x != ov[i] {
				return false
			}
		}

		return true
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

func unlazy(v Value) (Value, error) {
	var lzv LazyValue
	var isLazy bool
	var err error

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

func unlazyEval(fr *Frame, node *astNode) (Value, error) {
	asLazy, err := eval(fr, node)
	if err != nil {
		return nil, err
	}
	return unlazy(asLazy)
}
