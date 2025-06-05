package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"k8s.io/apimachinery/pkg/util/yaml"

	v1alpha1 "github.com/avarei/yoke-test/mylist/v1alpha1"
)

func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(stdin io.Reader, stdout io.Writer) error {
	cluster := &v1alpha1.MyList{}
	if err := yaml.NewYAMLToJSONDecoder(stdin).Decode(cluster); err != nil && err != io.EOF {
		return err
	}
	resources, err := reconcile(cluster)
	if err != nil {
		return err
	}
	return json.NewEncoder(stdout).Encode(resources)
}

func reconcile(cluster *v1alpha1.MyList) ([]any, error) {
	var resources []any
	resources = append(resources, cluster)

	for _, resource := range cluster.Spec.Items {
		resources = append(resources, resource.Extras)
	}

	return resources, nil
}
