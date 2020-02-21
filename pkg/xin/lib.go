package xin

import (
	"fmt"
	"os"
	osPath "path"

	"github.com/rakyll/statik/fs"
	_ "github.com/thesephist/xin/statik"
)

func loadStandardLibrary(vm *Vm) InterpreterError {
	statikFs, err := fs.New()
	if err != nil {
		fmt.Println("Standard library error:", err.Error())
	}

	// import order matters here, later libs
	// have dependency on the preceding ones
	libFiles := []string{
		"std",
		"math",
		"vec",
		"map",
		"str",
		"src",
		"stat",
		"test",
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

func importFrame(fr *Frame, importPath string) (*Frame, InterpreterError) {
	importFile, osErr := os.Open(importPath)
	if osErr != nil {
		return nil, RuntimeError{
			reason: fmt.Sprintf("Could not open imported file %s", importPath),
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
	importCwd := osPath.Dir(importPath)
	importFrame.cwd = &importCwd
	_, err = unlazyEval(importFrame, &rootNode)
	if err != nil {
		return nil, err
	}

	return importFrame, nil
}

func deduplicatedImportFrame(fr *Frame, importPath string) (*Frame, InterpreterError) {
	importMap := fr.Vm.imports

	if dedupFrame, prs := importMap[importPath]; prs {
		return dedupFrame, nil
	}

	dedupFrame, err := importFrame(fr, importPath)
	if err != nil {
		return nil, err
	}

	importMap[importPath] = dedupFrame

	return dedupFrame, nil
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

	importPath := osPath.Join(*fr.cwd, string(cleanPath)+".xin")
	importFramePtr, err := deduplicatedImportFrame(fr, importPath)
	if err != nil {
		return nil, err
	}
	importFrame := *importFramePtr

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
