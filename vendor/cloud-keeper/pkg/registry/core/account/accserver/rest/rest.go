package rest

import (
	"strconv"

	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/errors"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/runtime"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/account"
)

// REST implements the REST endpoint for usertoken
type REST struct {
	accRegistry account.Registry
}

// NewREST returns a RESTStorage object that will work with testtype.
func NewREST(accRegistry account.Registry) *REST {
	return &REST{accRegistry}
}

func (r *REST) New() runtime.Object {
	return &api.AccServer{}
}

func (*REST) NewList() runtime.Object {
	return &api.AccServerList{}
}

// ConnectMethods returns the list of HTTP methods that can be proxied
// func (r *REST) ConnectMethods() []string {
// 	return []string{"GET"}
// }
//
// // NewConnectOptions returns versioned resource that represents proxy parameters
// func (r *REST) NewConnectOptions() (runtime.Object, bool, string) {
// 	return nil, false, ""
// }
//
// // Connect returns a handler for the pod proxy
// func (r *REST) Connect(ctx freezerapi.Context, id string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
// 	return NewConnector(responder), nil
// }

func (r *REST) newCloudServer(ctx freezerapi.Context, accsrv *api.AccServer, accName string) error {
	collectorHandler, err := r.accRegistry.GetCloudProvider(ctx, accName)
	if err != nil {
		return errors.NewInternalError(err)
	}
	err = collectorHandler.CreateServer(accsrv)
	if err != nil {
		return errors.NewInternalError(err)
	}
	return nil
}

func (r *REST) deleteCloudServer(ctx freezerapi.Context, accsrv *api.AccServer, accName string) error {

	collectorHandler, err := r.accRegistry.GetCloudProvider(ctx, accName)
	if err != nil {
		return errors.NewInternalError(err)
	}

	id, err := strconv.Atoi(accsrv.Spec.ID)
	if err != nil {
		return errors.NewInternalError(err)
	}
	err = collectorHandler.DeleteServer(int64(id))
	if err != nil {
		return errors.NewInternalError(err)
	}

	return nil

}

func (r *REST) getServer(ctx freezerapi.Context, accName string) (*api.AccServer, error) {
	// collectorHandler, err := r.accRegistry.GetCloudProvider(ctx, accName)
	//
	// if err != nil {
	// 	return nil, errors.NewInternalError(err)
	// }
	//
	// id, err := strconv.ParseInt(accsrv.Spec.ID, 10, 63)
	// if err != nil {
	// 	return nil, errors.NewInternalError(err)
	// }
	//
	// server, err := collectorHandler.GetServer(int(id))
	// if err != nil {
	// 	return nil, errors.NewInternalError(err)
	// }
	//
	// return server, nil
	return nil, nil
}

//bool is not found will be true
func (r *REST) getServers(ctx freezerapi.Context, accName string, options *freezerapi.ListOptions) ([]api.AccServer, error) {
	collectorHandler, err := r.accRegistry.GetCloudProvider(ctx, accName)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	srvSpec, err := collectorHandler.GetServers(options.PageSelector)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	return srvSpec, nil
}

// func (r *REST) Get(ctx freezerapi.Context, name string) (runtime.Object, error) {
// 	return r.store.Get(ctx, name)
// }

func (r *REST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {
	accsrvlist := &api.AccServerList{}

	if name, ok := options.FieldSelector.RequiresExactMatch("metadata.name"); ok {
		srvs, err := r.getServers(ctx, name, options)
		if err != nil {
			return nil, err
		}

		accsrvlist.Items = append(accsrvlist.Items, srvs...)
		return accsrvlist, nil

	} else {
		return nil, errors.NewBadRequest("need a 'metadata.name' filed selector")
	}
}

func (r *REST) Create(ctx freezerapi.Context, obj runtime.Object) (runtime.Object, error) {
	accsrv := obj.(*api.AccServer)
	err := r.newCloudServer(ctx, accsrv, accsrv.Spec.AccName)

	return accsrv, err
}

func (r *REST) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {

	obj, err := objInfo.UpdatedObject(ctx, nil)
	if err != nil {
		return nil, false, err
	}

	accsrv := obj.(*api.AccServer)
	if len(accsrv.Spec.ID) == 0 {
		err = r.newCloudServer(ctx, accsrv, name)
	} else {
		err = r.deleteCloudServer(ctx, accsrv, name)
	}

	return obj, true, err
}
