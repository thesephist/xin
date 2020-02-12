package xin

import (
	"fmt"
	"io"
	"os"
	osPath "path"
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
		Frame:   newFrame(nil), // no parent frame
		imports: make(map[string]*Frame),
	}
	vm.Frame.Vm = vm

	cwd, osErr := os.Getwd()
	if osErr != nil {
		return nil, RuntimeError{
			reason: "Cannot find working directory",
		}
	}
	vm.Frame.cwd = &cwd

	loadAllDefaultValues(vm)
	loadAllNativeForms(vm)
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

	val, err := unlazyEval(vm.Frame, &rootNode)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (vm *Vm) Exec(path string) InterpreterError {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return RuntimeError{
			reason: fmt.Sprintf("Error opening file: %s", err),
		}
	}

	cwd := osPath.Dir(path)
	vm.Frame.cwd = &cwd
	_, ierr := vm.Eval(path, file)
	if ierr != nil {
		return ierr
	}

	return nil
}
