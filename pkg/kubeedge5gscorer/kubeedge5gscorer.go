package kubeedge5gscorer

import (
	"context"
	"fmt"
	"math/rand"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
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
	klog.InfoS("====KubeEdge5GScorer starts====")
	klog.InfoS("====Try to get custom resource object====")
	clientset := k.handle.ClientSet()
	data, err := clientset.AppsV1().RESTClient().
		Get().
		AbsPath("/apis/stable.example.com/v1").
		Namespace("default").
		Resource("networkinfos").
		Name("my-new-network-info-object").
		DoRaw(context.TODO())
	if err != nil {
		klog.ErrorS(err, "====****Failed to get custom resource object****====")
		return 0, framework.AsStatus(fmt.Errorf("getting custom resource for node %q: %w", nodeName, err))
	}
	klog.InfoS("====Get success====")
	klog.Info(data)
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
