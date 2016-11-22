package userservice

import (
	"cloud-keeper/pkg/api"
	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/rest"
)

// Registry is an interface for things that know how to store node.
type Registry interface {
	UpdateUserServiceByNodeUser(ctx freezerapi.Context, user *api.UserService) error
	CreateUserService(ctx freezerapi.Context, user *api.UserService) (*api.UserService, error)
	GetUserService(ctx freezerapi.Context, name string) (*api.UserService, error)
	ListUserServices(ctx freezerapi.Context, options *freezerapi.ListOptions) (*api.UserServiceList, error)
}

// storage puts strong typing around storage calls
type storage struct {
	rest.Getter
	rest.Lister
	rest.Creater
	rest.Updater
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched
// types will panic.
func NewRegistry(s rest.Getter, l rest.Lister, c rest.Creater, u rest.Updater) Registry {
	return &storage{
		Getter:  s,
		Lister:  l,
		Creater: c,
		Updater: u,
	}
}

func (s *storage) GetUserService(ctx freezerapi.Context, name string) (*api.UserService, error) {
	obj, err := s.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return obj.(*api.UserService), nil
}

func (s *storage) ListUserServices(ctx freezerapi.Context, options *freezerapi.ListOptions) (*api.UserServiceList, error) {
	obj, err := s.List(ctx, options)
	if err != nil {
		return nil, err
	}
	return obj.(*api.UserServiceList), nil
}

func (s *storage) CreateUserService(ctx freezerapi.Context, user *api.UserService) (*api.UserService, error) {
	obj, err := s.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return obj.(*api.UserService), nil
}

func (s *storage) UpdateUserServiceByNodeUser(ctx freezerapi.Context, user *api.UserService) error {
	_, _, err := s.Update(ctx, user.Name, rest.DefaultUpdatedObjectInfo(user, api.Scheme))
	return err
}
