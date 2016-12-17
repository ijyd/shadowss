package rest

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/node"
	"cloud-keeper/pkg/registry/core/nodeuser"
	"cloud-keeper/pkg/registry/core/user/dynamodb"
	"cloud-keeper/pkg/registry/core/user/etcd"
	"cloud-keeper/pkg/registry/core/user/mysql"
	"fmt"
	"gofreezer/pkg/api/errors"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/fields"
	"gofreezer/pkg/runtime"

	"github.com/golang/glog"

	freezerapi "gofreezer/pkg/api"
)

type UserREST struct {
	etcd   *etcd.REST
	mysql  *mysql.REST
	dynamo *dynamodb.REST

	node     node.Registry
	nodeuser nodeuser.Registry
	//userservice userservice.Registry
	//user user.Registry
}

func NewREST(etcdHandler *etcd.REST, mysqlHandler *mysql.REST, dynamo *dynamodb.REST, node node.Registry, nodeuser nodeuser.Registry) *UserREST {
	return &UserREST{
		etcd:     etcdHandler,
		mysql:    mysqlHandler,
		dynamo:   dynamo,
		node:     node,
		nodeuser: nodeuser,
	}
}

func (*UserREST) New() runtime.Object {
	return &api.User{}
}

func (*UserREST) NewList() runtime.Object {
	return &api.UserList{}
}

func (r *UserREST) mergeUser(left *api.User, out *api.User) {
	out.Annotations = make(map[string]string)
	for k, v := range left.Annotations {
		out.Annotations[k] = v
	}

	out.Spec.UserService.NodeCnt = left.Spec.UserService.NodeCnt
	out.Spec.UserService.Status = left.Spec.UserService.Status

	out.Spec.UserService.Nodes = make(map[string]api.NodeReferences)
	for nodeName, nodeRefer := range left.Spec.UserService.Nodes {
		refer := api.NodeReferences{
			Host: nodeRefer.Host,
		}
		refer.User.DownloadTraffic = nodeRefer.User.DownloadTraffic
		refer.User.UploadTraffic = nodeRefer.User.UploadTraffic
		refer.User.Name = nodeRefer.User.Name
		refer.User.EnableOTA = nodeRefer.User.EnableOTA
		refer.User.ID = nodeRefer.User.ID
		refer.User.Method = nodeRefer.User.Method
		refer.User.Port = nodeRefer.User.Port
		refer.User.Password = nodeRefer.User.Password

		out.Spec.UserService.Nodes[nodeName] = refer
	}

}

func (r *UserREST) Get(ctx freezerapi.Context, name string) (runtime.Object, error) {

	obj, err := r.mysql.Get(ctx, name)
	if err != nil {
		glog.Infof("Get resource failure %v\r\n", err)
		return nil, err
	}
	user := obj.(*api.User)
	user.Name = user.Spec.DetailInfo.Name

	etcdobj, err := r.dynamo.Get(ctx, name)
	if err == nil {
		left := etcdobj.(*api.User)
		// user.Annotations = make(map[string]string)
		// for k, v := range etcdUser.Annotations {
		// 	user.Annotations[k] = v
		// }
		r.mergeUser(left, user)
	}

	return user, nil
}

func (r *UserREST) FilterUserWithNodeName(ctx freezerapi.Context, nodename string, options *freezerapi.ListOptions) (*api.UserList, error) {
	options = &freezerapi.ListOptions{}

	selectorStr := fmt.Sprintf("spec.userService.nodes.%s", nodename)
	options.FieldSelector = fields.ParseSelectorOrDie(selectorStr)
	obj, err := r.dynamo.List(ctx, options)
	if err != nil {
		return nil, err
	}
	userlist := obj.(*api.UserList)

	var newItems []api.User
	for _, user := range userlist.Items {
		_, ok := user.Spec.UserService.Nodes[nodename]
		if ok {
			newItems = append(newItems, user)
		}
	}

	userlist.Items = newItems

	return userlist, nil
}

