package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/thesephist/xin/pkg/xin"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a Xin program file",
	Long:  "Execute a Xin program from a file specified by a path",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			run(args[0])
			return
		}

		fmt.Println("Program file unspecified")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func run(path string) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		color.Red("Error opening file: %s\n", err)
	}

	vm := xin.NewVm()
	_, ierr := vm.Eval(path, file)
	if ierr != nil {
		color.Red("Error: %s\n", xin.FormatError(ierr))
		return
	}
}
