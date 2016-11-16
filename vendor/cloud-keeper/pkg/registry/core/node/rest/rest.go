package rest

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/node/etcd"
	"cloud-keeper/pkg/registry/core/node/mysql"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/runtime"

	freezerapi "gofreezer/pkg/api"
)

type NodeREST struct {
	etcd  *etcd.REST
	mysql *mysql.REST
}

func NewREST(etcd *etcd.REST, mysql *mysql.REST) *NodeREST {
	return &NodeREST{
		etcd:  etcd,
		mysql: mysql,
	}
}

func (*NodeREST) New() runtime.Object {
	return &api.Node{}
}

func (*NodeREST) NewList() runtime.Object {
	return &api.NodeList{}
}

func (rs *NodeREST) Get(ctx freezerapi.Context, name string) (runtime.Object, error) {
	return rs.mysql.Get(ctx, name)
}

func (rs *NodeREST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {
	return rs.mysql.List(ctx, options)
}

func (r *NodeREST) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	return r.mysql.Update(ctx, name, objInfo)
}
