package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"

	applicationv1alpha1 "github.com/avarei/yoke-test/argocd-aoa/apis/v1alpha1"
	"github.com/yokecd/yoke/pkg/flight"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/utils/ptr"
)

func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type ArgoCDInput struct {
	Revision string `json:"revision"`
}

func run(stdin io.Reader, stdout io.Writer) error {
	input := &ArgoCDInput{}
	if err := yaml.NewYAMLToJSONDecoder(stdin).Decode(input); err != nil && err != io.EOF {
		return err
	}

	resources, err := reconcile(input)
	if err != nil {
		return err
	}
	return json.NewEncoder(stdout).Encode(resources)
}

var (
	flightClusterImage string = "ghcr.io/avarei/yoke-test/flight-cluster"
	repoUrl            string = "https://github.com/avarei/yoke-test"
)

func reconcile(input *ArgoCDInput) ([]applicationv1alpha1.Application, error) {
	revision := input.Revision

	if revision == "" {
		return nil, fmt.Errorf("expected --revision to be set")
	}

	var apps []applicationv1alpha1.Application
	appCluster, err := createAppCluster(revision)
	if err != nil {
		return nil, err
	}

	appAtc, err := createAppATC()
	if err != nil {
		return nil, err
	}

	apps = append(apps, appCluster, appAtc)

	return apps, nil
}

func createAppATC() (applicationv1alpha1.Application, error) {
	return applicationv1alpha1.Application{
		ObjectMeta: v1.ObjectMeta{
			Name:      "atc",
			Namespace: flight.Namespace(),
		},
		Spec: applicationv1alpha1.ApplicationSpec{
			Project: "default",
			Sources: applicationv1alpha1.ApplicationSources{
				applicationv1alpha1.ApplicationSource{
					RepoURL:        "https://github.com/yokecd/yoke",
					Path:           "./cmd/atc-installer",
					TargetRevision: "atc-installer/v0.12.3",
					Plugin: &applicationv1alpha1.ApplicationSourcePlugin{
						Name: "yokecd",
						Parameters: applicationv1alpha1.ApplicationSourcePluginParameters{
							applicationv1alpha1.ApplicationSourcePluginParameter{
								Name:    "build",
								String_: ptr.To("true"),
							},
						},
					},
				},
			},
			Destination: applicationv1alpha1.ApplicationDestination{
				Name:      "in-cluster",
				Namespace: "atc",
			},
			SyncPolicy: &applicationv1alpha1.SyncPolicy{
				Automated: &applicationv1alpha1.SyncPolicyAutomated{
					Prune:    true,
					SelfHeal: true,
				},
			},
		},
	}, nil
}

func createAppCluster(revision string) (applicationv1alpha1.Application, error) {
	return applicationv1alpha1.Application{
		TypeMeta: v1.TypeMeta{
			APIVersion: "argoproj.io/v1alpha1",
			Kind:       "Application",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "cluster",
			Namespace: flight.Namespace(),
		},
		Spec: applicationv1alpha1.ApplicationSpec{
			Project: "default",
			Sources: applicationv1alpha1.ApplicationSources{
				applicationv1alpha1.ApplicationSource{
					RepoURL:        repoUrl,
					Path:           "./cluster/airway",
					TargetRevision: "main",
					Plugin: &applicationv1alpha1.ApplicationSourcePlugin{
						Name: "yokecd",
						Parameters: applicationv1alpha1.ApplicationSourcePluginParameters{
							applicationv1alpha1.ApplicationSourcePluginParameter{
								Name:    "build",
								String_: ptr.To("true"),
							},
							applicationv1alpha1.ApplicationSourcePluginParameter{
								Name: "args",
								OptionalArray: &applicationv1alpha1.OptionalArray{
									Array: []string{fmt.Sprintf("--flight=%s:v0.0.0-%s", flightClusterImage, revision)},
								},
							},
						},
					},
				},
			},
			Destination: applicationv1alpha1.ApplicationDestination{
				Name:      "in-cluster",
				Namespace: "argocd",
			},
			SyncPolicy: &applicationv1alpha1.SyncPolicy{
				Automated: &applicationv1alpha1.SyncPolicyAutomated{
					Prune:    true,
					SelfHeal: true,
				},
			},
		},
	}, nil
}
