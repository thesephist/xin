package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/thesephist/xin/pkg/xin"
)

func repl() {
	vm, err := xin.NewVm()
	if err != nil {
		color.Red("Error creating Xin VM: %s\n", xin.FormatError(err))
		return
	}

	replCount := 0
	reader := bufio.NewReader(os.Stdin)
	for {
		// TODO: capture SIGINT and just continue loop
		// this requires running the VM in a Context{}

		fmt.Printf("%d ) ", replCount)

		text, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			color.Red("Repl error: %s\n\n", err.Error())
			continue
		}

		result, ierr := vm.Eval(fmt.Sprintf("input %d", replCount), strings.NewReader(text))
		if ierr != nil {
			color.Red("Eval error: %s\n\n", xin.FormatError(ierr))
			continue
		}

		vm.Frame.Put(
			fmt.Sprintf("_%d", replCount),
			result,
		)
		color.Yellow("%d ) %s\n\n", replCount, result.Repr())

		replCount++
	}
}
