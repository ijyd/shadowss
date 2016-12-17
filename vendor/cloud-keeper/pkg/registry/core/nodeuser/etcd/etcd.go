package etcd

import (
	etcdregistry "apistack/pkg/registry/generic/registry/etcds"
	"sync"

	"time"

	"github.com/golang/glog"

	"apistack/pkg/registry/generic"
	"apistack/pkg/registry/generic/registry"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/cache"
	"cloud-keeper/pkg/registry/core/node"
	"cloud-keeper/pkg/registry/core/nodeuser"
	"cloud-keeper/pkg/registry/core/user"

	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/errors"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/watch"
)

type nodeUserSync struct {
	mutex sync.Mutex
}

// REST implements the REST endpoint for usertoken
type REST struct {
	store *etcdregistry.Store

	nodeRegistry node.Registry
	user         user.Registry

	nodeUserLock *cache.LRUExpireCache
}

const (
	//ttl = 1800
	ttl = 120
)

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
		PredicateFunc: nodeuser.MatchNodeUser,
		TTLFunc: func(runtime.Object, uint64, bool) (uint64, error) {
			return ttl, nil
		},
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
		store:        etcdregistry.NewStore(*store),
		nodeUserLock: cache.NewLRUExpireCache(256),
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

func (r *REST) Get(ctx freezerapi.Context, name string) (runtime.Object, error) {
	return r.store.Get(ctx, name)
}

func (r *REST) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {

	obj, err := objInfo.UpdatedObject(ctx, nil)
	if err != nil {
		return nil, false, err
	}

	nodeuser := obj.(*api.NodeUser)

	//merge exist node user lables
	nodeuser.Labels = make(map[string]string)
	oldObj, err := r.Get(ctx, name)
	if err == nil {
		oldUser := oldObj.(*api.NodeUser)
		for k, v := range oldUser.Labels {
			nodeuser.Labels[k] = v
		}
	}

	switch nodeuser.Spec.Phase {
	case api.NodeUserPhaseAdd:
		//nodeuser.Spec.Phase = api.NodeUserPhaseAdd
		nodeuser.Labels[nodeuser.Spec.NodeName] = api.NodeUserPhaseAdd
	case api.NodeUserPhaseDelete:
		//nodeuser.Spec.Phase = api.NodeUserPhaseDelete
		nodeuser.Labels[nodeuser.Spec.NodeName] = api.NodeUserPhaseDelete
	case api.NodeUserPhaseUpdate:
		userSyncObj, ok := r.nodeUserLock.Get(name)
		var userSync *nodeUserSync
		if !ok {
			userSync = &nodeUserSync{}
			r.nodeUserLock.Add(name, userSync, time.Duration(time.Second*5))
		} else {
			userSync = userSyncObj.(*nodeUserSync)
		}

		userSync.mutex.Lock()
		defer userSync.mutex.Unlock()

		glog.V(5).Infof("update node user %v \r\n", *nodeuser)
		err = r.user.UpdateUserByNodeUser(ctx, nodeuser)
		if err != nil {
			return nil, false, err
		}
		return nodeuser, true, nil
	}

	glog.V(5).Infof("update node user %v \r\n", *nodeuser)
	return r.store.Update(ctx, name, rest.DefaultUpdatedObjectInfo(nodeuser, api.Scheme))
}

func (r *REST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {
	if name, ok := options.FieldSelector.RequiresExactMatch("metadata.name"); ok {
		options := &freezerapi.ListOptions{}
		users, err := r.user.ListUserByNodeName(ctx, name, options)
		if err != nil {
			return nil, err
		}

		nodeUserList := &api.NodeUserList{}
		for _, v := range users.Items {

			userNodeInfo, ok := v.Spec.UserService.Nodes[name]
			if !ok {
				glog.Warningf("not found node(%s) in user(%s), but in filter list result?\r\n", name, v.Name)
				continue
			}

			nodeUser := api.NodeUser{}
			nodeUser.Name = v.Name
			nodeUser.Spec.NodeName = name
			nodeUser.Spec.User = userNodeInfo.User
			//nodeUser.Spec.Phase = userNodeInfo.
			nodeUserList.Items = append(nodeUserList.Items, nodeUser)
		}

		return nodeUserList, nil
	} else {
		return nil, errors.NewBadRequest("need a 'metadata.name' filed selector")
	}

}

func (r *REST) Watch(ctx freezerapi.Context, options *freezerapi.ListOptions) (watch.Interface, error) {
	w, err := r.store.Watch(ctx, options)
	glog.Infof("watch node user result:%v\r\n", err)
	return w, err
}
