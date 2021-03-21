package kubenventoree

import (
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

type Cmdoptions struct {
	AllTests       bool   `help:"Run all queries"`
	KubeConfigPath string `short:"k" default:"" help:"path to the kubeconfig file" type:"path"`
	Output         string `short:"o" required help:"output file name" type:"path"`
	OutputFormat   string `short:"f" default:yaml help:"format of the output (json, yaml)"`
}

type Kubenventoree struct {
	Options   *Cmdoptions
	ClientSet *kubernetes.Clientset
}

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
