package xin

import (
	"io"
)

type Vm struct {
	Stack []Value
}

func NewVm() *Vm {
	vm := Vm{
		Stack: make([]Value, 8),
	}
	return &vm
}

func (vm *Vm) Eval(io.Reader) Value {
	return Int(2)
}

func (vm *Vm) pop() Value {
	popped := vm.Stack[len(vm.Stack)-1]
	vm.Stack = vm.Stack[0 : len(vm.Stack)-1]
	return popped
}
