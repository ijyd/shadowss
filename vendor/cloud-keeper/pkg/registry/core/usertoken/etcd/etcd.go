package etcd

import (
	"apistack/pkg/registry/generic"
	"apistack/pkg/registry/generic/registry"
	etcdregistry "apistack/pkg/registry/generic/registry/etcds"

	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/runtime"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/usertoken"
)

// REST implements the REST endpoint for usertoken
type REST struct {
	store *etcdregistry.Store
}

const (
	ttl = 600
)

// NewREST returns a RESTStorage object that will work with testtype.
func NewREST(opts generic.RESTOptions) *REST {
	prefix := "/" + opts.ResourcePrefix
	newListFunc := func() runtime.Object { return &api.UserTokenList{} }

	storageConfig := opts.StorageConfig
	storageInterface, dFunc := generic.NewRawStorage(storageConfig)

	store := &registry.Store{
		NewFunc:     func() runtime.Object { return &api.UserToken{} },
		NewListFunc: newListFunc,
		KeyRootFunc: func(ctx freezerapi.Context) string {
			return prefix
		},
		KeyFunc: func(ctx freezerapi.Context, name string) (string, error) {
			return registry.NoNamespaceKeyFunc(ctx, prefix, name)
		},
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*api.UserToken).Name, nil
		},
		TTLFunc: func(runtime.Object, uint64, bool) (uint64, error) {
			return ttl, nil
		},
		PredicateFunc:           usertoken.MatchUserToken,
		QualifiedResource:       api.Resource("usertokens"),
		EnableGarbageCollection: opts.EnableGarbageCollection,
		DeleteCollectionWorkers: opts.DeleteCollectionWorkers,

		CreateStrategy:      usertoken.Strategy,
		UpdateStrategy:      usertoken.Strategy,
		DeleteStrategy:      usertoken.Strategy,
		ReturnDeletedObject: true,

		Storage:     storageInterface,
		DestroyFunc: dFunc,
	}

	return &REST{
		etcdregistry.NewStore(*store),
	}
}

func (r *REST) New() runtime.Object {
	return &api.UserToken{}
}

func (r *REST) NewList() runtime.Object {
	return &api.UserToken{}
}

func (r *REST) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	return r.store.Update(ctx, name, objInfo)
}

func (r *REST) Get(ctx freezerapi.Context, name string) (runtime.Object, error) {
	return r.store.Get(ctx, name)
}

func (r *REST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {
	return r.store.List(ctx, options)
}
