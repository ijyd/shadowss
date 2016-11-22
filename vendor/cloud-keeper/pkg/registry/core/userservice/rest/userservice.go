package rest

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/node"
	"cloud-keeper/pkg/registry/core/nodeuser"
	"cloud-keeper/pkg/registry/core/user"
	"cloud-keeper/pkg/registry/core/userservice/etcd"
	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/errors"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/watch"

	"github.com/golang/glog"
)

type UserServiceREST struct {
	user     user.Registry
	node     node.Registry
	nodeuser nodeuser.Registry

	userservice *etcd.REST
}

func NewREST(userservice *etcd.REST, user user.Registry, node node.Registry, nodeuser nodeuser.Registry) *UserServiceREST {
	return &UserServiceREST{
		user:        user,
		userservice: userservice,
		node:        node,
		nodeuser:    nodeuser,
	}
}

func (*UserServiceREST) New() runtime.Object {
	return &api.UserService{}
}

func (*UserServiceREST) NewList() runtime.Object {
	return &api.UserServiceList{}
}

func (rs *UserServiceREST) Get(ctx freezerapi.Context, name string) (runtime.Object, error) {
	return rs.userservice.Get(ctx, name)
}

func (rs *UserServiceREST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {
	return rs.userservice.List(ctx, options)
}

func (r *UserServiceREST) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	return r.userservice.Update(ctx, name, objInfo)
}

func (r *UserServiceREST) Delete(ctx freezerapi.Context, name string, options *freezerapi.DeleteOptions) (runtime.Object, error) {
	return r.userservice.Delete(ctx, name, options)
}

func (r *UserServiceREST) Create(ctx freezerapi.Context, obj runtime.Object) (runtime.Object, error) {
	usersrv, ok := obj.(*api.UserService)
	if !ok {
		return nil, errors.NewBadRequest("not a User object")
	}
	nodelist, err := r.node.GetAPINodes(ctx, nil)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	defaultNode, ok := usersrv.Spec.Nodes[api.UserServicetDefaultNode]
	if !ok {
		return nil, errors.NewBadRequest("not found default node in nods")
	}
	delete(usersrv.Spec.Nodes, api.UserServicetDefaultNode)

	var i int
	for _, v := range nodelist.Items {

		nodeUser := &api.NodeUser{
			Spec: api.NodeUserSpec{
				User:     defaultNode.User,
				NodeName: v.Name,
				Phase:    api.NodeUserPhase(api.NodeUserPhaseAdd),
			},
		}
		nodeUser.Name = defaultNode.User.Name
		_, err := r.nodeuser.UpdateNodeUser(ctx, nodeUser)
		if err == nil {
			noderefer := api.NodeReferences{
				Host: v.Spec.Server.Host,
				User: defaultNode.User,
			}
			usersrv.Spec.Nodes[v.Name] = noderefer
		} else {
			glog.Warningf("create node user:%+v failure:%v\r\n", err)
		}

		//only give four api node
		if i > 4 {
			break
		}
		i++
	}

	return r.userservice.Create(ctx, usersrv)
}

func (r *UserServiceREST) Watch(ctx freezerapi.Context, options *freezerapi.ListOptions) (watch.Interface, error) {
	return nil, errors.NewMethodNotSupported(api.Resource("users"), "WATCH")
}
