package main

import (
	"testing"
	"time"

	v1alpha1 "github.com/avarei/yoke-test/cluster/v1alpha1"
	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRun(t *testing.T) {
	cluster := &v1alpha1.Cluster{
		TypeMeta: v1.TypeMeta{
			APIVersion: v1alpha1.APIVersion,
			Kind:       v1alpha1.KindCluster,
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "foo",
			Namespace: "bar",
		},
		Spec: v1alpha1.ClusterSpec{
			Type: v1alpha1.ClusterTypeVCluster,
		},
	}

	resources, err := reconcile(cluster)
	if err != nil {
		t.Error(err)
	}

	for _, resource := range resources {
		outCluster, ok := resource.(*v1alpha1.Cluster)
		if !ok {
			continue
		}

		want := &v1alpha1.Cluster{
			TypeMeta: v1.TypeMeta{
				APIVersion: v1alpha1.APIVersion,
				Kind:       v1alpha1.KindCluster,
			},
			ObjectMeta: v1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
			},
			Spec: v1alpha1.ClusterSpec{
				Type: v1alpha1.ClusterTypeVCluster,
			},
		}
		meta.SetStatusCondition((*[]v1.Condition)(&want.Status.Conditions), v1.Condition{
			Type:    "Ready",
			Status:  v1.ConditionUnknown,
			Message: "I am curious what happens next",
		})
		time := v1.Date(2025, 06, 03, 0, 0, 0, 0, time.UTC)

		for i := range outCluster.Status.Conditions {
			outCluster.Status.Conditions[i].LastTransitionTime = time
		}

		want.Status.Conditions[0].LastTransitionTime = time

		if diff := cmp.Diff(want, outCluster); diff != "" {
			t.Errorf("\n%s\nreconcile(...): -want, +got:\n%s", "fooooo", diff)
		}

	}
}
