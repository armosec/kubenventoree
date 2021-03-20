package kubenventoree

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestWorkloadInventory(t *testing.T) {
	cs, err := GetK8sClientset("")
	if err != nil {
		t.Error()
		return
	}
	k := Kubenventoree{ClientSet: cs}
	imageList, err := k.ReadWorkloadInventory()
	if err != nil {
		t.Error()
		return
	}
	b, err := json.Marshal(imageList)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

func TestClusterInfo(t *testing.T) {
	cs, err := GetK8sClientset("")
	if err != nil {
		t.Error()
		return
	}
	k := Kubenventoree{ClientSet: cs}
	imageList, err := k.ReadClusterInfo()
	if err != nil {
		t.Error()
		return
	}
	b, err := json.Marshal(imageList)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}
