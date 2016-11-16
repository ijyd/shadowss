package apiserver

import (
	"cloud-keeper/pkg/api"
	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/rest"
)

// Registry is an interface for things that know how to store node.
type Registry interface {
	GetAPIServer(ctx freezerapi.Context, name string) (*api.APIServer, error)
	ListAPIServer(ctx freezerapi.Context, options *freezerapi.ListOptions) (*api.APIServerList, error)
	CreateAPIServer(ctx freezerapi.Context, user *api.APIServer) error
	DeleteAPIServer(ctx freezerapi.Context, name string) error
}

// storage puts strong typing around storage calls
type storage struct {
	rest.Getter
	rest.Lister
	rest.Creater
	rest.GracefulDeleter
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched
// types will panic.
func NewRegistry(g rest.Getter, l rest.Lister, c rest.Creater, d rest.GracefulDeleter) Registry {
	return &storage{
		Getter:          g,
		Lister:          l,
		Creater:         c,
		GracefulDeleter: d,
	}
}

func (s *storage) GetAPIServer(ctx freezerapi.Context, name string) (*api.APIServer, error) {
	obj, err := s.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return obj.(*api.APIServer), nil
}

func (s *storage) ListAPIServer(ctx freezerapi.Context, options *freezerapi.ListOptions) (*api.APIServerList, error) {
	obj, err := s.List(ctx, options)
	if err != nil {
		return nil, err
	}
	return obj.(*api.APIServerList), nil
}

func (s *storage) CreateAPIServer(ctx freezerapi.Context, obj *api.APIServer) error {
	_, err := s.Create(ctx, obj)
	return err
}

func (s *storage) DeleteAPIServer(ctx freezerapi.Context, name string) error {
	_, err := s.Delete(ctx, name, nil)
	return err
}
