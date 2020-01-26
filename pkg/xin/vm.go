package xin

import (
	"fmt"
	"io"
)

type Vm struct {
	Frame *Frame
}

func NewVm() *Vm {
	vm := &Vm{
		Frame: newFrame(),
	}
	loadAllDefaultValues(vm.Frame)
	loadAllDefaultForms(vm.Frame)

	return vm
}

func (vm *Vm) Eval(r io.Reader) (Value, error) {
	toks := lex(r)
	rootNode, err := parse(toks)
	if err != nil {
		fmt.Printf("There was an error: %s", err.Error())
	}

	return unlazyEval(vm.Frame, &rootNode)
}
