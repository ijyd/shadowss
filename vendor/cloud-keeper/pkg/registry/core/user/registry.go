package user

import (
	"cloud-keeper/pkg/api"
	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/runtime"
)

// Registry is an interface for things that know how to store node.
type Registry interface {
	CreateUser(ctx freezerapi.Context, user *api.User) (runtime.Object, error)
	UpdateUser(ctx freezerapi.Context, svc *api.User) error
	Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (*api.User, bool, error)
	GetUser(ctx freezerapi.Context, name string) (*api.User, error)
	ListUsers(ctx freezerapi.Context, options *freezerapi.ListOptions) (*api.UserList, error)
	DeleteUser(ctx freezerapi.Context, name string) (runtime.Object, error)
}

// storage puts strong typing around storage calls
type storage struct {
	rest.Getter
	rest.GracefulDeleter
	rest.Lister
	rest.Updater
	rest.Creater
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched
// types will panic.
func NewRegistry(g rest.Getter, d rest.GracefulDeleter, l rest.Lister, u rest.Updater, c rest.Creater) Registry {
	return &storage{
		Getter:          g,
		GracefulDeleter: d,
		Lister:          l,
		Updater:         u,
		Creater:         c,
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

func (s *storage) UpdateUser(ctx freezerapi.Context, user *api.User) error {
	_, _, err := s.Update(ctx, user.Name, rest.DefaultUpdatedObjectInfo(user, api.Scheme))
	return err
}

func (s *storage) Update(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (*api.User, bool, error) {
	return s.Update(ctx, name, objInfo)
}

func (s *storage) CreateUser(ctx freezerapi.Context, user *api.User) (runtime.Object, error) {
	obj, err := s.Create(ctx, user)
	return obj, err
}

func (s *storage) DeleteUser(ctx freezerapi.Context, name string) (runtime.Object, error) {
	obj, err := s.Delete(ctx, name, nil)
	return obj, err
}
