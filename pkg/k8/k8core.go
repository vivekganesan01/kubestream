package k8

import (
	"context"
	"errors"
	"fmt"
	"kubestream/pkg/utilitycore"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	errfailedToUpdateKubeconfig   = errors.New("failed to update kube config file")
	errfailedToInitiateRestClient = errors.New("failed to initiate k8 rest client")
)

var defaultPath = "default"

func NewKubeClient(cn string, p string) (*kubernetes.Clientset, error) {
	fn := utilitycore.GetFn(NewKubeClient)
	// Define the context you want to switch to
	newContext := cn
	// Get the kubeconfig file path
	var kubeconfig string
	if p == "default" {
		kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	} else {
		kubeconfig = filepath.Join(p, ".kube", "config")
	}

	// Load the existing kubeconfig
	config, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		fmt.Printf("Failed to load kubeconfig: %v\n", err)
		return nil, clientcmd.NewEmptyConfigError("failed to laod the kube config file")
	}

	// Check if the context exists
	if _, exists := config.Contexts[newContext]; !exists {
		fmt.Printf("Context %s does not exist in kubeconfig\n", newContext)
		return nil, clientcmd.NewEmptyConfigError("given cluster context doesn't exists")
	}

	// Set the current context
	config.CurrentContext = newContext

	// Save the modified kubeconfig
	err = clientcmd.WriteToFile(*config, kubeconfig)
	if err != nil {
		fmt.Printf("Failed to write kubeconfig: %v\n", err)
		return nil, errfailedToUpdateKubeconfig
	}

	currentContextConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("[%s] couldn't able to build the kube config file", fn))
		return nil, errfailedToUpdateKubeconfig
	}

	clientset, err := kubernetes.NewForConfig(currentContextConfig)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("[%s] couldn't able to create configuration", fn))
		return nil, errfailedToInitiateRestClient
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

func ListDeployments(client *kubernetes.Clientset, groupAlias string, namespace string) {
	deployments, err := client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing deployments: %v\n", err)
		return
	}
	fmt.Printf("Deployments in namespace %s in alias %s:\n", namespace, groupAlias)
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"name", "desired", "actual"})
	for _, deployment := range deployments.Items {
		t.AppendRow([]interface{}{deployment.Name, *deployment.Spec.Replicas, deployment.Status.AvailableReplicas})
		t.AppendSeparator()
	}
	t.Render()
}

func ListStatefulset(client *kubernetes.Clientset, namespace string) {
	deployments, err := client.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing deployments: %v\n", err)
		return
	}
	fmt.Printf("Deployments in namespace %s:\n", namespace)
	for _, deployment := range deployments.Items {
		fmt.Printf(" - name: %s | replicas: %d vs %d\n", deployment.Name, *deployment.Spec.Replicas, deployment.Status.AvailableReplicas)
	}
}

func ListDaemonset(client *kubernetes.Clientset, namespace string) {
	daemonsets, err := client.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing deployments: %v\n", err)
		return
	}
	fmt.Printf("Deployments in namespace %s:\n", namespace)
	for _, daemonset := range daemonsets.Items {
		desired := daemonset.Status.DesiredNumberScheduled
		current := daemonset.Status.CurrentNumberScheduled
		ready := daemonset.Status.NumberReady
		upToDate := daemonset.Status.UpdatedNumberScheduled
		fmt.Printf(" - name: %s | Desired: %d, Current: %d, Ready: %d, Up-to-date: %d\n", daemonset.Name, desired, current, ready, upToDate)
	}
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
				ac, err := NewKubeClient(item.NameAlias, "default")
				if err != nil {
					log.Error().Msg("Error: failed to establish kubeclient with group")
					return
				}
				ListDeployments(ac, item.NameAlias, namespace)
			}(item, *namespace)
		}

		wg.Wait()

	case "statefulsets":
		ac, err := NewKubeClient(*group, defaultPath)
		if err != nil {
			log.Error().Msg("Error: failed to establish kubeclient with group")
			return
		}
		ListStatefulset(ac, *namespace)
	case "daemonsets":
		ac, err := NewKubeClient(*group, defaultPath)
		if err != nil {
			log.Error().Msg("Error: failed to establish kubeclient with group")
			return
		}
		ListDaemonset(ac, *namespace)
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
