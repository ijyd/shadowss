package mysql

import (
	mysqlregistry "apistack/pkg/registry/generic/registry/mysqls"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/storage/storagebackend"

	"apistack/pkg/registry/generic"
	"apistack/pkg/registry/generic/registry"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/node"

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
	newListFunc := func() runtime.Object { return &api.NodeList{} }

	storageConfig := opts.StorageConfig
	storageConfig.Type = storagebackend.StorageTypeMysql
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
		AfterCreate:         node.PadObj,
		AfterUpdate:         node.PadObj,
		AfterDelete:         node.PadObj,

		Storage:     storageInterface,
		DestroyFunc: dFunc,
	}

	return &REST{
		store: mysqlregistry.NewStore(*store),
	}

}

func (r *REST) New() runtime.Object {
	return &api.Node{}
}

func (r *REST) NewList() runtime.Object {
	return &api.NodeList{}
}

// func (r *REST) Create(ctx freezerapi.Context, obj runtime.Object) (runtime.Object, error) {
// 	//check user info is right
// 	return r.store.Create(ctx, obj)
// }

func (r *REST) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	return r.store.Update(ctx, name, objInfo)
}

func (r *REST) Get(ctx freezerapi.Context, name string) (runtime.Object, error) {
	obj, err := r.store.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	objNode := obj.(*api.Node)

	objNode.Name = objNode.Spec.Server.Name

	return objNode, nil
}

func (rs *REST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {
	return rs.store.List(ctx, options)
}
