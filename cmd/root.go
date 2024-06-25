/*
Copyright © 2024 vivekganesan01@gmail.com
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubestream",
	Short: "kubestream is an utility to stream k8 objects from multiple cluster",
	Long: `kubestream is a cli library to stream information and metadata of kubernetes object
from difference kubernetes cluster as a single entity simultaneously. For example:

./kubestream --help

kubestream get --api_resource=deployment --namespace="all" --group_by="aws-us-east"
kubestream get --api_resource=daemonset --namespace="default" --group_by="local"
kubestream get --api_resource=statefulset --namespace="default" --group_by="ibm-us-south"
kubestream get --api_resource=statefulset --namespace="kube-system" --group_by="all"
.`,

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(getCmd)
}
