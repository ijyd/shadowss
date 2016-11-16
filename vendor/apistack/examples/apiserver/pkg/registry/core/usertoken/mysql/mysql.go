package mysql

import (
	mysqlregistry "apistack/pkg/registry/generic/registry/mysqls"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/storage/storagebackend"

	"apistack/pkg/registry/generic"
	"apistack/pkg/registry/generic/registry"

	"apistack/examples/apiserver/pkg/api"
	"apistack/examples/apiserver/pkg/registry/core/usertoken"

	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/runtime"
)

// REST implements the REST endpoint for usertoken
type REST struct {
	store *mysqlregistry.Store
}

// NewREST returns a RESTStorage object that will work with testtype.
func NewREST(opts generic.RESTOptions) *REST {
	prefix := "/" + opts.ResourcePrefix
	newListFunc := func() runtime.Object { return &api.UserTokenList{} }

	storageConfig := opts.StorageConfig
	storageConfig.Type = storagebackend.StorageTypeMysql
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
		PredicateFunc:           usertoken.MatchUserToken,
		QualifiedResource:       api.Resource("usertokens"),
		EnableGarbageCollection: opts.EnableGarbageCollection,
		DeleteCollectionWorkers: opts.DeleteCollectionWorkers,

		CreateStrategy:      usertoken.Strategy,
		UpdateStrategy:      usertoken.Strategy,
		DeleteStrategy:      usertoken.Strategy,
		ReturnDeletedObject: true,
		AfterCreate:         usertoken.PadObj,
		AfterUpdate:         usertoken.PadObj,
		AfterDelete:         usertoken.PadObj,

		Storage:     storageInterface,
		DestroyFunc: dFunc,
	}

	return &REST{
		store: mysqlregistry.NewStore(*store),
	}

}

func (r *REST) New() runtime.Object {
	return &api.UserToken{}
}

// func (r *REST) Create(ctx freezerapi.Context, obj runtime.Object) (runtime.Object, error) {
// 	//check user info is right
// 	return r.store.Create(ctx, obj)
// }
func (r *REST) NewList() runtime.Object {
	return &api.UserTokenList{}
}

func (r *REST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {
	return r.store.List(ctx, options)
}

func (r *REST) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	return r.store.Update(ctx, name, objInfo)
}

func (r *REST) Get(ctx freezerapi.Context, name string) (runtime.Object, error) {
	return r.store.Get(ctx, name)
}
