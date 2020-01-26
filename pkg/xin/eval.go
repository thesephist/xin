package xin

import (
	"fmt"
)

type Frame struct {
	Scope  map[string]Value
	Parent *Frame
}

func newFrame() *Frame {
	return &Frame{
		Scope: make(map[string]Value),
	}
}

type UndefinedNameError struct {
	name string
}

func (e UndefinedNameError) Error() string {
	return fmt.Sprintf("Undefined name %s", e.name)
}

func (fr *Frame) Get(name string) (Value, error) {
	if val, prs := fr.Scope[name]; prs {
		return val, nil
	}
	if fr.Parent != nil {
		return fr.Parent.Get(name)
	}
	return nil, UndefinedNameError{name: name}
}

func (fr *Frame) Put(name string, val Value) {
	fr.Scope[name] = val
}

func eval(fr *Frame, node astNode) (Value, error) {
	return IntValue(42), nil
}
