package k8

import (
	"context"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// todo: add filter native to kubernetes client api
type K8StandardApiResources interface {
	List()
	Update()
}

type Deployments struct {
	Client         *kubernetes.Clientset
	Namespace      string
	GroupNameAlias string
}

func (d *Deployments) List() {
	if d.Namespace == "all" {
		d.Namespace = ""
	}
	deployments, err := d.Client.AppsV1().Deployments(d.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing deployments: %v\n", err)
		return
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"cluster", "type", "namespace", "name", "desired", "current"})
	for _, deployment := range deployments.Items {
		t.AppendRow([]interface{}{d.GroupNameAlias, "deployment", deployment.Namespace, deployment.Name, *deployment.Spec.Replicas, deployment.Status.AvailableReplicas})
	}
	t.Render()
}

func (d *Deployments) Update() {
	fmt.Println("updating all the deployment resources")
}

type Daemonsets struct {
	Client         *kubernetes.Clientset
	Namespace      string
	GroupNameAlias string
}

func (ds *Daemonsets) List() {
	daemonsets, err := ds.Client.AppsV1().DaemonSets(ds.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing deployments: %v\n", err)
		return
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"cluster", "type", "namespace", "name", "desired", "current", "ready", "up-to-date"})
	for _, daemonset := range daemonsets.Items {
		desired := daemonset.Status.DesiredNumberScheduled
		current := daemonset.Status.CurrentNumberScheduled
		ready := daemonset.Status.NumberReady
		upToDate := daemonset.Status.UpdatedNumberScheduled
		// todo: load ds.Namespace from response
		t.AppendRow([]interface{}{ds.GroupNameAlias, "daemonset", daemonset.Namespace, daemonset.Name, desired, current, ready, upToDate})
	}
	t.Render()
}
func (ds *Daemonsets) Update() {
	fmt.Println("updating all the daemonsets resources")
}

type StatefulSets struct {
	Client         *kubernetes.Clientset
	Namespace      string
	GroupNameAlias string
}

func (ss *StatefulSets) List() {
	statefulsets, err := ss.Client.AppsV1().StatefulSets(ss.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing deployments: %v\n", err)
		return
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"cluster", "type", "namespace", "name", "desired", "actual"})
	for _, st := range statefulsets.Items {
		t.AppendRow([]interface{}{ss.GroupNameAlias, "statefulset", st.Namespace, st.Name, *st.Spec.Replicas, st.Status.AvailableReplicas})
	}
	t.Render()
}

func (ss *StatefulSets) Update() {
	fmt.Println("updating all the statefulset resources")
}

func FetchStandardAPIResources(apiResource K8StandardApiResources) {
	apiResource.List()
}

func PatchStandardAPIResources(apiResource K8StandardApiResources) {
	apiResource.Update()
}
