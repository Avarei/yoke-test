package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	meta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/yokecd/yoke/pkg/flight"
	"github.com/yokecd/yoke/pkg/helm"

	v1alpha1 "github.com/avarei/yoke-test/cluster/v1alpha1"
)

//go:embed vcluster-0.25.0.tgz
var chart []byte

func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(stdin io.Reader, stdout io.Writer) error {
	cluster := &v1alpha1.Cluster{}
	if err := yaml.NewYAMLToJSONDecoder(stdin).Decode(cluster); err != nil && err != io.EOF {
		return err
	}
	resources, err := reconcile(cluster)
	if err != nil {
		return err
	}
	return json.NewEncoder(stdout).Encode(resources)
}

func reconcile(cluster *v1alpha1.Cluster) ([]any, error) {
	if !meta.IsStatusConditionPresentAndEqual(cluster.Status.Conditions, "Ready", metav1.ConditionTrue) {
		meta.SetStatusCondition((*[]metav1.Condition)(&cluster.Status.Conditions), metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionTrue,
			Reason:             "Everything is Good",
			Message:            "Weather is nice...",
			ObservedGeneration: cluster.Generation,
		})
	}

	cluster.Status = v1alpha1.ClusterStatus{
		Conditions: []metav1.Condition{
			{
				Type:               "Ready",
				Status:             metav1.ConditionUnknown,
				Message:            "I am curious what happens next",
				ObservedGeneration: cluster.ObjectMeta.Generation,
				LastTransitionTime: metav1.Now(),
			},
		},
	}
	var resources []any

	switch cluster.Spec.Type {
	case v1alpha1.ClusterTypeVCluster:
		var err error
		vclusterResources, err := createVClusterHelm(cluster)
		if err != nil {
			return nil, err
		}
		for _, resource := range vclusterResources {
			resources = append(resources, resource)
		}

	default:
		return nil, errors.New("not yet implemented cluster type")
	}

	resources = append(resources, cluster)

	return resources, nil
}

func createVClusterHelm(cluster *v1alpha1.Cluster) ([]*unstructured.Unstructured, error) {
	chart, err := helm.LoadChartFromZippedArchive(chart)
	if err != nil {
		return nil, err
	}

	return chart.Render(
		cluster.GetName(),
		flight.Namespace(),
		map[string]any{},
	)
}
