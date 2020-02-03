package xin

import (
	"fmt"
	"io"
	"sync"
)

// TODO: implement call stack unwinding on interpreter error
type stackRecord struct {
	parent *stackRecord
	node   *astNode
}

func (sr stackRecord) String() string {
	if sr.parent == nil {
		return fmt.Sprintf("%s at %s",
			sr.node.String()[0:32],
			sr.node.position)
	}

	return fmt.Sprintf("%s at %s\nfrom %s",
		sr.node.String()[0:32],
		sr.node.position,
		sr.parent)
}

type Vm struct {
	Frame *Frame

	stack   *stackRecord
	imports map[string]*Frame

	sync.Mutex
	waiter sync.WaitGroup
}

func NewVm() (*Vm, InterpreterError) {
	vm := &Vm{
		Frame: newFrame(nil), // no parent frame
	}
	vm.Frame.Vm = vm

	loadAllDefaultValues(vm)
	loadAllDefaultForms(vm)
	err := loadStandardLibrary(vm)
	if err != nil {
		return nil, err
	}

	return vm, nil
}

func (vm *Vm) pushStack(node *astNode) {
	vm.stack = &stackRecord{
		parent: vm.stack.parent,
		node:   node,
	}
}

func (vm *Vm) popStack() {
	if vm.stack == nil {
		panic("Attempted to unwind an empty call stack!")
	}

	vm.stack = vm.stack.parent
}

func (vm *Vm) Eval(path string, r io.Reader) (Value, InterpreterError) {
	defer vm.waiter.Wait()

	toks, err := lex(path, r)
	if err != nil {
		return nil, err
	}
	rootNode, err := parse(toks)
	if err != nil {
		return nil, err
	}

	vm.Lock()
	defer vm.Unlock()

	// every file runs in a child frame of the global
	// root frame
	fr := newFrame(vm.Frame)
	val, err := unlazyEval(fr, &rootNode)
	if err != nil {
		return nil, err
	}

	return val, nil
}
