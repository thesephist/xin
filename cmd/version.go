package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of the Xin CLI",
	Long:  "Print the version number of the Xin CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("xin v%s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
