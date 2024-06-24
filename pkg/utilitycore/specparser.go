package utilitycore

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const filePath = "/config/kubeobject.yaml"

type KubeConfigMetadata struct {
	NameAlias  string `yaml:"name_alias"`
	Kubeconfig string `yaml:"kubeconfig"`
	GroupBy    string `yaml:"group_by"`
}
type ObjectMetadata struct {
	KubernetesCluster []KubeConfigMetadata `yaml:"kubernetes_cluster"`
}

// type specParser interface {
// 	parseConfigObject(filePath string) *objectMetadata
// 	queryConfigObject(groupBy string) *objectMetadata
// }

var globalConfig ObjectMetadata

func ParseConfigObject() *ObjectMetadata {
	// Read the YAML file
	cwd, _ := os.Getwd()
	fp := cwd + filePath
	file, err := os.Open(fp)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Unmarshal the YAML into the Config struct
	err = yaml.Unmarshal(data, &globalConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// todo: make readonly object
	return &globalConfig
}

func QueryConfigObject(groupBy string) *ObjectMetadata {
	if groupBy == "all" {
		return &globalConfig
	}
	var currentContext ObjectMetadata
	for _, items := range globalConfig.KubernetesCluster {
		if items.GroupBy == groupBy {
			currentContext.KubernetesCluster = append(currentContext.KubernetesCluster, items)
		}
	}
	if len(currentContext.KubernetesCluster) == 0 {
		return &ObjectMetadata{}
	}
	return &currentContext
}
