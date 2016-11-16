package etcd

import (
	"fmt"
	"strconv"

	"apistack/pkg/registry/generic"
	"apistack/pkg/registry/generic/registry"
	etcdregistry "apistack/pkg/registry/generic/registry/etcds"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/account"
	"cloud-keeper/pkg/registry/core/accserver"

	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/errors"
	"gofreezer/pkg/runtime"
)

// REST implements the REST endpoint for usertoken
type REST struct {
	store       *etcdregistry.Store
	accRegistry account.Registry
}

// NewREST returns a RESTStorage object that will work with testtype.
func NewREST(opts generic.RESTOptions, accRegistry account.Registry) *REST {
	prefix := "/" + opts.ResourcePrefix
	newListFunc := func() runtime.Object { return &api.UserList{} }

	storageConfig := opts.StorageConfig
	storageInterface, dFunc := generic.NewRawStorage(storageConfig)

	store := &registry.Store{
		NewFunc:     func() runtime.Object { return &api.User{} },
		NewListFunc: newListFunc,
		KeyRootFunc: func(ctx freezerapi.Context) string {
			return prefix
		},
		KeyFunc: func(ctx freezerapi.Context, name string) (string, error) {
			return registry.NoNamespaceKeyFunc(ctx, prefix, name)
		},
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*api.User).Name, nil
		},
		PredicateFunc:           accserver.MatchAccServer,
		QualifiedResource:       api.Resource("accservers"),
		EnableGarbageCollection: opts.EnableGarbageCollection,
		DeleteCollectionWorkers: opts.DeleteCollectionWorkers,

		CreateStrategy:      accserver.Strategy,
		UpdateStrategy:      accserver.Strategy,
		DeleteStrategy:      accserver.Strategy,
		ReturnDeletedObject: true,
		// AfterCreate:         node.PadObj,
		// AfterUpdate:         node.PadObj,
		// AfterDelete:         node.PadObj,

		Storage:     storageInterface,
		DestroyFunc: dFunc,
	}

	return &REST{
		store:       etcdregistry.NewStore(*store),
		accRegistry: accRegistry,
	}
}

func (r *REST) New() runtime.Object {
	return &api.AccServer{}
}

func (*REST) NewList() runtime.Object {
	return &api.AccServerList{}
}

func (r *REST) Create(ctx freezerapi.Context, obj runtime.Object) (runtime.Object, error) {
	accsrv := obj.(*api.AccServer)
	accName := accsrv.Spec.AccName
	collectorHandler, err := r.accRegistry.GetCloudProvider(ctx, accName)
	if err == nil {
		err := collectorHandler.CreateServer(accsrv)
		if err != nil {
			return nil, errors.NewBadRequest(err.Error())
		}

		return r.store.Create(ctx, obj)
	}
	return nil, errors.NewBadRequest(fmt.Sprintf("not found cloud provider by account(%v)", accName))
}

func (r *REST) Delete(ctx freezerapi.Context, name string, options *freezerapi.DeleteOptions) (runtime.Object, error) {
	obj, err := r.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	accsrv := obj.(*api.AccServer)

	accName := accsrv.Spec.AccName
	collectorHandler, err := r.accRegistry.GetCloudProvider(ctx, accName)
	if err == nil {
		id, err := strconv.Atoi(accsrv.Spec.ID)
		if err != nil {
			return nil, errors.NewBadRequest(err.Error())
		}
		err = collectorHandler.DeleteServer(int64(id))
		if err != nil {
			return nil, errors.NewBadRequest(err.Error())
		}
		return r.store.Delete(ctx, name, options)
	}
	return nil, errors.NewBadRequest(fmt.Sprintf("not found cloud provider by account(%v)", accName))
}

func (r *REST) Get(ctx freezerapi.Context, name string) (runtime.Object, error) {
	obj, err := r.store.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	accsrv := obj.(*api.AccServer)

	err = r.getServer(ctx, accsrv)
	if err != nil {
		return nil, err
	}

	return accsrv, nil
}

func (r *REST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {
	obj, err := r.store.List(ctx, options)
	if err != nil {
		return nil, err
	}
	accsrvlist := obj.(*api.AccServerList)

	for i, v := range accsrvlist.Items {
		if err == nil {
			err = r.getServer(ctx, &v)
			if err != nil {
				return nil, err
			}
			accsrvlist.Items[i] = v
		} else {
			return nil, errors.NewBadRequest(fmt.Sprintf("not found cloud provider by account(%v)", v.Spec.AccName))
		}
	}

	return accsrvlist, nil
}

func (r *REST) getServer(ctx freezerapi.Context, accsrv *api.AccServer) error {
	accName := accsrv.Spec.AccName
	collectorHandler, err := r.accRegistry.GetCloudProvider(ctx, accName)

	if err == nil {
		id, err := strconv.ParseInt(accsrv.Spec.ID, 10, 63)
		if err != nil {
			return errors.NewInternalError(err)
		}

		subsrv, err := collectorHandler.GetServer(int(id))
		if err != nil {
			return errors.NewInternalError(err)
		}

		accsrv.Spec.DigitalOcean = subsrv.DigitalOcean
		accsrv.Spec.Vultr = subsrv.Vultr
		return nil
	}

	return errors.NewBadRequest(fmt.Sprintf("not found cloud provider by account(%v)", accName))
}
