package kubenventoree

import (
	"context"
	"log"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type WorkloadInventory struct {
	PodCount                   int `json:"pods"`
	DeploymentCount            int `json:"deployments"`
	StateFullSetsCount         int `json:"statefullsets"`
	DaemonSetsCount            int `json:"daemonsets"`
	ReplicaSetsCount           int `json:"replicasets"`
	ControlerRevisionCount     int `json:"controlerrevision"`
	ReplicationControllerCount int `json:"replicationcontroller"`
}

type ClusterInfo struct {
	ClusterProvider  string `json:"provider,omitempty"`
	Cloud            string `json:"cloud,omitempty"`
	HasOperators     bool   `json:"has_operators"`
	HasHelm          bool   `json:"has_helm"`
	SaWithPullSecret bool   `json:"serviceaccount_has_pullsecret"`
	UsingIngress     bool   `json:"has_ingress"`
}

func (k *Kubenventoree) ReadWorkloadInventory() (*WorkloadInventory, error) {
	var wlInventory WorkloadInventory

	pods, err := k.ClientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Cannot read pods\n")
		return nil, err
	}
	wlInventory.PodCount = len(pods.Items)

	replicationcontrollers, err := k.ClientSet.CoreV1().ReplicationControllers("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Cannot read replicationcontrollers\n")
		return nil, err
	}
	wlInventory.ReplicationControllerCount = len(replicationcontrollers.Items)

	deployments, err := k.ClientSet.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Cannot read deployments\n")
		return nil, err
	}
	wlInventory.DeploymentCount = len(deployments.Items)

	daemonsets, err := k.ClientSet.AppsV1().DaemonSets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Cannot read daemonsets\n")
		return nil, err
	}
	wlInventory.DaemonSetsCount = len(daemonsets.Items)

	statefulsets, err := k.ClientSet.AppsV1().StatefulSets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Cannot read statefulsets\n")
		return nil, err
	}
	wlInventory.StateFullSetsCount = len(statefulsets.Items)

	replicasets, err := k.ClientSet.AppsV1().ReplicaSets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Cannot read replicasets\n")
		return nil, err
	}
	wlInventory.ReplicaSetsCount = len(replicasets.Items)

	controlerrevisions, err := k.ClientSet.AppsV1().ControllerRevisions("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Cannot read controlerrevisions\n")
		return nil, err
	}
	wlInventory.ControlerRevisionCount = len(controlerrevisions.Items)

	return &wlInventory, nil
}

func (k *Kubenventoree) ReadClusterInfo() (*ClusterInfo, error) {
	var clusterInfo ClusterInfo

	// operator
	pods, err := k.ClientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Cannot read pods\n")
		return nil, err
	}

	clusterInfo.HasHelm = false
	clusterInfo.HasOperators = false
	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, "operator") {
			clusterInfo.HasOperators = true
		}
		if strings.Contains(pod.Name, "tiller") {
			clusterInfo.HasHelm = true
		}
	}

	secrets, err := k.ClientSet.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Cannot read secrets\n")
		return nil, err
	}
	for _, secret := range secrets.Items {
		if strings.Contains(secret.Name, ".helm.") {
			clusterInfo.HasHelm = true
			break
		}
	}

	// service account imagepull /  helm
	sas, err := k.ClientSet.CoreV1().ServiceAccounts("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Cannot read sas\n")
		return nil, err
	}

	clusterInfo.SaWithPullSecret = false
	for _, sa := range sas.Items {
		if 0 < len(sa.ImagePullSecrets) {
			clusterInfo.SaWithPullSecret = true
			break
		}
	}

	cRoleBindings, err := k.ClientSet.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Cannot read clusterrolebindings\n")
		return nil, err
	}
	clusterInfo.ClusterProvider = "Unknown"
	for _, crb := range cRoleBindings.Items {
		if strings.HasPrefix(crb.Name, "eks") {
			clusterInfo.ClusterProvider = "EKS"
			break
		}
		if strings.Contains(crb.Name, "gke") {
			clusterInfo.ClusterProvider = "GKE"
			break
		}
	}

	nodes, err := k.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if 0 < len(nodes.Items) {
		node := nodes.Items[0]
		if strings.Contains(node.Name, "compute.internal") {
			clusterInfo.Cloud = "AWS"
		} else if strings.HasPrefix(node.Name, "gke") {
			clusterInfo.Cloud = "GCP"
		} else if strings.HasPrefix(node.Name, "aks") {
			clusterInfo.Cloud = "Azure"
			clusterInfo.ClusterProvider = "AKS"
		} else {
			clusterInfo.Cloud = "Unknown"
		}
	}

	clusterInfo.UsingIngress = false
	ingresses, err := k.ClientSet.NetworkingV1().Ingresses("").List(context.TODO(), metav1.ListOptions{})
	if err == nil && 0 < len(ingresses.Items) {
		clusterInfo.UsingIngress = true
	}

	return &clusterInfo, nil
}
