package nodeuser

import (
	"cloud-keeper/pkg/api"
	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/watch"
)

// Registry is an interface for things that know how to store node.
type Registry interface {
	//CreateNodeUser(ctx freezerapi.Context, user *api.NodeUser) (*api.NodeUser, error)
	// DeleteNodeUser(ctx freezerapi.Context, name string) error
	UpdateNodeUser(ctx freezerapi.Context, token *api.NodeUser) (*api.NodeUser, error)
	WatchNodeUsers(ctx freezerapi.Context, options *freezerapi.ListOptions) (watch.Interface, error)
}

// storage puts strong typing around storage calls
type storage struct {
	//rest.Creater
	// rest.GracefulDeleter
	rest.Updater
	rest.Watcher
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched
// types will panic.
func NewRegistry(w rest.Watcher, u rest.Updater) Registry {
	return &storage{
		Watcher: w,
		Updater: u,
	}
}

//
// func (s *storage) CreateNodeUser(ctx freezerapi.Context, user *api.NodeUser) (*api.NodeUser, error) {
// 	obj, err := s.Create(ctx, user)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return obj.(*api.NodeUser), nil
// }

//
// func (s *storage) DeleteNodeUser(ctx freezerapi.Context, name string) error {
// 	_, err := s.Delete(ctx, name, nil)
// 	return err
// }

func (s *storage) UpdateNodeUser(ctx freezerapi.Context, user *api.NodeUser) (*api.NodeUser, error) {
	obj, _, err := s.Update(ctx, user.Spec.NodeName, rest.DefaultUpdatedObjectInfo(user, api.Scheme))
	if err != nil {
		return nil, err
	}
	return obj.(*api.NodeUser), nil
}

func (s *storage) WatchNodeUsers(ctx freezerapi.Context, options *freezerapi.ListOptions) (watch.Interface, error) {
	return s.Watch(ctx, options)
}
