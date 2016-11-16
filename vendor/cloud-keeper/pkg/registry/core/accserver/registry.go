package accserver

import (
	"cloud-keeper/pkg/api"
	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/rest"
)

// Registry is an interface for things that know how to store node.
type Registry interface {
	GetAccServer(ctx freezerapi.Context, name string) (*api.AccServer, error)
	// ListAccounts(ctx freezerapi.Context, options *freezerapi.ListOptions) (*api.AccountList, error)
	// CreateAccount(ctx freezerapi.Context, acc *api.Account) (*api.Account, error)
	// DeleteAccount(ctx freezerapi.Context, name string) error
}

// storage puts strong typing around storage calls
type storage struct {
	rest.Getter
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched
// types will panic.
func NewRegistry(g rest.Getter) Registry {
	return &storage{
		Getter: g,
	}
}

func (s *storage) GetAccServer(ctx freezerapi.Context, name string) (*api.AccServer, error) {
	obj, err := s.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return obj.(*api.AccServer), nil
}

// func (s *storage) ListAccounts(ctx freezerapi.Context, options *freezerapi.ListOptions) (*api.AccountList, error) {
// 	obj, err := s.List(ctx, options)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return obj.(*api.AccountList), nil
// }
//
// func (s *storage) CreateAccount(ctx freezerapi.Context, acc *api.Account) (*api.Account, error) {
// 	obj, err := s.Create(ctx, acc)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return obj.(*api.Account), nil
// }
//
// func (s *storage) DeleteAccount(ctx freezerapi.Context, name string) error {
// 	_, err := s.Delete(ctx, name)
// 	return err
// }
