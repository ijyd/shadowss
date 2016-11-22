package rest

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/node/etcd"
	"cloud-keeper/pkg/registry/core/node/mysql"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/runtime"

	"github.com/golang/glog"

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
	mysqlobj, err := rs.mysql.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	node := mysqlobj.(*api.Node)

	node.Spec.Users = make(map[string]api.NodeUser)
	etcdObj, err := rs.etcd.Get(ctx, name)
	if err == nil {
		etcdNode := etcdObj.(*api.Node)
		node.ObjectMeta = etcdNode.ObjectMeta
		for k, v := range etcdNode.Spec.Users {
			node.Spec.Users[k] = v
		}
	} else {
		node.Name = node.Spec.Server.Name
	}

	return node, nil
}

func (rs *NodeREST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {
	mysqlobj, err := rs.mysql.List(ctx, options)
	if err != nil {
		return nil, err
	}

	nodeList := mysqlobj.(*api.NodeList)
	for i, v := range nodeList.Items {
		etcdObj, err := rs.etcd.Get(ctx, v.Spec.Server.Name)
		if err == nil {
			etcdNode := etcdObj.(*api.Node)
			nodeList.Items[i].ObjectMeta = etcdNode.ObjectMeta
			nodeList.Items[i].Spec.Users = make(map[string]api.NodeUser)
			for k, v := range etcdNode.Spec.Users {
				nodeList.Items[i].Spec.Users[k] = v
			}
		} else {
			nodeList.Items[i].Name = nodeList.Items[i].Spec.Server.Name
		}

	}

	return nodeList, nil
}

func (r *NodeREST) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	glog.V(5).Infof("update etcd name:%v\r\n", name)
	obj, upflag, err := r.etcd.Update(ctx, name, objInfo)
	if err != nil {
		glog.V(5).Infof("update etcd err:%v\r\n", err)
		return obj, upflag, err
	}
	glog.V(5).Infof("update mysql name:%v\r\n", name)

	return r.mysql.Update(ctx, name, objInfo)
}
