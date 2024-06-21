package main

import (
	"fmt"
	"kubestream/pkg/utilitycore"
)

func main() {
	var ccContext *utilitycore.ObjectMetadata
	utilitycore.ParseConfigObject()
	ccContext = utilitycore.QueryConfigObject("optest")
	if ccContext == (&utilitycore.ObjectMetadata{}) {
		fmt.Println("Current context is empty")
	}
	for _, v := range ccContext.KubernetesCluster {
		fmt.Println("Chosen:", v.NameAlias)
	}
}
