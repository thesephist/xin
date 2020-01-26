package xin

import (
	"fmt"
	"io"
)

type Vm struct {
	Frame *Frame
}

func NewVm() *Vm {
	return &Vm{
		Frame: newFrame(),
	}
}

func (vm *Vm) Eval(r io.Reader) (Value, error) {
	toks := lex(r)
	rootNode, err := parse(toks)
	if err != nil {
		fmt.Printf("There was an error: %s", err.Error())
	}

	return eval(vm.Frame, rootNode)
}
