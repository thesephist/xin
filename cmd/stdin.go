package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/thesephist/xin/pkg/xin"
)

func stdin() {
	vm, err := xin.NewVm()
	if err != nil {
		color.Red("Error creating Xin VM: %s\n", xin.FormatError(err))
		return
	}

	_, err = vm.Eval("stdin", os.Stdin)
	if err != nil {
		color.Red("Error: %s\n", xin.FormatError(err))
	}
}
