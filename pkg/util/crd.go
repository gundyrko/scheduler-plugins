package util

import (
	"context"
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

var gvr = schema.GroupVersionResource{
	Group:    "stable.example.com",
	Version:  "v1",
	Resource: "crontabs",
}

type CrontabSpec struct {
	CronSpec string `json:"cronSpec"`
	Image    string `json:"image"`
}

type Crontab struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CrontabSpec `json:"spec,omitempty"`
}

type CrontabList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Crontab `json:"items"`
}

func ListCrontabs(client dynamic.Interface, namespace string) (*CrontabList, error) {
	list, err := client.Resource(gvr).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	data, err := list.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var ctList CrontabList
	if err := json.Unmarshal(data, &ctList); err != nil {
		return nil, err
	}
	return &ctList, nil
}

func GetCrontab(client dynamic.Interface, namespace string, name string) (*Crontab, error) {
	utd, err := client.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	data, err := utd.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var ct Crontab
	if err := json.Unmarshal(data, &ct); err != nil {
		return nil, err
	}
	return &ct, nil
}

func CreateCrontabWithYaml(client dynamic.Interface, namespace string, yamlData string) (*Crontab, error) {
	decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	if _, _, err := decoder.Decode([]byte(yamlData), nil, obj); err != nil {
		return nil, err
	}

	utd, err := client.Resource(gvr).Namespace(namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	data, err := utd.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var ct Crontab
	if err := json.Unmarshal(data, &ct); err != nil {
		return nil, err
	}
	return &ct, nil
}

func UpdateCrontabWithYaml(client dynamic.Interface, namespace string, yamlData string) (*Crontab, error) {
	decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	if _, _, err := decoder.Decode([]byte(yamlData), nil, obj); err != nil {
		return nil, err
	}
	name := obj.GetName()
	utd, err := client.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	obj.SetResourceVersion(utd.GetResourceVersion())
	utd, err = client.Resource(gvr).Namespace(namespace).Update(context.TODO(), obj, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	data, err := utd.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var ct Crontab
	if err := json.Unmarshal(data, &ct); err != nil {
		return nil, err
	}
	return &ct, nil
}

func PatchCrontab(client dynamic.Interface, namespace, name string, pt types.PatchType, data []byte) error {
	_, err := client.Resource(gvr).Namespace(namespace).Patch(context.TODO(), name, pt, data, metav1.PatchOptions{})
	return err
}

func DeleteCrontab(client dynamic.Interface, namespace string, name string) error {
	return client.Resource(gvr).Namespace(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}
