package mysql

import (
	mysqlregistry "apistack/pkg/registry/generic/registry/mysqls"
	"gofreezer/pkg/storage/storagebackend"

	"apistack/pkg/registry/generic"
	"apistack/pkg/registry/generic/registry"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/account"

	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/runtime"
)

// REST implements the REST endpoint for usertoken
type REST struct {
	*mysqlregistry.Store
}

// NewREST returns a RESTStorage object that will work with testtype.
func NewREST(opts generic.RESTOptions) *REST {
	prefix := "/" + opts.ResourcePrefix
	newListFunc := func() runtime.Object { return &api.AccountList{} }

	storageConfig := opts.StorageConfig
	storageConfig.Type = storagebackend.StorageTypeMysql
	storageInterface, dFunc := generic.NewRawStorage(storageConfig)

	store := &registry.Store{
		NewFunc:     func() runtime.Object { return &api.Account{} },
		NewListFunc: newListFunc,
		KeyRootFunc: func(ctx freezerapi.Context) string {
			return prefix
		},
		KeyFunc: func(ctx freezerapi.Context, name string) (string, error) {
			return registry.NoNamespaceKeyFunc(ctx, prefix, name)
		},
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*api.Account).Name, nil
		},
		PredicateFunc:           account.MatchAccount,
		QualifiedResource:       api.Resource("accounts"),
		EnableGarbageCollection: opts.EnableGarbageCollection,
		DeleteCollectionWorkers: opts.DeleteCollectionWorkers,

		CreateStrategy:      account.Strategy,
		UpdateStrategy:      account.Strategy,
		DeleteStrategy:      account.Strategy,
		ReturnDeletedObject: true,
		AfterCreate:         account.PadObj,
		AfterUpdate:         account.PadObj,
		AfterDelete:         account.PadObj,

		Storage:     storageInterface,
		DestroyFunc: dFunc,
	}

	return &REST{mysqlregistry.NewStore(*store)}
}
