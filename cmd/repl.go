package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
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
		}

		vm.Frame.Put(
			fmt.Sprintf("_%d", replCount),
			result,
		)
		color.Yellow("%d ) %s\n\n", replCount, result)

		replCount++
	}
}
