package kubenventoree

import (
	"context"
	"fmt"
	"strings"

	"github.com/distribution/distribution/reference"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ImageDescription struct {
	ImageRepository string `json:"image_repository,omitempty"`
	ImageTag        string `json:"image_tag,omitempty"`
	ImagePullable   string `json:"image_pullable,omitempty"`
	Count           int    `json:"count,omitempty"`
	HasSecret       bool   `json:"has_secret"`
}

type ClusterInventoryMetaInfo struct {
	UsingPrivateRepos bool `json:"using_private_image_repositories"`
	UsingEcr          bool `json:"using_ecr"`
	UsingGcr          bool `json:"using_gcr"`
}

type ClusterImageInventory struct {
	ImageList []ImageDescription       `json:"images"`
	MetaInfo  ClusterInventoryMetaInfo `json:"extra_info"`
}

func (k *Kubenventoree) readAllImages() ([]ImageDescription, error) {
	pods, err := k.ClientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	clusterImageMap := make(map[string]ImageDescription)

	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			if imageDescription, found := clusterImageMap[container.Image]; found {
				imageDescription.Count = imageDescription.Count + 1
			} else {
				imageNamed, err := reference.ParseNamed(container.Image)
				if err != nil {
					continue
				}
				secretFlag := false
				if len(pod.Spec.ImagePullSecrets) > 0 {
					secretFlag = true
				}
				imageDescription = ImageDescription{
					ImageRepository: reference.Domain(imageNamed),
					ImageTag:        imageNamed.Name(),
					Count:           1,
					HasSecret:       secretFlag,
				}
				clusterImageMap[container.Image] = imageDescription
			}
		}
		for _, container := range pod.Spec.InitContainers {
			if imageDescription, found := clusterImageMap[container.Image]; found {
				imageDescription.Count = imageDescription.Count + 1
			} else {
				imageNamed, err := reference.ParseNamed(container.Image)
				if err != nil {
					continue
				}
				secretFlag := false
				if len(pod.Spec.ImagePullSecrets) > 0 {
					secretFlag = true
				}
				imageDescription = ImageDescription{
					ImageRepository: reference.Domain(imageNamed),
					ImageTag:        imageNamed.Name(),
					Count:           1,
					HasSecret:       secretFlag,
				}
				clusterImageMap[container.Image] = imageDescription
			}
		}
	}

	imageList := []ImageDescription{}
	for _, value := range clusterImageMap {
		imageList = append(imageList, value)
	}

	return imageList, nil
}

func isUsingEcr(imageList []ImageDescription) bool {
	for _, imageDesc := range imageList {
		if strings.HasPrefix(imageDesc.ImageRepository, "ecr") {
			return true
		}
	}

	return false
}

func isUsingGcr(imageList []ImageDescription) bool {
	for _, imageDesc := range imageList {
		if strings.HasSuffix(imageDesc.ImageRepository, "gcr.io") {
			return true
		}
	}
	return false
}

func isUsingPrivateRepo(imageList []ImageDescription) bool {
	for _, imageDesc := range imageList {
		if imageDesc.HasSecret {
			return true
		}
	}
	return false
}

func (k *Kubenventoree) ReadImageInventory() (*ClusterImageInventory, error) {
	imageList, err := k.readAllImages()
	if err != nil {
		return nil, err
	}

	fmt.Printf("ReadImageInventory runing\n")

	var extraInfo ClusterInventoryMetaInfo
	extraInfo.UsingEcr = isUsingEcr(imageList)
	extraInfo.UsingGcr = isUsingGcr(imageList)
	extraInfo.UsingPrivateRepos = isUsingPrivateRepo(imageList)

	return &ClusterImageInventory{ImageList: imageList, MetaInfo: extraInfo}, nil
}
