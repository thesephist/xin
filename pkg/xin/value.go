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

type LazyValue struct {
	node *astNode
}

func (v LazyValue) String() string {
	return "(<lazy> " + v.node.String() + ")"
}

func unlazy(fr *Frame, v Value) (Value, error) {
	if asLazy, isLazy := v.(LazyValue); isLazy {
		return eval(fr, asLazy.node)
	}

	return v, nil
}
