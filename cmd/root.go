package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const version = "0.1"

var rootCmd = &cobra.Command{
	Use:   "xin [files]",
	Short: "Run Xin programs",
	Long:  "Xin is an extensible, functional, lisp-like general-purpose programming language.",
	Example: strings.Join([]string{
		"  xin\t\t\tstart a repl",
		"  xin prog.xin\t\trun prog.xin",
		"  echo file | xin\trun from stdin",
	}, "\n"),
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		// check if running files
		if len(args) >= 1 {
			run(args[0])
			return
		}

		// check if there's stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			stdin()
			return
		}

		// if all else fails, start repl
		repl()
	},
}

func Execute() error {
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "v%s" .Version}}
`)

	return rootCmd.Execute()
}
