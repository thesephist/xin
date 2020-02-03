package xin

import (
	"fmt"
	"os"

	"github.com/rakyll/statik/fs"
	_ "github.com/thesephist/xin/statik"
)

func loadStandardLibrary(vm *Vm) InterpreterError {
	statikFs, err := fs.New()
	if err != nil {
		fmt.Println("Standard library error:", err.Error())
	}

	libFiles := []string{
		"std",
		"math",
		"vec",
		"str",
	}

	vm.Lock()
	defer vm.Unlock()

	for _, path := range libFiles {
		alias := path
		// everything in std.xin is assumed not to be prefixed
		if path == "std" {
			alias = ""
		}

		libFile, ferr := statikFs.Open("/" + path + ".xin")
		if ferr != nil {
			return RuntimeError{
				reason: fmt.Sprintf("Stdlib error loading %s: %s", path, ferr.Error()),
			}
		}
		defer libFile.Close()

		// std import
		toks, err := lex(fmt.Sprintf("(std) %s", alias), libFile)
		if err != nil {
			return err
		}
		rootNode, err := parse(toks)
		if err != nil {
			return err
		}

		libFrame := newFrame(vm.Frame)
		_, err = unlazyEval(libFrame, &rootNode)
		if err != nil {
			return err
		}

		rootFrame := vm.Frame
		if alias == "" {
			for name, value := range libFrame.Scope {
				rootFrame.Put(name, value)
			}
		} else {
			for name, value := range libFrame.Scope {
				rootFrame.Put(alias+"::"+name, value)
			}
		}
	}

	return nil
}

func evalImportForm(fr *Frame, args []*astNode) (Value, InterpreterError) {
	if len(args) == 0 {
		return nil, InvalidImportError{nodes: args}
	}

	pathNode := args[0]
	path, err := unlazyEval(fr, pathNode)
	if err != nil {
		return nil, err
	}

	cleanPath, ok := path.(StringValue)
	if !ok {
		return nil, InvalidImportError{nodes: args}
	}

	importPath := string(cleanPath) + ".xin"

	importFile, osErr := os.Open(importPath)
	if osErr != nil {
		return nil, RuntimeError{
			reason:   fmt.Sprintf("Could not open imported file %s", importPath),
			position: pathNode.position,
		}
	}
	toks, err := lex(importPath, importFile)
	if err != nil {
		return nil, err
	}
	rootNode, err := parse(toks)
	if err != nil {
		return nil, err
	}

	// import runs in a new top-level frame in
	// the same VM (execution lock)
	importFrame := newFrame(fr.Vm.Frame)
	_, err = unlazyEval(importFrame, &rootNode)
	if err != nil {
		return nil, err
	}

	alias := ""
	if len(args) > 1 {
		aliasNode := args[1]
		if aliasNode.token.kind == tkName {
			alias = aliasNode.token.value
		}
	}

	if alias == "" {
		for name, value := range importFrame.Scope {
			fr.Put(name, value)
		}
	} else {
		for name, value := range importFrame.Scope {
			fr.Put(alias+"::"+name, value)
		}
	}

	return zeroValue, nil
}
