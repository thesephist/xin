package xin

import (
	"fmt"
	"strconv"
	"strings"
)

type Value interface {
	String() string
}

type StringValue string

func (v StringValue) String() string {
	return "'" + string(v) + "'"
}

type IntValue int64

func (v IntValue) String() string {
	return strconv.FormatInt(int64(v), 10)
}

type FracValue float64

func (v FracValue) String() string {
	return fmt.Sprintf("%.8f", float64(v))
}

type FormValue struct {
	arguments  []string
	definition *astNode
}

func (v FormValue) String() string {
	return "(<form> " + strings.Join(v.arguments, " ") + ") " + v.definition.String()
}

type VecValue []Value

func (v VecValue) String() string {
	ss := make([]string, len(v))
	for i, item := range v {
		ss[i] = item.String()
	}
	return "(<vec> " + strings.Join(ss, " ") + ")"
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

type LazyValue struct {
	node *astNode
}

func (v LazyValue) String() string {
	return "(<lazy> " + v.node.String() + ")"
}

func unlazy(fr *Frame, v Value) (Value, error) {
	if asLazy, isLazy := v.(LazyValue); isLazy {
		tmp, err := eval(fr, asLazy.node)
		if err != nil {
			return nil, err
		}
		return unlazy(fr, tmp)
	}

	return v, nil
}

func unlazyEval(fr *Frame, node *astNode) (Value, error) {
	asLazy, err := eval(fr, node)
	if err != nil {
		return nil, err
	}
	return unlazy(fr, asLazy)
}
