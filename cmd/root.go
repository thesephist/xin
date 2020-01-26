package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "xin",
	Short: "Xin interpreter",
	Long:  "Xin is a purely functional, lisp-like concurrent general-purpose programming language.",
}

func Execute() error {
	return rootCmd.Execute()
}
