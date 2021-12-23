package resources

import (
	"context"
	"github.com/itchyny/gojq"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type Generic struct {
	Dyn       dynamic.Interface
	Ctx       context.Context
	Group     string
	Version   string
	Kind      string
	Namespace string
	Jq        string
}

func NewGeneric(ctx context.Context, dynamic dynamic.Interface, group, version, kind, namespace, jq string) *Generic {
	return &Generic{
		Dyn:       dynamic,
		Ctx:       ctx,
		Group:     group,
		Version:   version,
		Kind:      kind,
		Namespace: namespace,
		Jq:        jq,
	}
}

func (g *Generic) GetResourcesDynamically() ([]unstructured.Unstructured, error) {

	resourceId := schema.GroupVersionResource{
		Group:    g.Group,
		Version:  g.Version,
		Resource: g.Kind,
	}
	list, err := g.Dyn.Resource(resourceId).Namespace(g.Namespace).
		List(g.Ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (g *Generic) GetResourcesByJq() ([]unstructured.Unstructured, error) {

	resources := make([]unstructured.Unstructured, 0)

	query, err := gojq.Parse(g.Jq)
	if err != nil {
		return nil, err
	}

	items, err := g.GetResourcesDynamically()
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		// Convert object to raw JSON
		var rawJson interface{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, &rawJson)
		if err != nil {
			return nil, err
		}

		// Evaluate jq against JSON

		iter := query.Run(rawJson)
		for {
			result, ok := iter.Next()
			if !ok {
				break
			}
			if err, ok := result.(error); ok {
				if err != nil {
					return nil, err
				}
			} else {
				resources = append(resources, item)
			}
		}

	}
	return resources, nil
}
