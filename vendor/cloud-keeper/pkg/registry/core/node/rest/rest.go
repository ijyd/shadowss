package rest

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/node/etcd"
	"cloud-keeper/pkg/registry/core/node/mysql"
	"cloud-keeper/pkg/registry/core/nodeuser"
	"cloud-keeper/pkg/registry/core/user"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/runtime"
	"strings"

	"github.com/golang/glog"

	freezerapi "gofreezer/pkg/api"
)

type NodeREST struct {
	etcd  *etcd.REST
	mysql *mysql.REST

	user user.Registry
	//nodeUser nodeuser.Registry
}

func NewREST(etcd *etcd.REST, mysql *mysql.REST) *NodeREST {
	return &NodeREST{
		etcd:  etcd,
		mysql: mysql,
	}
}

func (r *NodeREST) SetRequireRegistry(user user.Registry, nodeUser nodeuser.Registry) {
	r.user = user
	//r.nodeUser = nodeUser
}

func (*NodeREST) New() runtime.Object {
	return &api.Node{}
}

func (*NodeREST) NewList() runtime.Object {
	return &api.NodeList{}
}

func (r *NodeREST) mergeNoSqlNode(left *api.Node, out *api.Node) {
	out.UID = left.UID
	out.CreationTimestamp = left.CreationTimestamp
	out.SelfLink = left.SelfLink
	out.Name = left.Name
	out.ResourceVersion = left.ResourceVersion

	out.Labels = make(map[string]string)
	for lk, lv := range left.Labels {
		out.Labels[lk] = lv
	}

	out.Annotations = make(map[string]string)
	for ak, av := range left.Annotations {
		out.Annotations[ak] = av
	}

	out.Spec.Users = make(map[string]api.NodeUserSpec)
	for k, v := range left.Spec.Users {
		out.Spec.Users[k] = v
	}
}

func (r *NodeREST) mergeSqlNode(left *api.Node, out *api.Node) {
	out.Spec.Server = left.Spec.Server
}

func (r *NodeREST) Get(ctx freezerapi.Context, name string) (runtime.Object, error) {
	mysqlobj, err := r.mysql.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	node := mysqlobj.(*api.Node)

	node.Spec.Users = make(map[string]api.NodeUserSpec)
	etcdObj, err := r.etcd.Get(ctx, name)
	if err == nil {
		etcdNode := etcdObj.(*api.Node)
		r.mergeNoSqlNode(etcdNode, node)
	} else {
		node.Name = node.Spec.Server.Name
	}

	return node, nil
}

//need merge mysql node info into etcd. because of we use the lables options for ListOptions
func (r *NodeREST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {

	mysqlobj, err := r.mysql.List(ctx, options)
	if err != nil {
		return nil, err
	}

	nodeList := mysqlobj.(*api.NodeList)
	mysqlNodeMap := make(map[string]*api.Node)
	for k, v := range nodeList.Items {
		mysqlNodeMap[v.Spec.Server.Name] = &nodeList.Items[k]
		nodeList.Items[k].Name = v.Spec.Server.Name
	}

	etcdObj, err := r.etcd.List(ctx, options)
	if err != nil {
		return nil, err
	}
	etcdNodeList := etcdObj.(*api.NodeList)
	etcdNodeMap := make(map[string]*api.Node)
	for k, v := range etcdNodeList.Items {
		etcdNodeMap[v.Name] = &etcdNodeList.Items[k]
	}

	//merge NodeUser from etcd
	for k, v := range nodeList.Items {
		etcdNode, ok := etcdNodeMap[v.Name]
		if ok {
			// nodeList.Items[k].Spec.Users = make(map[string]api.NodeUserSpec)
			// for userName, noedUser := range etcdNode.Spec.Users {
			// 	nodeList.Items[k].Spec.Users[userName] = noedUser
			// }
			// nodeList.Items[k].ResourceVersion = etcdNode.ResourceVersion
			r.mergeNoSqlNode(etcdNode, &nodeList.Items[k])
		}
	}

	var newNodes []api.Node
	if options.LabelSelector != nil && !options.LabelSelector.Empty() {
		for i, v := range etcdNodeList.Items {
			name := v.Name
			mysqlNode, ok := mysqlNodeMap[name]
			if ok {
				r.mergeSqlNode(mysqlNode, &etcdNodeList.Items[i])
				newNodes = append(newNodes, etcdNodeList.Items[i])
			} else {
				//if not found in mysql delete this item
				glog.Warningf("not found node(%v) in etcd\r\n", name)
			}
		}
		etcdNodeList.Items = etcdNodeList.Items[:0]
		etcdNodeList.Items = append(etcdNodeList.Items, newNodes...)
		return etcdNodeList, nil
	}

	return nodeList, nil
}

func (r *NodeREST) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	obj, upflag, err := r.etcd.Update(ctx, name, objInfo)
	if err != nil {
		glog.V(5).Infof("update etcd err:%v\r\n", err)
		return nil, upflag, err
	}

	node := obj.(*api.Node)

	if strings.Compare(node.Annotations[api.NodeAnnotationRefreshCnt], "0") == 0 {
		glog.Infof("refresh new node(%s) need to sync node user", node.Name)
		//ListUserServicesByNodeName
		go r.user.DumpNodeUser(ctx, name)
	}

	return r.mysql.Update(ctx, name, objInfo)
}
