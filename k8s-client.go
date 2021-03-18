package kubenventoree

import (
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func GetK8sClientset(kubeconfig_path string) (*kubernetes.Clientset, error) {
	var kubeconfig string
	if kubeconfig_path == "" {
		home := homedir.HomeDir()
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		kubeconfig = kubeconfig_path
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}
