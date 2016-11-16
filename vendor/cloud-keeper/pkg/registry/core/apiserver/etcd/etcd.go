package etcd

import (
	etcdregistry "apistack/pkg/registry/generic/registry/etcds"

	"apistack/pkg/registry/generic"
	"apistack/pkg/registry/generic/registry"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/apiserver"

	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/runtime"
)

// REST implements the REST endpoint for usertoken
type REST struct {
	Store *etcdregistry.Store
}

// NewREST returns a RESTStorage object that will work with testtype.
func NewREST(opts generic.RESTOptions) *REST {
	prefix := "/" + opts.ResourcePrefix
	newListFunc := func() runtime.Object { return &api.APIServerList{} }

	storageConfig := opts.StorageConfig
	storageInterface, dFunc := generic.NewRawStorage(storageConfig)

	store := &registry.Store{
		NewFunc:     func() runtime.Object { return &api.APIServer{} },
		NewListFunc: newListFunc,
		KeyRootFunc: func(ctx freezerapi.Context) string {
			return prefix
		},
		KeyFunc: func(ctx freezerapi.Context, name string) (string, error) {
			return registry.NoNamespaceKeyFunc(ctx, prefix, name)
		},
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*api.APIServer).Name, nil
		},
		PredicateFunc:           apiserver.MatchAPIServer,
		QualifiedResource:       api.Resource("apiservers"),
		EnableGarbageCollection: opts.EnableGarbageCollection,
		DeleteCollectionWorkers: opts.DeleteCollectionWorkers,

		CreateStrategy:      apiserver.Strategy,
		UpdateStrategy:      apiserver.Strategy,
		DeleteStrategy:      apiserver.Strategy,
		ReturnDeletedObject: true,

		Storage:     storageInterface,
		DestroyFunc: dFunc,
	}

	return &REST{
		etcdregistry.NewStore(*store),
	}

}

func (*REST) New() runtime.Object {
	return &api.Node{}
}

func (*REST) NewList() runtime.Object {
	return &api.NodeList{}
}

func (rs *REST) Get(ctx freezerapi.Context, name string) (runtime.Object, error) {
	return rs.Store.Get(ctx, name)
}

func (rs *REST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {
	return rs.Store.List(ctx, options)
}
