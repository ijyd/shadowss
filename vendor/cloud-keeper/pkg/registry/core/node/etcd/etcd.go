package etcd

import (
	etcdregistry "apistack/pkg/registry/generic/registry/etcds"

	"apistack/pkg/registry/generic"
	"apistack/pkg/registry/generic/registry"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/node"

	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/runtime"
)

// REST implements the REST endpoint for usertoken
type REST struct {
	*etcdregistry.Store
}

// NewREST returns a RESTStorage object that will work with testtype.
func NewREST(opts generic.RESTOptions) *REST {
	prefix := "/" + opts.ResourcePrefix
	newListFunc := func() runtime.Object { return &api.NodeList{} }

	storageConfig := opts.StorageConfig
	storageInterface, dFunc := generic.NewRawStorage(storageConfig)

	store := &registry.Store{
		NewFunc:     func() runtime.Object { return &api.Node{} },
		NewListFunc: newListFunc,
		KeyRootFunc: func(ctx freezerapi.Context) string {
			return prefix
		},
		KeyFunc: func(ctx freezerapi.Context, name string) (string, error) {
			return registry.NoNamespaceKeyFunc(ctx, prefix, name)
		},
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*api.Node).Name, nil
		},
		PredicateFunc:           node.MatchNode,
		QualifiedResource:       api.Resource("nodes"),
		EnableGarbageCollection: opts.EnableGarbageCollection,
		DeleteCollectionWorkers: opts.DeleteCollectionWorkers,

		CreateStrategy:      node.Strategy,
		UpdateStrategy:      node.Strategy,
		DeleteStrategy:      node.Strategy,
		ReturnDeletedObject: true,
		// AfterCreate:         node.PadObj,
		// AfterUpdate:         node.PadObj,
		// AfterDelete:         node.PadObj,

		Storage:     storageInterface,
		DestroyFunc: dFunc,
	}

	return &REST{
		etcdregistry.NewStore(*store),
	}

}
