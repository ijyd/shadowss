package usertoken

import (
	"cloud-keeper/pkg/api"
	"fmt"
	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/fields"
)

// Registry is an interface for things that know how to store node.
type Registry interface {
	GetUserTokenByToken(ctx freezerapi.Context, token string) (*api.UserToken, error)
	UpdateUserToken(ctx freezerapi.Context, token *api.UserToken) (*api.UserToken, error)
}

// storage puts strong typing around storage calls
type storage struct {
	rest.Getter
	rest.Updater
	rest.Lister
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched
// types will panic.
func NewRegistry(s rest.Getter, c rest.Updater, l rest.Lister) Registry {
	return &storage{
		Getter:  s,
		Updater: c,
		Lister:  l,
	}
}

func (s *storage) GetUserTokenByToken(ctx freezerapi.Context, token string) (*api.UserToken, error) {
	options := &freezerapi.ListOptions{}
	options.FieldSelector = fields.ParseSelectorOrDie(fmt.Sprintf("spec.token=%s", token))
	obj, err := s.List(ctx, options)
	if err != nil {
		return nil, err
	}

	list := obj.(*api.UserTokenList)
	if len(list.Items) == 1 {
		return &list.Items[0], nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func (s *storage) UpdateUserToken(ctx freezerapi.Context, token *api.UserToken) (*api.UserToken, error) {
	obj, _, err := s.Update(ctx, token.Name, rest.DefaultUpdatedObjectInfo(token, api.Scheme))

	if err != nil {
		return nil, err
	}
	newToken, ok := obj.(*api.UserToken)
	if !ok {
		return nil, fmt.Errorf("not a UserToken Object")
	}

	return newToken, nil
}
