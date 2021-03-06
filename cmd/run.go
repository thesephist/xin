package cmd

import (
	"github.com/fatih/color"
	"github.com/thesephist/xin/pkg/xin"
)

func run(path string) {
	vm, err := xin.NewVm()
	if err != nil {
		color.Red("Error creating Xin VM: %s\n", xin.FormatError(err))
		return
	}

	err = vm.Exec(path)
	if err != nil {
		color.Red("Error: %s\n", xin.FormatError(err))
	}
}
