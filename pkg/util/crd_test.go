package util

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func getClient() dynamic.Interface {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return client
}

func TestCreateNetworkInfoWithYaml(t *testing.T) {
	client := getClient()
	createData := `
apiVersion: "crd.k8s5g.com/v1"
kind: NetworkInfo
metadata:
  name: test-info
spec:
  location: 43
`
	ct, err := CreateNetworkInfoWithYaml(client, "scheduler-plugins", createData)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s %s %d\n", ct.Namespace, ct.Name, ct.Spec.Location)
}

func TestListNetworkInfos(t *testing.T) {
	client := getClient()
	list, err := ListNetworkInfos(client, "scheduler-plugins")
	if err != nil {
		panic(err)
	}
	for _, t := range list.Items {
		fmt.Printf("%s %s %d\n", t.Namespace, t.Name, t.Spec.Location)
	}
}

func TestGetNetworkInfo(t *testing.T) {
	client := getClient()
	ct, err := GetNetworkInfo(client, "scheduler-plugins", "test-info")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s %s %d\n", ct.Namespace, ct.Name, ct.Spec.Location)
}

func TestUpdateNetworkInfoWithYaml(t *testing.T) {
	client := getClient()
	updateData := `
apiVersion: "stable.example.com/v1"
kind: NetworkInfo
metadata:
  name: test-info
spec:
  location: 47
`
	ct, err := UpdateNetworkInfoWithYaml(client, "scheduler-plugins", updateData)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s %s %d\n", ct.Namespace, ct.Name, ct.Spec.Location)

}

func TestDeleteNetworkInfo(t *testing.T) {
	client := getClient()
	if err := DeleteNetworkInfo(client, "scheduler-plugins", "test-info"); err != nil {
		panic(err)
	}
}
