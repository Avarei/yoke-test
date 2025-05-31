package v1

import (
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	APIVersion  = "example.com/v1alpha1"
	KindCluster = "Cluster"
)

type Cluster struct {
	metav1.TypeMeta
	metav1.ObjectMeta `json:"metadata"`
	Spec              ClusterSpec `json:"spec"`
}

// Our Backend Specification
type ClusterSpec struct {
	// ClusterType specifies what kind of cluster to deploy it
	Type ClusterType `json:"type" Enum:"vCluster,ClusterAPI,Gardener"`
}

type ClusterType string

var (
	ClusterTypeVCluster   ClusterType = "vCluster"
	ClusterTypeClusterAPI ClusterType = "ClusterAPI"
	ClusterTypeGardener   ClusterType = "Gardener"
)

// Custom Marshalling Logic so that users do not need to explicity fill out the Kind and ApiVersion.
func (cluster Cluster) MarshalJSON() ([]byte, error) {
	cluster.Kind = KindCluster
	cluster.APIVersion = APIVersion

	type clusterAlt Cluster
	return json.Marshal(clusterAlt(cluster))
}

// Custom Unmarshalling to raise an error if the ApiVersion or Kind does not match.
func (cluster *Cluster) UnmarshalJSON(data []byte) error {
	type ClusterAlt Cluster
	if err := json.Unmarshal(data, (*ClusterAlt)(cluster)); err != nil {
		return err
	}
	if cluster.APIVersion != APIVersion {
		return fmt.Errorf("unexpected api version: expected %s but got %s", APIVersion, cluster.APIVersion)
	}
	if cluster.Kind != KindCluster {
		return fmt.Errorf("unexpected kind: expected %s but got %s", KindCluster, cluster.Kind)
	}
	return nil
}
