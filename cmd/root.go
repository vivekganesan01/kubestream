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
from difference kubernetes cluster as a single entity. For example:

kubestream get --api_resource=deployment --namespace=default --group_by=all
kubestream get --api_resource=deployment --namespace=default --group_by=${CLUSTER_NAME}
kubestream get --api_resource=statefulset --namespace=default --group_by=${REGULAR_EXPRESSION}
kubestream get --api_resource=pod --namespace=default --group_by=${REGULAR_EXPRESSION} [--condition='deployment,statefulset' --filter='crashloop'] 
.`,

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(getCmd)
}
