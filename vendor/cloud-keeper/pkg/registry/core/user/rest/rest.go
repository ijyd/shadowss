package rest

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/user"
	"cloud-keeper/pkg/registry/core/userservice"
	apierr "gofreezer/pkg/api/errors"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/runtime"

	freezerapi "gofreezer/pkg/api"
)

type UserREST struct {
	userservice userservice.Registry
	user        user.Registry
}

func NewREST(user user.Registry, userservice userservice.Registry) *UserREST {
	return &UserREST{
		userservice: userservice,
		user:        user,
	}
}

func (*UserREST) New() runtime.Object {
	return &api.User{}
}

func (*UserREST) NewList() runtime.Object {
	return &api.UserList{}
}

func (rs *UserREST) Get(ctx freezerapi.Context, name string) (runtime.Object, error) {
	return rs.user.GetUser(ctx, name)
}

func (rs *UserREST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {
	return rs.user.ListUsers(ctx, options)
}

func (r *UserREST) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	return r.user.Update(ctx, name, objInfo)
}

func (r *UserREST) Delete(ctx freezerapi.Context, name string, options *freezerapi.DeleteOptions) (runtime.Object, error) {
	return r.user.DeleteUser(ctx, name)
}

func (r *UserREST) Create(ctx freezerapi.Context, obj runtime.Object) (runtime.Object, error) {
	user, ok := obj.(*api.User)
	if !ok {
		return nil, apierr.NewBadRequest("not a User object")
	}

	//for user default info
	userRefer := api.UserReferences{
		ID:        user.Spec.DetailInfo.ID,
		Name:      user.Name,
		Port:      0,
		Method:    string("aes-256-cfb"),
		Password:  user.Spec.DetailInfo.Passwd,
		EnableOTA: true,
	}
	spec := api.UserServiceSpec{
		Nodes: map[string]api.NodeReferences{
			api.UserServicetDefaultNode: api.NodeReferences{
				User: userRefer,
			},
		},
		NodeCnt: 0,
	}

	userSrv := &api.UserService{}
	userSrv.Spec = spec
	userSrv.Name = user.Name

	_, err := r.userservice.CreateUserService(ctx, userSrv)
	if err != nil {
		return nil, apierr.NewInternalError(err)
	}

	return r.user.CreateUser(ctx, user)
}
