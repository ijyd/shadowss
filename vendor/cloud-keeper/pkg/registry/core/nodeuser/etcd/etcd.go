package etcd

import (
	etcdregistry "apistack/pkg/registry/generic/registry/etcds"

	"github.com/golang/glog"

	"apistack/pkg/registry/generic"
	"apistack/pkg/registry/generic/registry"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/node"
	"cloud-keeper/pkg/registry/core/nodeuser"
	"cloud-keeper/pkg/registry/core/userservice"

	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/watch"
)

// REST implements the REST endpoint for usertoken
type REST struct {
	store *etcdregistry.Store

	userService  userservice.Registry
	nodeRegistry node.Registry
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
		store: etcdregistry.NewStore(*store),
	}
}

func (r *REST) SetRequireRegistry(userSrvReg userservice.Registry, nodeReg node.Registry) {
	r.userService = userSrvReg
	r.nodeRegistry = nodeReg
}

func (*REST) New() runtime.Object {
	return &api.NodeUser{}
}

// func (*REST) NewList() runtime.Object {
// 	return &api.NodeUserList{}
// }

// func (r *REST) Delete(ctx freezerapi.Context, name string, options *freezerapi.DeleteOptions) (runtime.Object, error) {
// 	return r.store.Delete(ctx, name, options)
// }
//
// func (r *REST) Create(ctx freezerapi.Context, obj runtime.Object) (runtime.Object, error) {
// 	return r.store.Create(ctx, obj)
// }

func (r *REST) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {

	obj, err := objInfo.UpdatedObject(ctx, nil)
	if err != nil {
		return nil, false, err
	}

	user := obj.(*api.NodeUser)
	nodename := user.Spec.NodeName
	node, err := r.nodeRegistry.GetNode(ctx, nodename)

	glog.V(5).Infof("update node user %v \r\n", *node)

	switch user.Spec.Phase {
	case api.NodeUserPhaseAdd:
		glog.V(5).Infof("add node user %v \r\n", *user)
		node.Spec.Users[user.Name] = *user
	case api.NodeUserPhaseDelete:
		glog.V(5).Infof("delete node user %v \r\n", user)
		delete(node.Spec.Users, user.Name)
	case api.NodeUserPhaseUpdate:
		glog.V(5).Infof("update node user %v \r\n", user)
		user := &api.UserService{}
		err = r.userService.UpdateUserServiceByNodeUser(ctx, user)
		if err != nil {
			return nil, false, err
		}
		return obj, true, nil
	}

	_, _, err = r.nodeRegistry.UpdateNode(ctx, nodename, rest.DefaultUpdatedObjectInfo(node, api.Scheme))
	if err != nil {
		return nil, false, err
	}

	return obj, true, nil
}

func (r *REST) Watch(ctx freezerapi.Context, options *freezerapi.ListOptions) (watch.Interface, error) {
	return r.store.Watch(ctx, options)
}
