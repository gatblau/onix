package server

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

type K8S struct {
	cfg     *rest.Config
	decoder runtime.Serializer
}

func NewK8S() (*K8S, error) {
	config, err := getKubeConfig()
	if err != nil {
		return nil, err
	}
	return &K8S{
		cfg:     config,
		decoder: yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme),
	}, nil
}

func (k *K8S) apply(yamlResource string, ctx context.Context) error {
	// 1. Prepare a RESTMapper to find GVR
	dc, err := discovery.NewDiscoveryClientForConfig(k.cfg)
	if err != nil {
		return err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// 2. Prepare the dynamic client
	dyn, err := dynamic.NewForConfig(k.cfg)
	if err != nil {
		return err
	}

	// 3. Decode YAML manifest into unstructured.Unstructured
	obj := &unstructured.Unstructured{}
	_, gvk, err := k.decoder.Decode([]byte(yamlResource), nil, obj)
	if err != nil {
		return err
	}

	// 4. Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	// 5. Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}

	// 6. Marshal object into JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	// 7. Create or Update the object with SSA
	//     types.ApplyPatchType indicates SSA.
	//     FieldManager specifies the field owner ID.
	_, err = dr.Patch(ctx, obj.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{
		FieldManager: "build-manager-controller",
	})
	return err
}

func (k *K8S) NewImagePipeline(appName, namespace string) error {
	pr, err := merge(imageBuilderRun, newPipelineRunData(appName, namespace))
	if err != nil {
		return fmt.Errorf("cannot merge image pipeline template: %v", err)
	}
	err = k.apply(pr, context.Background())
	if err != nil {
		return fmt.Errorf("cannot apply pipeline run: %v", err)
	}
	return nil
}
