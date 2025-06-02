package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/yokecd/yoke/pkg/flight"
	"github.com/yokecd/yoke/pkg/helm"

	v1alpha1 "github.com/avarei/yoke-test/cluster/v1alpha1"
)

//go:embed vcluster-0.25.0.tgz
var chart []byte

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	cluster := &v1alpha1.Cluster{}
	if err := yaml.NewYAMLToJSONDecoder(os.Stdin).Decode(cluster); err != nil && err != io.EOF {
		return err
	}

	switch cluster.Spec.Type {
	case v1alpha1.ClusterTypeVCluster:
		var err error
		resources, err := createVClusterHelm(cluster)
		if err != nil {
			return err
		}
		return json.NewEncoder(os.Stdout).Encode(resources)

	default:
		return errors.New("not yet implemented cluster type")
	}
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
