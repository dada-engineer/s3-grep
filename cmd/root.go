package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/dabdada/s3-grep/cli"
	"github.com/dabdada/s3-grep/config"
	"github.com/dabdada/s3-grep/s3"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   `s3-grep search query --bucket --profile [--version] [-i] [--help]`,
	Short: "Grep contents of an object in S3",
	Long:  "Grep contents of an object in S3",
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

		if help {
			cmd.Usage()
			return
		}

		session, err := config.NewAWSSession(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if ok := s3.IsBucket(*session, bucketName); !ok {
			fmt.Printf("The bucket name `%s` was not found in profile `%s`\n", bucketName, profile)
			return
		} else {
			cli.Grep(session, bucketName, args[0], ignoreCase)
			return
		}
	},
}

// Execute the root command s3-grep
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	help       bool
	version    bool
	profile    string
	bucketName string
	ignoreCase bool
)

func init() {
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "Print the version of s3-grep")
	rootCmd.Flags().BoolVarP(&version, "help", "h", false, "Print the usage of s3-grep")
	rootCmd.Flags().StringVarP(&profile, "profile", "p", "", "The AWS profile the S3 bucketName is hosted in")
	rootCmd.Flags().StringVarP(&bucketName, "bucket", "b", "", "The bucketName name to grep in")
	rootCmd.Flags().BoolVarP(&ignoreCase, "", "i", false, "Ignore case of the search query while grepping")

	rootCmd.MarkFlagRequired("profile")
	rootCmd.MarkFlagRequired("bucket")
}
