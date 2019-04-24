package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version of s3-grep
var Version = "dev"

// Print the version of s3-grep
func PrintVersion() {
	fmt.Printf("%s/n", Version)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of s3-grep",
	Long:  "Print the version of s3-grep",
	Run: func(cmd *cobra.Command, args []string) {
		PrintVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
