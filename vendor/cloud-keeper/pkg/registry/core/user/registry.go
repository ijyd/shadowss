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
	AddNodeUser(ctx freezerapi.Context, updateUser *api.User, userservice *api.UserService) error
	DumpNodeUserToNode(ctx freezerapi.Context, nodename string) error
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
	ReInitNodeUser(ctx freezerapi.Context, user *api.User) error
	AddNodeUserByUserService(ctx freezerapi.Context, updateUser *api.User, userservice *api.UserService) error
	DumpNodeUser(ctx freezerapi.Context, nodename string) error
	CreateNodeUser(user *api.UserReferences, nodeName string) *api.NodeUser
	MigrateUser() error
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

	glog.V(5).Infof("UpdateUserByNodeUser *********begin \r\n")
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
	glog.V(5).Infof("user %+v update by node user %+v\r\n", objuser, nodeuser)
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

func (s *storage) ReInitNodeUser(ctx freezerapi.Context, user *api.User) error {
	return s.InitNodeUser(ctx, user)
}

func (s *storage) AddNodeUserByUserService(ctx freezerapi.Context, updateUser *api.User, usersrv *api.UserService) error {
	return s.AddNodeUser(ctx, updateUser, usersrv)
}

func (s *storage) DumpNodeUser(ctx freezerapi.Context, nodename string) error {
	return s.DumpNodeUserToNode(ctx, nodename)
}

func (s *storage) CreateNodeUser(user *api.UserReferences, nodeName string) *api.NodeUser {
	return s.NewNodeUser(user, nodeName)
}

func (s *storage) MigrateUser() error {
	return s.MigrateUserToDynamodb()
}
