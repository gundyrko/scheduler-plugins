package kubeedge5gscorer

import (
	"context"
	"fmt"
	"math/rand"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
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
	klog.InfoS("====KubeEdge5GScorer starts -v2====")
	klog.InfoS("====Try to get custom resource object====")

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
		return 0, framework.AsStatus(fmt.Errorf("2. getting custom resource for node %q: %w", nodeName, err))
	}

	clientset, err := kubernetes.NewForConfig(config)

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	// 	createData := `
	// apiVersion: "stable.example.com/v1"
	// kind: CronTab
	// metadata:
	//   name: cron-4
	// spec:
	//   cronSpec: "* * * * */15"
	//   image: my-awesome-cron-image-4
	// `
	// ct, err := util.CreateCrontabWithYaml(client, "default", createData)
	// if err != nil {
	// 	klog.Error(err)
	// 	return 0, framework.AsStatus(fmt.Errorf("3. getting custom resource for node %q: %w", nodeName, err))
	// }
	// fmt.Printf("%s %s %s %s\n", ct.Namespace, ct.Name, ct.Spec.CronSpec, ct.Spec.Image)

	list, err := util.ListCrontabs(client, "default")
	if err != nil {
		klog.Error(err)
		return 0, framework.AsStatus(fmt.Errorf("4. getting custom resource for node %q: %w", nodeName, err))
	}
	for _, t := range list.Items {
		klog.Info("%s %s %s %s\n", t.Namespace, t.Name, t.Spec.CronSpec, t.Spec.Image)
	}

	// clientset := k.handle.ClientSet()
	// data, err := clientset.AppsV1().RESTClient().
	// 	Get().
	// 	AbsPath("/apis/stable.example.com/v1/namespaces/*/networkinfos/my-new-network-info-object").
	// 	// AbsPath("/apis/stable.example.com/v1").
	// 	// Namespace("default").
	// 	// Resource("networkinfos").
	// 	// Name("my-new-network-info-object").
	// 	DoRaw(context.TODO())

	// var config *rest.Config
	// var err error
	// // v1alpha1.AddToScheme(scheme.Scheme)
	// config, err = rest.InClusterConfig()
	// if err != nil {
	// 	klog.ErrorS(err, "====****Failed to get generate config****====")
	// }
	// crdConfig := *config
	// crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: "stable.example.com", Version: "v1"}
	// crdConfig.APIPath = "/apis"
	// crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	// crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()
	// exampleRestClient, err := rest.UnversionedRESTClientFor(&crdConfig)
	// data, err := exampleRestClient.
	// 	Get().
	// 	Resource("networkinfos").
	// 	DoRaw(context.TODO())
	if err != nil {
		klog.ErrorS(err, "====****Failed to get custom resource object****====")
		return 0, framework.AsStatus(fmt.Errorf("getting custom resource for node %q: %w", nodeName, err))
	}
	klog.InfoS("====Get success====")
	// klog.Info(data)
	score := rand.Int()
	klog.Infof("====Score = %d====", score)
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
