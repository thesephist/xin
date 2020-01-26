package xin

import (
	"fmt"
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

	vm.Lock()
	return vm
}

func (vm *Vm) Eval(r io.Reader) (Value, error) {
	defer vm.waiter.Wait()

	toks := lex(r)
	rootNode, err := parse(toks)
	if err != nil {
		fmt.Printf("There was an error: %s", err.Error())
	}

	vm.Unlock()
	defer vm.Lock()

	val, err := unlazyEval(vm.Frame, &rootNode)
	if err != nil {
		return nil, err
	}

	return val, nil
}
