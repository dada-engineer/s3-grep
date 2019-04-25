package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/dabdada/s3-grep/config"
	"github.com/dabdada/s3-grep/s3"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "s3-grep [search query]",
	Short: "Grep contents of an object in S3",
	Long: "Grep contents of an object in S3",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("s3-grep requires a search query argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			PrintVersion()
			return
		}
		if ok := s3.IsBucket(); !ok {
			fmt.Printf("The bucket name `%s` was not found in profile `%s`\n", bucketName, profile)
			return
		} else {
			// run the actual grep command
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
var profile string
var bucketName string
var path string

func init() {
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "Print the version of s3-grep")
	rootCmd.Flags().StringVarP(&profile, "profile", "", "", "The AWS profile the S3 bucketName is hosted in")
	rootCmd.Flags().StringVarP(&bucketName, "bucket", "b", "", "The bucketName name to grep in")
	rootCmd.Flags().StringVarP(&bucketName, "path", "p", "/", "The path to grep in")

	rootCmd.MarkFlagRequired("profile")
	rootCmd.MarkFlagRequired("bucket")
}
