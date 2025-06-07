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
	"k8s.io/utils/ptr"
)

func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(_ io.Reader, stdout io.Writer) error {
	resources, err := reconcile()
	if err != nil {
		return err
	}
	return json.NewEncoder(stdout).Encode(resources)
}

func reconcile() ([]applicationv1alpha1.Application, error) {
	var apps []applicationv1alpha1.Application
	appCluster, err := createAppAppOfApps()
	if err != nil {
		return nil, err
	}

	apps = append(apps, appCluster)

	return apps, nil
}

func createAppAppOfApps() (applicationv1alpha1.Application, error) {
	return applicationv1alpha1.Application{
		TypeMeta: v1.TypeMeta{
			APIVersion: "argoproj.io/v1alpha1",
			Kind:       "Application",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "app-of-apps",
			Namespace: flight.Namespace(),
		},
		Spec: applicationv1alpha1.ApplicationSpec{
			Project: "default",
			Sources: applicationv1alpha1.ApplicationSources{
				applicationv1alpha1.ApplicationSource{
					RepoURL:        "https://github.com/avarei/yoke-test",
					Path:           "./argocd-aoa/app-of-apps",
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
									Array: []string{"--revision=$ARGOCD_APP_REVISION_SHORT_8"},
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
