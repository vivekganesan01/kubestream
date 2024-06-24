package k8

import (
	"context"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

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
	deployments, err := d.Client.AppsV1().Deployments(d.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing deployments: %v\n", err)
		return
	}
	fmt.Printf("Deployments in namespace %s for alias %s:\n", d.Namespace, d.GroupNameAlias)
	for _, deployment := range deployments.Items {
		fmt.Printf(" - name: %s | replicas: %d vs %d\n", deployment.Name, *deployment.Spec.Replicas, deployment.Status.AvailableReplicas)
	}
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
	// fmt.Printf("Deployments in namespace %s in alias %s:\n", namespace, groupNameAlias)
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"cluster", "type", "name", "desired", "current", "ready", "up-to-date"})
	for _, daemonset := range daemonsets.Items {
		desired := daemonset.Status.DesiredNumberScheduled
		current := daemonset.Status.CurrentNumberScheduled
		ready := daemonset.Status.NumberReady
		upToDate := daemonset.Status.UpdatedNumberScheduled
		t.AppendRow([]interface{}{ds.GroupNameAlias, "daemonset", daemonset.Name, desired, current, ready, upToDate})
		t.AppendSeparator()
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
	// fmt.Printf("Deployments in namespace %s in alias %s:\n", namespace, groupNameAlias)
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"cluster", "type", "name", "desired", "actual"})
	for _, st := range statefulsets.Items {
		t.AppendRow([]interface{}{ss.GroupNameAlias, "statefulset", st.Name, *st.Spec.Replicas, st.Status.AvailableReplicas})
		t.AppendSeparator()
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
