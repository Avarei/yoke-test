package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/yokecd/yoke/pkg/apis/airway/v1alpha1"
	"github.com/yokecd/yoke/pkg/openapi"

	listv1alpha1 "github.com/avarei/yoke-test/mylist/v1alpha1"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var flight string
	flag.StringVar(&flight, "flight", "", "wasm URL of the v1alpha1 Flight for ")
	flag.Parse()

	if flight == "" {
		return errors.New("flight URL not specified")
	}

	return json.NewEncoder(os.Stdout).Encode(v1alpha1.Airway{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mylists.example.com",
		},
		Spec: v1alpha1.AirwaySpec{
			Mode: v1alpha1.AirwayModeDynamic,
			WasmURLs: v1alpha1.WasmURLs{
				Flight: flight,
			},
			Template: apiextv1.CustomResourceDefinitionSpec{
				Group: "example.com",
				Names: apiextv1.CustomResourceDefinitionNames{
					Plural:   "mylists",
					Singular: "mylist",
					Kind:     "MyList",
				},
				Scope: apiextv1.NamespaceScoped,
				Versions: []apiextv1.CustomResourceDefinitionVersion{
					{
						Name:    strings.Split(listv1alpha1.APIVersion, "/")[1],
						Served:  true,
						Storage: true,
						Schema: &apiextv1.CustomResourceValidation{
							OpenAPIV3Schema: openapi.SchemaFrom(reflect.TypeFor[listv1alpha1.MyList]()),
						},
					},
				},
			},
		},
	})
}
