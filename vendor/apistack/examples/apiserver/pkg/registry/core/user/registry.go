package user

import (
	"apistack/examples/apiserver/pkg/api"
	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/rest"
)

// Registry is an interface for things that know how to store node.
type Registry interface {
	GetUser(ctx freezerapi.Context, userName string) (*api.User, error)
}

// storage puts strong typing around storage calls
type storage struct {
	rest.StandardStorage
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched
// types will panic.
func NewRegistry(s rest.StandardStorage) Registry {
	return &storage{s}
}

func (s *storage) GetUser(ctx freezerapi.Context, name string) (*api.User, error) {
	obj, err := s.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return obj.(*api.User), nil
}
