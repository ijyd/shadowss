package user

import (
	"cloud-keeper/pkg/api"
	"fmt"
	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/runtime"

	"github.com/golang/glog"
)

type NodeUserHandler interface {
	InitNodeUser(ctx freezerapi.Context, user *api.User) error
	AddNodeUser(ctx freezerapi.Context, updateUser *api.User, userservice *api.UserService, update bool, syncToNode bool) error
	DelNodeUser(ctx freezerapi.Context, user *api.User, nodeName string, update bool, syncToNode bool) error
	DelAllNodeUser(ctx freezerapi.Context, user *api.User, update bool, syncToNode bool) error
	DumpNodeUserToNode(ctx freezerapi.Context, nodename string) ([]*api.NodeUser, error)
	NewNodeUser(user *api.UserReferences, nodeName string) *api.NodeUser
	MigrateUserToDynamodb() error
}

type NoSQLUserHandler interface {
	FilterUserWithNodeName(ctx freezerapi.Context, nodename string, options *freezerapi.ListOptions) (*api.UserList, error)
}

// Registry is an interface for things that know how to store node.
type Registry interface {
	CreateUser(ctx freezerapi.Context, user *api.User) (runtime.Object, error)
	UpdateUser(ctx freezerapi.Context, svc *api.User) (*api.User, error)
	GetUser(ctx freezerapi.Context, name string) (*api.User, error)
	ListUsers(ctx freezerapi.Context, options *freezerapi.ListOptions) (*api.UserList, error)
	DeleteUser(ctx freezerapi.Context, name string) (runtime.Object, error)

	ListUserByNodeName(ctx freezerapi.Context, nodename string, options *freezerapi.ListOptions) (*api.UserList, error)
	UpdateUserByNodeUser(ctx freezerapi.Context, svc *api.NodeUser) error
	ReInitNodeToUser(ctx freezerapi.Context, user *api.User) error
	AddNodeToUser(ctx freezerapi.Context, updateUser *api.User, userservice *api.UserService, update bool, syncToNode bool) error
	DelNodeFromUser(ctx freezerapi.Context, user *api.User, nodeName string, update bool, syncToNode bool) error
	DelAllNodeFromUser(ctx freezerapi.Context, user *api.User, update bool, syncToNode bool) error
	DumpNodeUser(ctx freezerapi.Context, nodename string) ([]*api.NodeUser, error)
	CreateNodeUser(user *api.UserReferences, nodeName string) *api.NodeUser
	MigrateUser() error

	//resume user,clear traffic and enable user
	ResumeUsers(ctx freezerapi.Context, name []string) error
}

