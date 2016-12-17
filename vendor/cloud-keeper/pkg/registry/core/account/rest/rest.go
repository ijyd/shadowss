package rest

import (
	"cloud-keeper/pkg/registry/core/account/mysql"

	"cloud-keeper/pkg/api"

	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/runtime"
)

// REST implements the REST endpoint for usertoken
type REST struct {
	store *mysql.REST
}

func NewREST(mysql *mysql.REST) *REST {
	return &REST{
		store: mysql,
	}
}

func (r *REST) New() runtime.Object {
	return &api.Account{}
}

func (r *REST) NewList() runtime.Object {
	return &api.AccountList{}
}

func (r *REST) Get(ctx freezerapi.Context, name string) (runtime.Object, error) {
	return r.store.Get(ctx, name)
}

func (r *REST) Create(ctx freezerapi.Context, obj runtime.Object) (runtime.Object, error) {
	return r.store.Create(ctx, obj)
}

func (r *REST) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	return r.store.Update(ctx, name, objInfo)
}

func (r *REST) Delete(ctx freezerapi.Context, name string) (runtime.Object, error) {
	return r.store.Delete(ctx, name, nil)
}

func (r *REST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {
	return r.store.List(ctx, options)
}
