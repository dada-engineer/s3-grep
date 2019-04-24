package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "s3-grep",
	Short: "Grep contents of an object in S3",
	Long: "Grep contents of an object in S3",
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			PrintVersion()
			return
		}
		cmd.Usage()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var version bool

func init() {
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "print the version of s3-grep")
}
