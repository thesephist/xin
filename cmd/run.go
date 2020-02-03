package cmd

import (
	"fmt"

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
