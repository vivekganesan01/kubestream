/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"kubestream/pkg/k8"

	"github.com/spf13/cobra"
)

var (
	apiResource string
	namespace   string
	group       string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Fetch API resource information",
	Long:  `Command to fetch information and status about the live kubernetes resource object.`,
	Run: func(cmd *cobra.Command, args []string) {
		k8.GetResourceInformation(&apiResource, &namespace, &group)
	},
}

func init() {
	getCmd.Flags().StringVarP(&apiResource, "apiresource", "a", "all", "name api resource to be fetched")
	getCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "namespace")
	getCmd.Flags().StringVarP(&group, "groupby", "g", "", "group by")
}
