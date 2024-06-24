package k8

import (
	"context"
	"errors"
	"fmt"

	"kubestream/pkg/utilitycore"
	"net/http"
	"sync"

	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// var (
// 	errfailedToUpdateKubeconfig   = errors.New("failed to update kube config file")
// 	errfailedToInitiateRestClient = errors.New("failed to initiate k8 rest client")
// )

var defaultPath = "default"

// var filepath string = "kubeconfig.yaml"

func NewKubeClient(cn string, p string) (*kubernetes.Clientset, error) {
	fn := utilitycore.GetFn(NewKubeClient)

	var kubeconfig *string = &p
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("[%s] couldn't able to build the kube config file", fn))
		return nil, errors.New("couldn't able to build the kube config file")
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("[%s] couldn't able to create configuration", fn))
		return nil, errors.New("couldn't able to create configuration")
	}
	return clientset, nil
}

func CloseKubeClient(client *kubernetes.Clientset) bool {
	if client == nil {
		return false
	}
	rest, ok := client.CoreV1().RESTClient().(*restclient.RESTClient)
	if !ok || rest.Client == nil || rest.Client.Transport == nil {
		return false
	}
	if transport, ok := rest.Client.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
		// remove kubeconfig file if required
		return true
	}
	return false
}

func ListPods(client *kubernetes.Clientset, namespace string) {
	pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing pods: %v\n", err)
		return
	}
	fmt.Printf("Pods in namespace %s:\n", namespace)
	for _, pod := range pods.Items {
		fmt.Printf(" - %s\n", pod.Name)
	}
}

func ListNode(client *kubernetes.Clientset) {
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing nodes: %v\n", err)
	}
	for _, v := range nodes.Items {
		fmt.Printf(" - %s\n", v.Name)
	}
}

func GetResourceInformation(apiResourceType, namespace *string, group *string) {
	// todo: streamline the pre-requisite
	utilitycore.ParseConfigObject()
	ccContext := utilitycore.QueryConfigObject(*group)
	if ccContext == (&utilitycore.ObjectMetadata{}) {
		fmt.Println("group by current context is filtered as empty")
		return
	}
	switch *apiResourceType {
	// todo: come up with a better case switcher, goal is to support all the k8 objects
	case "deployments":
		var wg sync.WaitGroup
		for _, item := range ccContext.KubernetesCluster {
			wg.Add(1)
			go func(item utilitycore.KubeConfigMetadata, namespace string) {
				defer wg.Done()
				ac, err := NewKubeClient(item.NameAlias, item.Kubeconfig)
				if err != nil {
					log.Error().Msg("Error: failed to establish kubeclient with group")
					return
				}
				dMetaData := Deployments{
					Client:         ac,
					GroupNameAlias: item.NameAlias,
					Namespace:      namespace,
				}
				FetchStandardAPIResources(&dMetaData)
			}(item, *namespace)
		}

		wg.Wait()

	case "statefulsets":

		var wg sync.WaitGroup
		for _, item := range ccContext.KubernetesCluster {
			wg.Add(1)
			go func(item utilitycore.KubeConfigMetadata, namespace string) {
				defer wg.Done()
				ac, err := NewKubeClient(item.NameAlias, item.Kubeconfig)
				if err != nil {
					log.Error().Msg("Error: failed to establish kubeclient with group")
					return
				}
				ssMetadata := StatefulSets{
					Client:         ac,
					GroupNameAlias: item.NameAlias,
					Namespace:      namespace,
				}
				FetchStandardAPIResources(&ssMetadata)
			}(item, *namespace)
		}

		wg.Wait()

	case "daemonsets":

		var wg sync.WaitGroup
		for _, item := range ccContext.KubernetesCluster {
			wg.Add(1)
			go func(item utilitycore.KubeConfigMetadata, namespace string) {
				defer wg.Done()
				ac, err := NewKubeClient(item.NameAlias, item.Kubeconfig)
				if err != nil {
					log.Error().Msg("Error: failed to establish kubeclient with group")
					return
				}
				ds := Daemonsets{
					Client:         ac,
					GroupNameAlias: item.GroupBy,
					Namespace:      namespace,
				}
				FetchStandardAPIResources(&ds)
			}(item, *namespace)
		}

		wg.Wait()

	case "pods":
		ac, err := NewKubeClient(*group, defaultPath)
		if err != nil {
			log.Error().Msg("Error: failed to establish kubeclient with group")
			return
		}
		ListPods(ac, *namespace)
	default:
	}
}

// todo: move all the goroutine logic to
