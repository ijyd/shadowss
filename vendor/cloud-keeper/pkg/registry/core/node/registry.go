package node

import (
	"cloud-keeper/pkg/api"
	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/labels"
	"gofreezer/pkg/watch"
)

// Registry is an interface for things that know how to store node.
type Registry interface {
	GetAPINodes(ctx freezerapi.Context, options *freezerapi.ListOptions) (*api.NodeList, error)
	GetNode(ctx freezerapi.Context, name string) (*api.Node, error)
	ListNodes(ctx freezerapi.Context, options *freezerapi.ListOptions) (*api.NodeList, error)
	UpdateNode(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (*api.Node, bool, error)
	WatchNodes(ctx freezerapi.Context, options *freezerapi.ListOptions) (watch.Interface, error)
	//UpdateNodeUser(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (*api.Node, bool, error)
}

// storage puts strong typing around storage calls
type storage struct {
	rest.Updater
	rest.Getter
	rest.Lister
	rest.Watcher
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched
// types will panic.
func NewRegistry(u rest.Updater, g rest.Getter, l rest.Lister, w rest.Watcher) Registry {
	return &storage{
		Updater: u,
		Getter:  g,
		Lister:  l,
		Watcher: w,
	}
}

func (s *storage) UpdateNode(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (*api.Node, bool, error) {
	obj, flag, err := s.Update(ctx, name, objInfo)
	if err != nil {
		return nil, false, err
	}
	return obj.(*api.Node), flag, nil
}

// func (s *storage) UpdateNodeUser(ctx freezerapi.Context, name string, objInfo rest.UpdatedObjectInfo) (*api.Node, bool, error) {
// 	return s.Update(ctx, name, objInfo)
// }

func (s *storage) GetAPINodes(ctx freezerapi.Context, options *freezerapi.ListOptions) (*api.NodeList, error) {
	options = &freezerapi.ListOptions{}
	ls := labels.Set(map[string]string{
		api.NodeLablesUserSpace: api.NodeUserSpaceAPI,
	})
	options.LabelSelector = labels.SelectorFromSet(ls)
	obj, err := s.List(ctx, options)
	if err != nil {
		return nil, err
	}

	return obj.(*api.NodeList), nil
}

func (s *storage) GetNode(ctx freezerapi.Context, name string) (*api.Node, error) {
	obj, err := s.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return obj.(*api.Node), nil
}

func (s *storage) ListNodes(ctx freezerapi.Context, options *freezerapi.ListOptions) (*api.NodeList, error) {
	obj, err := s.List(ctx, options)
	if err != nil {
		return nil, err
	}

	return obj.(*api.NodeList), nil
}

func (s *storage) WatchNodes(ctx freezerapi.Context, options *freezerapi.ListOptions) (watch.Interface, error) {
	return s.Watch(ctx, options)
}
