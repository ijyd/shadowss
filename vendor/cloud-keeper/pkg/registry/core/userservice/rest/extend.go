package rest

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/userservice"

	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/runtime"
)

func NewExtendREST(userservice userservice.Registry) (*BindingNodeREST, *PropertiesREST) {
	return &BindingNodeREST{userservice}, &PropertiesREST{userservice}
}

type BindingNodeREST struct {
	userservice userservice.Registry
}

func (*BindingNodeREST) New() runtime.Object {
	return &api.UserService{}
}

func (*BindingNodeREST) NewList() runtime.Object {
	return &api.UserServiceList{}
}

func (rs *BindingNodeREST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {
	return nil, nil
}

func (r *BindingNodeREST) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	return nil, false, nil
}

type PropertiesREST struct {
	userservice userservice.Registry
}

func (*PropertiesREST) New() runtime.Object {
	return &api.UserService{}
}

func (r *PropertiesREST) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	//return r.mysqlregistry.Update(ctx, name, objInfo)
	return nil, false, nil
}
