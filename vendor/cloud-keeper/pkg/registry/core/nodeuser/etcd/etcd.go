package etcd

import (
	etcdregistry "apistack/pkg/registry/generic/registry/etcds"

	"github.com/golang/glog"

	"apistack/pkg/registry/generic"
	"apistack/pkg/registry/generic/registry"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/node"
	"cloud-keeper/pkg/registry/core/nodeuser"
	"cloud-keeper/pkg/registry/core/user"

	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/errors"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/watch"
)

// REST implements the REST endpoint for usertoken
type REST struct {
	store *etcdregistry.Store

	nodeRegistry node.Registry
	user         user.Registry
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

func (r *REST) SetRequireRegistry(nodeReg node.Registry, user user.Registry) {
	r.user = user
	r.nodeRegistry = nodeReg

}

func (*REST) New() runtime.Object {
	return &api.NodeUser{}
}

func (*REST) NewList() runtime.Object {
	return &api.NodeUserList{}
}

func (r *REST) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {

	obj, err := objInfo.UpdatedObject(ctx, nil)
	if err != nil {
		return nil, false, err
	}

	nodeuser := obj.(*api.NodeUser)
	nodename := nodeuser.Spec.NodeName
	node, err := r.nodeRegistry.GetNode(ctx, nodename)
	userName := nodeuser.Spec.User.Name
	glog.V(5).Infof("update node user %+v \r\n", *node)
	switch nodeuser.Spec.Phase {
	case api.NodeUserPhaseAdd:
		glog.V(5).Infof("add node user %v \r\n", *nodeuser)
		node.Spec.Users[userName] = nodeuser.Spec
	case api.NodeUserPhaseDelete:
		glog.V(5).Infof("delete node user %v \r\n", *nodeuser)
		delete(node.Spec.Users, userName)
	case api.NodeUserPhaseUpdate:
		glog.V(5).Infof("update node user %v \r\n", *nodeuser)
		err = r.user.UpdateUserByNodeUser(ctx, nodeuser)
		if err != nil {
			return nil, false, err
		}
		glog.V(5).Infof("update node user done \r\n")
		node.Spec.Users[userName] = nodeuser.Spec
	}

	_, _, err = r.nodeRegistry.UpdateNode(ctx, nodename, rest.DefaultUpdatedObjectInfo(node, api.Scheme))
	if err != nil {
		return nil, false, err
	}

	return obj, true, nil
}

func (r *REST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {
	if name, ok := options.FieldSelector.RequiresExactMatch("metadata.name"); ok {
		node, err := r.nodeRegistry.GetNode(ctx, name)
		if err != nil {
			return nil, err
		}

		nodeUserList := &api.NodeUserList{}
		for _, v := range node.Spec.Users {
			nodeUser := api.NodeUser{}
			nodeUser.Name = v.User.Name
			nodeUser.Spec.NodeName = v.NodeName
			nodeUser.Spec.User = v.User
			nodeUser.Spec.Phase = v.Phase
			nodeUserList.Items = append(nodeUserList.Items, nodeUser)
		}

		return nodeUserList, nil
	} else {
		return nil, errors.NewBadRequest("need a 'metadata.name' filed selector")
	}

}

func (r *REST) Watch(ctx freezerapi.Context, options *freezerapi.ListOptions) (watch.Interface, error) {
	return r.store.Watch(ctx, options)
}
