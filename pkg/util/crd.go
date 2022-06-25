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
	Group:    "crd.k8s5g.com",
	Version:  "v1",
	Resource: "networkinfos",
}

type NetworkInfoSpec struct {
	Location int `json:"location"`
}

type NetworkInfo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec NetworkInfoSpec `json:"spec,omitempty"`
}

type NetworkInfoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []NetworkInfo `json:"items"`
}

func ListNetworkInfos(client dynamic.Interface, namespace string) (*NetworkInfoList, error) {
	list, err := client.Resource(gvr).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	data, err := list.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var ctList NetworkInfoList
	if err := json.Unmarshal(data, &ctList); err != nil {
		return nil, err
	}
	return &ctList, nil
}

func GetNetworkInfo(client dynamic.Interface, namespace string, name string) (*NetworkInfo, error) {
	utd, err := client.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	data, err := utd.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var ct NetworkInfo
	if err := json.Unmarshal(data, &ct); err != nil {
		return nil, err
	}
	return &ct, nil
}

func CreateNetworkInfoWithYaml(client dynamic.Interface, namespace string, yamlData string) (*NetworkInfo, error) {
	decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	if _, _, err := decoder.Decode([]byte(yamlData), nil, obj); err != nil {
		return nil, err
	}

	return createNetworkInfoWithUnstructured(client, namespace, obj)
}

func CreateNetworkInfo(client dynamic.Interface, namespace string, info *NetworkInfo) (*NetworkInfo, error) {
	decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	byteSlice, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	if _, _, err := decoder.Decode(byteSlice, nil, obj); err != nil {
		return nil, err
	}

	return createNetworkInfoWithUnstructured(client, namespace, obj)
}

func createNetworkInfoWithUnstructured(client dynamic.Interface, namespace string, obj *unstructured.Unstructured) (*NetworkInfo, error) {
	utd, err := client.Resource(gvr).Namespace(namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	data, err := utd.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var ct NetworkInfo
	if err := json.Unmarshal(data, &ct); err != nil {
		return nil, err
	}
	return &ct, nil
}

func UpdateNetworkInfo(client dynamic.Interface, namespace string, info *NetworkInfo) (*NetworkInfo, error) {
	decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	byteSlice, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	if _, _, err := decoder.Decode(byteSlice, nil, obj); err != nil {
		return nil, err
	}

	return updateNetworkInfoWithUnstructured(client, namespace, obj)
}

func UpdateNetworkInfoWithYaml(client dynamic.Interface, namespace string, yamlData string) (*NetworkInfo, error) {
	decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	if _, _, err := decoder.Decode([]byte(yamlData), nil, obj); err != nil {
		return nil, err
	}
	return updateNetworkInfoWithUnstructured(client, namespace, obj)
}

func updateNetworkInfoWithUnstructured(client dynamic.Interface, namespace string, obj *unstructured.Unstructured) (*NetworkInfo, error) {
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
	var ct NetworkInfo
	if err := json.Unmarshal(data, &ct); err != nil {
		return nil, err
	}
	return &ct, nil
}

func PatchNetworkInfo(client dynamic.Interface, namespace, name string, pt types.PatchType, data []byte) error {
	_, err := client.Resource(gvr).Namespace(namespace).Patch(context.TODO(), name, pt, data, metav1.PatchOptions{})
	return err
}

func DeleteNetworkInfo(client dynamic.Interface, namespace string, name string) error {
	return client.Resource(gvr).Namespace(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}
