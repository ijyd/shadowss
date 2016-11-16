package nodeuser

import (
	"cloud-keeper/pkg/api"
	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/rest"
)

// Registry is an interface for things that know how to store node.
type Registry interface {
	CreateNodeUser(ctx freezerapi.Context, user *api.NodeUser) (*api.NodeUser, error)
	DeleteNodeUser(ctx freezerapi.Context, name string) error
}

// storage puts strong typing around storage calls
type storage struct {
	rest.Creater
	rest.GracefulDeleter
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched
// types will panic.
func NewRegistry(c rest.Creater, d rest.GracefulDeleter) Registry {
	return &storage{
		Creater:         c,
		GracefulDeleter: d,
	}
}

func (s *storage) CreateNodeUser(ctx freezerapi.Context, user *api.NodeUser) (*api.NodeUser, error) {
	obj, err := s.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return obj.(*api.NodeUser), nil
}

func (s *storage) DeleteNodeUser(ctx freezerapi.Context, name string) error {
	_, err := s.Delete(ctx, name, nil)
	return err
}
