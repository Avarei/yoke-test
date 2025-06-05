package v1alpha1

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/yokecd/yoke/pkg/flight"
	"github.com/yokecd/yoke/pkg/openapi"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

const (
	APIVersion = "example.com/v1alpha1"
	KindMyList = "MyList"
)

// MyList is used to test what happens when two resources from different flights try to manage the same resource
// This will be tested in three ways:
// 1. Both times the resource to be created is identical
// 2. Both Flights control different fields of the resource
// 3. Both Flights want to create conflicting resources
type MyList struct {
	metav1.TypeMeta
	metav1.ObjectMeta `json:"metadata"`
	Spec              MyListSpec   `json:"spec"`
	Status            MyListStatus `json:"status,omitzero"`
}

type MyListSpec struct {
	Items []Item `json:"items"`
}

type Item struct {
	// metav1.TypeMeta   `json:",inline"`
	APIVersion        string `json:"apiVersion"`
	Kind              string `json:"kind"`
	metav1.ObjectMeta `json:"metadata"`
	Extras            map[string]any `json:"-"`
}

func (m *Item) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	m.Extras = raw
	return nil
}
func (m *Item) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Extras)
}

type MyListStatus struct {
	Conditions flight.Conditions `json:"conditions,omitempty"`
}

func (Item) OpenAPISchema() *apiext.JSONSchemaProps {
	type crd Item
	schema := openapi.SchemaFrom(reflect.TypeFor[crd]())
	schema.XPreserveUnknownFields = ptr.To(true)
	schema.Description = "define a kubernetes manifest here"
	return schema
}

// Custom Marshalling Logic so that users do not need to explicity fill out the Kind and ApiVersion.
func (cluster MyList) MarshalJSON() ([]byte, error) {
	cluster.Kind = KindMyList
	cluster.APIVersion = APIVersion

	type clusterAlt MyList
	return json.Marshal(clusterAlt(cluster))
}

// Custom Unmarshalling to raise an error if the ApiVersion or Kind does not match.
func (cluster *MyList) UnmarshalJSON(data []byte) error {
	type ClusterAlt MyList
	if err := json.Unmarshal(data, (*ClusterAlt)(cluster)); err != nil {
		return err
	}
	if cluster.APIVersion != APIVersion {
		return fmt.Errorf("unexpected api version: expected %s but got %s", APIVersion, cluster.APIVersion)
	}
	if cluster.Kind != KindMyList {
		return fmt.Errorf("unexpected kind: expected %s but got %s", KindMyList, cluster.Kind)
	}
	return nil
}
