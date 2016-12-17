package accserver

import (
	"cloud-keeper/pkg/api"
	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/rest"
)

// Registry is an interface for things that know how to store node.
type Registry interface {
	//GetAccServer(ctx freezerapi.Context, name string) (*api.AccServer, error)
	ListAccServers(ctx freezerapi.Context, options *freezerapi.ListOptions) (*api.AccServerList, error)
}

// storage puts strong typing around storage calls
type storage struct {
	//rest.Getter
	rest.Lister
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched
// types will panic.
func NewRegistry(l rest.Lister) Registry {
	return &storage{
		Lister: l,
	}
}

// func (s *storage) GetAccServer(ctx freezerapi.Context, name string) (*api.AccServer, error) {
// 	obj, err := s.Get(ctx, name)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return obj.(*api.AccServer), nil
// }

func (s *storage) ListAccServers(ctx freezerapi.Context, options *freezerapi.ListOptions) (*api.AccServerList, error) {
	obj, err := s.List(ctx, options)
	if err != nil {
		return nil, err
	}
	return obj.(*api.AccServerList), nil
}
