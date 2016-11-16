package etcd

import (
	"fmt"

	"apistack/pkg/registry/generic"
	"apistack/pkg/registry/generic/registry"
	etcdregistry "apistack/pkg/registry/generic/registry/etcds"

	"apistack/examples/apiserver/pkg/apis/testgroup"
	"apistack/examples/apiserver/pkg/registry/testgroup/testtype"

	api "gofreezer/pkg/api"
	"gofreezer/pkg/fields"
	"gofreezer/pkg/labels"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/storage"
)

type REST struct {
	*etcdregistry.Store
	//*registry.Store
}

// NewREST returns a RESTStorage object that will work with testtype.
func NewREST(opts generic.RESTOptions) *REST {
	prefix := opts.ResourcePrefix
	newListFunc := func() runtime.Object { return &testgroup.TestTypeList{} }

	storageConfig := opts.StorageConfig
	storageInterface, dFunc := generic.NewRawStorage(storageConfig)

	store := &registry.Store{
		NewFunc: func() runtime.Object { return &testgroup.TestType{} },
		// NewListFunc returns an object capable of storing results of an etcd list.
		NewListFunc: newListFunc,
		// Produces a path that etcd understands, to the root of the resource
		// by combining the namespace in the context with the given prefix.
		KeyRootFunc: func(ctx api.Context) string {
			return registry.NamespaceKeyRootFunc(ctx, prefix)
		},
		// Produces a path that etcd understands, to the resource by combining
		// the namespace in the context with the given prefix.
		KeyFunc: func(ctx api.Context, name string) (string, error) {
			return registry.NamespaceKeyFunc(ctx, prefix, name)
		},
		// Retrieve the name field of the resource.
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*testgroup.TestType).Name, nil
		},
		// Used to match objects based on labels/fields for list.
		PredicateFunc: func(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
			return storage.SelectionPredicate{
				Label: label,
				Field: field,
				GetAttrs: func(obj runtime.Object) (labels.Set, fields.Set, error) {
					testType, ok := obj.(*testgroup.TestType)
					if !ok {
						return nil, nil, fmt.Errorf("unexpected type of given object")
					}
					return labels.Set(testType.ObjectMeta.Labels), fields.Set{}, nil
				},
			}
		},

		CreateStrategy:      testtype.Strategy,
		UpdateStrategy:      testtype.Strategy,
		DeleteStrategy:      testtype.Strategy,
		ReturnDeletedObject: true,

		Storage:     storageInterface,
		DestroyFunc: dFunc,
	}

	return &REST{etcdregistry.NewStore(*store)}
}