// storage puts strong typing around storage calls
type storage struct {
	rest.Getter
	rest.GracefulDeleter
	rest.Lister
	rest.Updater
	rest.Creater
	NodeUserHandler
	NoSQLUserHandler
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched
// types will panic.
func NewRegistry(g rest.Getter, d rest.GracefulDeleter, l rest.Lister, u rest.Updater, c rest.Creater, nodeUser NodeUserHandler, nosql NoSQLUserHandler) Registry {
	return &storage{
		Getter:           g,
		GracefulDeleter:  d,
		Lister:           l,
		Updater:          u,
		Creater:          c,
		NodeUserHandler:  nodeUser,
		NoSQLUserHandler: nosql,
	}
}

func (s *storage) GetUser(ctx freezerapi.Context, name string) (*api.User, error) {
	obj, err := s.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return obj.(*api.User), nil
}

func (s *storage) ListUsers(ctx freezerapi.Context, options *freezerapi.ListOptions) (*api.UserList, error) {
	obj, err := s.List(ctx, options)
	if err != nil {
		return nil, err
	}
	return obj.(*api.UserList), nil
}

func (s *storage) UpdateUser(ctx freezerapi.Context, user *api.User) (*api.User, error) {
	obj, _, err := s.Update(ctx, user.Name, rest.DefaultUpdatedObjectInfo(user, api.Scheme))
	if err != nil {
		return nil, err
	}
	newuser := obj.(*api.User)

	return newuser, nil
}

func (s *storage) UpdateUserByNodeUser(ctx freezerapi.Context, nodeuser *api.NodeUser) error {
	nodename := nodeuser.Spec.NodeName
	username := nodeuser.Spec.User.Name
	objuser, err := s.GetUser(ctx, nodeuser.Spec.User.Name)
	if err != nil {
		return err
	}
	objuser.Spec.DetailInfo.DownloadTraffic += nodeuser.Spec.User.DownloadTraffic
	objuser.Spec.DetailInfo.UploadTraffic += nodeuser.Spec.User.UploadTraffic

	noderefer, ok := objuser.Spec.UserService.Nodes[nodename]
	if !ok {
		return fmt.Errorf("not found user(%v) in node(%v), may be need to disable user", username, nodename)
	}

	noderefer.User.Port = nodeuser.Spec.User.Port
	noderefer.User.DownloadTraffic = nodeuser.Spec.User.DownloadTraffic
	noderefer.User.UploadTraffic = nodeuser.Spec.User.UploadTraffic
	objuser.Spec.UserService.Nodes[nodename] = noderefer

	if lastActive, ok := nodeuser.Annotations[api.UserFakeAnnotationLastActiveTime]; ok {
		if objuser.Annotations == nil {
			objuser.Annotations = make(map[string]string)
		}
		objuser.Annotations[api.UserFakeAnnotationLastActiveTime] = lastActive
	}

	glog.V(5).Infof("UpdateUserByNodeUser user %+v \r\nupdate by node user %+v\r\n", *objuser, *nodeuser)
	_, err = s.UpdateUser(ctx, objuser)
	return err
}

func (s *storage) CreateUser(ctx freezerapi.Context, user *api.User) (runtime.Object, error) {
	obj, err := s.Create(ctx, user)
	return obj, err
}

func (s *storage) DeleteUser(ctx freezerapi.Context, name string) (runtime.Object, error) {
	obj, err := s.Delete(ctx, name, nil)
	return obj, err
}

func (s *storage) ListUserByNodeName(ctx freezerapi.Context, nodename string, options *freezerapi.ListOptions) (*api.UserList, error) {
	userlist, err := s.FilterUserWithNodeName(ctx, nodename, options)
	if err != nil {
		return nil, err
	}

	return userlist, nil
}

func (s *storage) ReInitNodeToUser(ctx freezerapi.Context, user *api.User) error {
	return s.InitNodeUser(ctx, user)
}

func (s *storage) AddNodeToUser(ctx freezerapi.Context, updateUser *api.User, usersrv *api.UserService, update bool, syncToNode bool) error {
	return s.AddNodeUser(ctx, updateUser, usersrv, update, syncToNode)
}

func (s *storage) DelNodeFromUser(ctx freezerapi.Context, user *api.User, nodeName string, update bool, syncToNode bool) error {
	return s.DelNodeUser(ctx, user, nodeName, update, syncToNode)
}

func (s *storage) DelAllNodeFromUser(ctx freezerapi.Context, user *api.User, update bool, syncToNode bool) error {
	return s.DelAllNodeUser(ctx, user, update, syncToNode)
}

func (s *storage) DumpNodeUser(ctx freezerapi.Context, nodename string) ([]*api.NodeUser, error) {
	return s.DumpNodeUserToNode(ctx, nodename)
}

func (s *storage) CreateNodeUser(user *api.UserReferences, nodeName string) *api.NodeUser {
	return s.NewNodeUser(user, nodeName)
}

func (s *storage) MigrateUser() error {
	return s.MigrateUserToDynamodb()
}

type updateTrans func(name string) (*api.User, error)

//resume user,clear traffic and enable user,add user to node
func (s *storage) getAndUpdate(ctx freezerapi.Context, name string, trans updateTrans) (*api.User, error) {
	return trans(name)
}

//resume user,clear traffic and enable user,add user to node
func (s *storage) ResumeUsers(ctx freezerapi.Context, names []string) error {

	trans := func(name string) (*api.User, error) {
		user, err := s.GetUser(ctx, name)
		if err != nil {
			return nil, err
		}
		user.Spec.DetailInfo.DownloadTraffic = 0
		user.Spec.DetailInfo.UploadTraffic = 0
		user.Spec.DetailInfo.Status = true
		newUser, err := s.UpdateUser(ctx, user)
		return newUser, err
	}

	for _, v := range names {
		user, err := s.getAndUpdate(ctx, v, trans)
		if err != nil {
			glog.Warningf("resume user(%s) failure : %v\r\n", v, err)
		} else {
			err = s.ReInitNodeToUser(ctx, user)
			if err != nil {
				glog.Warningf("reinit resume user(%v) failure:%v\r\n", v, err)
			}
		}
	}

	return nil
}
