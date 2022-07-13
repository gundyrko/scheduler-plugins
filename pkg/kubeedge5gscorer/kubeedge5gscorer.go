package kubeedge5gscorer

import (
	"context"
	"fmt"
	"math/rand"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"sigs.k8s.io/scheduler-plugins/pkg/util"
)

type KubeEdge5GScorer struct {
	handle framework.Handle
}

var _ framework.ScorePlugin = &KubeEdge5GScorer{}

const Name = "KubeEdge5GScorer"

func (k *KubeEdge5GScorer) Name() string {
	return Name
}

func (k *KubeEdge5GScorer) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {
	// TODO: use framework handle get CR and calculate score
	// klog.InfoS("====KubeEdge5GScorer starts -v2====")
	// klog.InfoS("====Try to get custom resource object====")

	// clientset := k.handle.ClientSet()
	// client := clientset.AppsV1().RESTClient()

	// kubeconfig := "/kube/config"
	// config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		klog.Error(err)
		return 0, framework.AsStatus(fmt.Errorf("generating k8s configs for node %q: %w", nodeName, err))
	}

	// 	createData := `
	// apiVersion: "crd.k8s5g.com/v1"
	// kind: NetworkInfo
	// metadata:
	//   name: test-info
	// spec:
	//   location: 43
	// `
	// 	ct, err := util.CreateNetworkInfoWithYaml(client, "scheduler-plugins", createData)
	namespace := "scheduler-plugins"
	info := &util.NetworkInfo{}
	info.APIVersion = "crd.k8s5g.com/v1"
	info.Kind = "NetworkInfo"
	info.Name = "test-info-" + nodeName
	info.Spec.Location = rand.Int() % 100
	_, err = util.CreateNetworkInfo(client, namespace, info)
	if err != nil {
		_, err = util.UpdateNetworkInfo(client, namespace, info)
		if err != nil {
			klog.Error(err, " when updating for node:", nodeName, ", pod:", p.Name)
			// return 0, framework.AsStatus(fmt.Errorf("creating custom resource for node %q: %w", nodeName, err))
		}
	}
	// if err == nil {
	// 	klog.Info("created/updated for node:", nodeName, ", pod:", p.Name, " in namespace: ", ct.Namespace, "crd name: ", ct.Name, ", Val: ", ct.Spec.Location)
	// }
	// list, err := util.ListNetworkInfos(client, namespace)
	// if err != nil {
	// 	klog.Error(err)
	// 	return 0, framework.AsStatus(fmt.Errorf("listing custom resources for node %q: %w", nodeName, err))
	// }
	// for _, t := range list.Items {
	// 	klog.Info("getting location for all nodes:", nodeName, ", pod:", p.Name, " in namespace: ", t.Namespace, "crd name: ", t.Name, ", Val: ", t.Spec.Location)
	// }

	newInfo, err := util.GetNetworkInfo(client, namespace, info.Name)
	if err != nil {
		klog.Error(err)
		return 0, framework.AsStatus(fmt.Errorf("getting custom resources for node %q: %w", nodeName, err))
	}
	// klog.InfoS("====Get success====")
	// klog.Info(data)
	score := newInfo.Spec.Location
	klog.Infof("Scheduling pod %q: node %q has score %d", p.Name, nodeName, score)
	return int64(score), nil
}

// ScoreExtensions of the Score plugin.
func (k *KubeEdge5GScorer) ScoreExtensions() framework.ScoreExtensions {
	return nil
}

// New initializes a new plugin and returns it.
func New(_ runtime.Object, h framework.Handle) (framework.Plugin, error) {
	return &KubeEdge5GScorer{handle: h}, nil
}
