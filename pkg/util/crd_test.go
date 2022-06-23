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

func TestCreateCrontabWithYaml(t *testing.T) {
	client := getClient()
	createData := `
apiVersion: "stable.example.com/v1"
kind: CronTab
metadata:
  name: cron-4
spec:
  cronSpec: "* * * * */15"
  image: my-awesome-cron-image-4
`
	ct, err := CreateCrontabWithYaml(client, "default", createData)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s %s %s %s\n", ct.Namespace, ct.Name, ct.Spec.CronSpec, ct.Spec.Image)
}

func TestListCrontabs(t *testing.T) {
	client := getClient()
	list, err := ListCrontabs(client, "default")
	if err != nil {
		panic(err)
	}
	for _, t := range list.Items {
		fmt.Printf("%s %s %s %s\n", t.Namespace, t.Name, t.Spec.CronSpec, t.Spec.Image)
	}
}

func TestGetCrontab(t *testing.T) {
	client := getClient()
	ct, err := GetCrontab(client, "default", "cron-4")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s %s %s %s\n", ct.Namespace, ct.Name, ct.Spec.CronSpec, ct.Spec.Image)
}

func TestUpdateCrontabWithYaml(t *testing.T) {
	client := getClient()
	updateData := `
apiVersion: "stable.example.com/v1"
kind: CronTab
metadata:
  name: cron-4
spec:
  cronSpec: "* * * * */8"
  image: my-awesome-cron-image-4-update
`
	ct, err := UpdateCrontabWithYaml(client, "default", updateData)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s %s %s %s\n", ct.Namespace, ct.Name, ct.Spec.CronSpec, ct.Spec.Image)

}

func TestDeleteCrontab(t *testing.T) {
	client := getClient()
	if err := DeleteCrontab(client, "default", "cron-4"); err != nil {
		panic(err)
	}
}
