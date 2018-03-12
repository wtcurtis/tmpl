package cmd

import (
	"github.com/spf13/cobra"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var region string
var awsSession *session.Session

var RootCommand = &cobra.Command{
	Use:   "tmpl",
	Long:  "Entry point for the app.  Used to run all app processes.",
}

func init() {
	RootCommand.PersistentFlags().StringVarP(&region, "region", "r", "", "AWS region to pull from")
	cobra.OnInitialize(configure)

	awsSession = session.Must(session.NewSession(&aws.Config{
		Region: &region,
	}))
}

func configure() {
}
