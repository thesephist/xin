package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thesephist/xin/pkg/xin"
)

var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Run a Xin repl",
	Long:  "Start an interactive Xin language interpreter session",
	Run: func(cmd *cobra.Command, args []string) {
		repl()
	},
}

func init() {
	rootCmd.AddCommand(replCmd)
}

func repl() {
	vm := xin.NewVm()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		text, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Repl error:", err.Error())
		}

		result, err := vm.Eval(strings.NewReader(text))
		if err != nil {
			fmt.Println("Eval error:", err.Error())
		} else {
			fmt.Println(result)
		}
	}
}
