package rest

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/node"
	"cloud-keeper/pkg/registry/core/userservice"
	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/runtime"
)

func NewExtendREST(node node.Registry, userservice userservice.Registry) (*BindUserREST, *APINodeREST) {
	return &BindUserREST{userservice}, &APINodeREST{node}
}

type BindUserREST struct {
	userservice userservice.Registry
}

func (*BindUserREST) New() runtime.Object {
	return &api.UserService{}
}

func (*BindUserREST) NewList() runtime.Object {

	return &api.UserServiceList{}
}

func (rs *BindUserREST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {
	return nil, nil
}

type APINodeREST struct {
	node node.Registry
}

func (*APINodeREST) New() runtime.Object {
	return &api.Node{}
}

func (*APINodeREST) NewList() runtime.Object {
	return &api.NodeList{}
}

//give a lables selector in request param like as label.selector
func (r *APINodeREST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {
	return r.List(ctx, options)
}