func (r *UserREST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {

	if (options == nil) || (options != nil && options.PageSelector == nil) ||
		(options.PageSelector != nil && options.PageSelector.Empty()) {
		return nil, errors.NewInternalError(fmt.Errorf("you must give a pagination selector.\r\n"))
	}

	_, perPage := options.PageSelector.RequirePage()
	if perPage > 20 {
		return nil, errors.NewInternalError(fmt.Errorf("you must give a (perPage < 15) for pagination selector.\r\n"))
	}

	mysqlobj, err := r.mysql.List(ctx, options)
	if err != nil {
		return nil, err
	}
	userlist := mysqlobj.(*api.UserList)

	// obj, err := r.dynamo.List(ctx, nil)
	// if err != nil {
	// 	return nil, err
	// }
	// dynamoUserList := obj.(*api.UserList)
	// dynamodbUserListMap := make(map[string]*api.User)
	// for k, v := range dynamoUserList.Items {
	// 	glog.V(5).Infof("found user(%v) in dynamodb\r\n", v.Name)
	// 	dynamodbUserListMap[v.Name] = &dynamoUserList.Items[k]
	// }

	for k, v := range userlist.Items {
		userlist.Items[k].Name = v.Spec.DetailInfo.Name
		name := v.Spec.DetailInfo.Name

		dynamoObj, err := r.dynamo.Get(ctx, name)
		if err != nil {
			glog.Warningf("not found user(%v) in dynamodb\r\n", name)
			continue
		}
		dynamoUser, ok := dynamoObj.(*api.User)

		//dynamoUser, ok := dynamodbUserListMap[name]
		if ok {
			r.mergeUser(dynamoUser, &userlist.Items[k])
		} else {
			glog.Warningf("not found user(%v) in dynamodb\r\n", name)
		}
	}

	return userlist, nil
}

func (r *UserREST) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {

	newObj, err := objInfo.UpdatedObject(ctx, nil)
	if err != nil {
		glog.Errorf("extact user failure:%v\r\n", err)
		return nil, false, err
	}
	newUser := newObj.(*api.User)

	if !r.CheckUser(newUser) {
		//need to disable user and delete user from node
		newUser.Spec.DetailInfo.Status = false
		newUser.Spec.UserService.Status = false
		if len(newUser.Spec.UserService.Nodes) > 0 {
			r.DelAllNodeUser(ctx, newUser, false, true)
		}
	}

	_, flag, err := r.dynamo.Update(ctx, name, rest.DefaultUpdatedObjectInfo(newUser, api.Scheme))
	if err != nil {
		glog.Errorf("update user(%+v) error:%v\r\n", name, err)
		return nil, false, err
	}

	user, flag, err := r.mysql.Update(ctx, name, rest.DefaultUpdatedObjectInfo(newUser, api.Scheme))
	if err != nil {
		glog.Errorf("update user(%+v) error:%v\r\n", name, err)
		return nil, false, err
	}

	glog.V(5).Infof("update user(%+v) done(%v)\r\n", name, err)

	return user, flag, nil
}

func (r *UserREST) Delete(ctx freezerapi.Context, name string, options *freezerapi.DeleteOptions) (runtime.Object, error) {
	obj, err := r.mysql.Get(ctx, name)
	if err != nil {
		return nil, errors.NewNotFound(api.Resource("users"), err.Error())
	}
	user := obj.(*api.User)

	user.Spec.DetailInfo.Status = false
	_, _, err = r.dynamo.Update(ctx, user.Name, rest.DefaultUpdatedObjectInfo(user, api.Scheme))
	if err != nil {
		return nil, err
	}

	// user := obj.(*api.User)
	// user.Spec.DetailInfo.Status = false

	//disable user
	upobj, _, err := r.mysql.Update(ctx, user.Name, rest.DefaultUpdatedObjectInfo(user, api.Scheme))
	return upobj, err
}

func (r *UserREST) Create(ctx freezerapi.Context, obj runtime.Object) (runtime.Object, error) {
	user, ok := obj.(*api.User)
	if !ok {
		return nil, errors.NewBadRequest("not a User object")
	}

	err := r.InitNodeUser(ctx, user)
	if err != nil {
		return nil, err
	}

	user.Spec.DetailInfo.Activation = true
	user.Spec.DetailInfo.Delete = false
	user.Spec.DetailInfo.IsAdmin = false
	user.Spec.DetailInfo.Status = true

	switch user.Spec.DetailInfo.PackageType {
	case api.PackageTypeDefault:
	default:
		user.Spec.DetailInfo.PackageType = api.PackageTypeDefault
	}

	switch user.Spec.DetailInfo.UserType {
	case api.UserTypeDesktopRouter:
	default:
		user.Spec.DetailInfo.UserType = api.UserTypeDesktopRouter
	}

	switch user.Spec.DetailInfo.Bandwidth {
	case api.UserBandwidthUnlimited:
	default:
		user.Spec.DetailInfo.Bandwidth = api.UserBandwidthUnlimited
	}

	_, err = r.mysql.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	resObj, err := r.dynamo.Create(ctx, user)
	if err != nil {
		glog.Errorf("create user error:%v will delete from mysql\r\n", err)
		r.mysql.Delete(ctx, user.Name, nil)
		return nil, err
	}

	return resObj, nil
}
