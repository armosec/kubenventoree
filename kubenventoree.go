package kubenventoree

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type cmdoptions struct {
	AllTests       bool   `help:"Run all queries"`
	KubeConfigPath string `short:"k" default:"" help:"path to the kubeconfig file" type:"path"`
	Output         string `short:"o" required help:"output file name" type:"path"`
	OutputFormat   string `short:"f" default:text help:"format of the output (json, yaml, text)"`
}

var cliOptions cmdoptions

type Kubenventoree struct {
	Options   *cmdoptions
	ClientSet *kubernetes.Clientset
}

type DiscoveryResult struct {
	ImageInventory   *ClusterImageInventory `json:"image_inventory"`
	WokloadInventory *WorkloadInventory     `json:"workload_inventory"`
	ClusterInfo      *ClusterInfo           `json:"cluster_info"`
}

func main() {
	ctx := kong.Parse(&cliOptions)

	//fmt.Printf("Command: %s\n", ctx.Command())
	//fmt.Printf("Struct: %s %s", cliOptions.Output, cliOptions.OutputFormat)
	cs, err := GetK8sClientset(cliOptions.KubeConfigPath)
	if err != nil {
		log.Fatalf("Cannot access kubernetes config file (%s)", err)
		ctx.Exit(1)
	}

	_, err = cs.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Cannot access kubernetes API (%s)", err)
		ctx.Exit(1)
	}

	result := DiscoveryResult{}
	k := Kubenventoree{Options: &cliOptions, ClientSet: cs}

	result.ImageInventory, err = k.ReadImageInventory()
	if err != nil {
		log.Fatalf("Cannot read image inventory (%s)", err)
		ctx.Exit(1)
	}

	result.WokloadInventory, err = k.ReadWorkloadInventory()
	if err != nil {
		log.Fatalf("Cannot read workload inventory (%s)", err)
		ctx.Exit(1)
	}

	result.ClusterInfo, err = k.ReadClusterInfo()
	if err != nil {
		log.Fatalf("Cannot read cluster information (%s)", err)
		ctx.Exit(1)
	}

	var result_text string
	switch cliOptions.OutputFormat {
	case "json":
		{
			result_bytes, err := json.Marshal(result)
			if err != nil {
				log.Fatalf("json conversion failed %s", err)
			}
			result_text = string(result_bytes)
		}
	case "yaml":
		{
			result_bytes, err := yaml.Marshal(result)
			if err != nil {
				log.Fatalf("yaml conversion failed %s", err)
			}
			result_text = string(result_bytes)
		}
	}

	if cliOptions.Output == "-" {
		fmt.Println(result_text)
	} else {
		f, err := os.Create(cliOptions.Output)
		if err != nil {
			log.Fatalf("cannot open output file (%s)", err)
		}
		defer f.Close()
		f.WriteString(result_text)
	}

	os.Exit(0)
}
