package etcd

import (
	etcdregistry "apistack/pkg/registry/generic/registry/etcds"

	"apistack/pkg/registry/generic"
	"apistack/pkg/registry/generic/registry"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/nodeuser"

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
	newListFunc := func() runtime.Object { return &api.NodeUserList{} }

	storageConfig := opts.StorageConfig
	storageInterface, dFunc := generic.NewRawStorage(storageConfig)

	store := &registry.Store{
		NewFunc:     func() runtime.Object { return &api.NodeUser{} },
		NewListFunc: newListFunc,
		KeyRootFunc: func(ctx freezerapi.Context) string {
			return prefix
		},
		KeyFunc: func(ctx freezerapi.Context, name string) (string, error) {
			return registry.NoNamespaceKeyFunc(ctx, prefix, name)
		},
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*api.NodeUser).Name, nil
		},
		PredicateFunc:           nodeuser.MatchNodeUser,
		QualifiedResource:       api.Resource("nodeusers"),
		EnableGarbageCollection: opts.EnableGarbageCollection,
		DeleteCollectionWorkers: opts.DeleteCollectionWorkers,

		CreateStrategy:      nodeuser.Strategy,
		UpdateStrategy:      nodeuser.Strategy,
		DeleteStrategy:      nodeuser.Strategy,
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
