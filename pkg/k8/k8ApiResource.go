package k8

import (
	"context"
	"fmt"
	"kubestream/pkg/utilitycore"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

// todo: add filter native to kubernetes client api
type K8StandardApiResources interface {
	List()
	Update()
	Scale(rN string, replicas int32)
}

func FetchStandardAPIResources(apiResource K8StandardApiResources) {
	apiResource.List()
}

func PatchStandardAPIResources(apiResource K8StandardApiResources) {
	apiResource.Update()
}

type Deployments struct {
	Client         *kubernetes.Clientset
	Namespace      string
	GroupNameAlias string
	ResourceName   string
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

func (d *Deployments) Scale(rN string, replicas int32) {
	if d.Namespace == "all" {
		log.Warn().Msg("warn: user passed all namespace to it will search for deployment in all the namespace and rollout")
	}
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := d.Client.AppsV1().Deployments(d.Namespace).Get(context.TODO(), rN, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}
		result.Spec.Replicas = utilitycore.Int32Ptr(replicas) // reduce replica count
		_, updateErr := d.Client.AppsV1().Deployments(d.Namespace).Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Updated deployment...")
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

func (ds *Daemonsets) Scale(rN string, replicas int32) {
	fmt.Println("scaling all the deployment resources")
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

func (ss *StatefulSets) Scale(rN string, replicas int32) {
	fmt.Println("scaling all the statefulset resources")
}

type Secrets struct {
	Client         *kubernetes.Clientset
	Namespace      string
	GroupNameAlias string
}

func (se *Secrets) List() {
	secrets, err := se.Client.CoreV1().Secrets(se.Namespace).List(context.Background(), metav1.ListOptions{})

	if err != nil {
		fmt.Printf("Error listing secrets: %v\n", err)
		return
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"cluster", "type", "namespace", "name"})
	for _, ses := range secrets.Items {
		t.AppendRow([]interface{}{se.GroupNameAlias, "secrets", ses.Namespace, ses.Name})
	}
	t.Render()
}

func (se *Secrets) Update() {
	fmt.Println("update secret is not permitted ....")
}

func (se *Secrets) Scale(rN string, replicas int32) {
	fmt.Println("scale secret is not permitted ....")
}

type Configmaps struct {
	Client         *kubernetes.Clientset
	Namespace      string
	GroupNameAlias string
}

func (cm *Configmaps) List() {
	configmap, err := cm.Client.CoreV1().ConfigMaps(cm.Namespace).List(context.Background(), metav1.ListOptions{})

	if err != nil {
		fmt.Printf("Error listing secrets: %v\n", err)
		return
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"cluster", "type", "namespace", "name"})
	for _, cme := range configmap.Items {
		t.AppendRow([]interface{}{cm.GroupNameAlias, "configmap", cme.Namespace, cme.Name})
	}
	t.Render()
}

// todo: add logic to patch configmaps
func (cm *Configmaps) Update() {
	fmt.Println("update cm is not permitted ....")
}

func (cm *Configmaps) Scale() {
	fmt.Println("update cm is not permitted ....")
}
