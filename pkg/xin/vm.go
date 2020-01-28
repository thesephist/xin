package xin

import (
	"io"
	"sync"
)

type Vm struct {
	Frame *Frame

	sync.Mutex
	waiter sync.WaitGroup
}

func NewVm() *Vm {
	vm := &Vm{
		Frame: newFrame(nil), // no parent frame
	}
	vm.Frame.Vm = vm

	loadAllDefaultValues(vm)
	loadAllDefaultForms(vm)
	loadStandardLibrary(vm)

	return vm
}

func (vm *Vm) Eval(r io.Reader) (Value, InterpreterError) {
	defer vm.waiter.Wait()

	toks, err := lex(r)
	if err != nil {
		return nil, err
	}
	rootNode, err := parse(toks)
	if err != nil {
		return nil, err
	}

	vm.Lock()
	defer vm.Unlock()

	val, err := unlazyEval(vm.Frame, &rootNode)
	if err != nil {
		return nil, err
	}

	return val, nil
}
