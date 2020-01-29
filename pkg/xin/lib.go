package xin

import (
	"fmt"

	"github.com/rakyll/statik/fs"
	_ "github.com/thesephist/xin/statik"
)

func loadStandardLibrary(vm *Vm) {
	statikFs, err := fs.New()
	if err != nil {
		fmt.Println("Standard library error:", err.Error())
	}

	libFiles := []string{
		"std.xin",
		"tensor.xin",
	}

	for _, path := range libFiles {
		libFile, err := statikFs.Open("/" + path)
		if err != nil {
			fmt.Println("Standard library read error:", err.Error())
		}
		defer libFile.Close()
		vm.Eval(fmt.Sprintf("<xin> %s", path), libFile)
	}
}
