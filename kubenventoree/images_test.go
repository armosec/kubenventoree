package kubenventoree

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestImageInventory(t *testing.T) {
	cs, err := GetK8sClientset("")
	if err != nil {
		t.Error()
		return
	}
	k := Kubenventoree{ClientSet: cs}
	imageList, err := k.ReadImageInventory()
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

	//t.Error() // to indicate test failed
}
