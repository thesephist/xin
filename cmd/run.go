package cmd

import (
	"fmt"
	"os"

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

func run(filePath string) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		fmt.Println("Error opening file:", err)
	}

	vm := xin.NewVm()
	_, err = vm.Eval(file)
	if err != nil {
		fmt.Println("Eval error:", err.Error())
		return
	}
}
