package master

import (
	"gofreezer/pkg/api"
	freezeruser "gofreezer/pkg/auth/user"

	corerest "cloud-keeper/pkg/registry/core/rest"
	"cloud-keeper/pkg/registry/core/user"
	"cloud-keeper/pkg/registry/core/usertoken"
)

type InnerHook struct {
	UserRegistry      user.Registry
	UserTokenRegistry usertoken.Registry
}

var InnerHookHandler = NewInnerHook()

func NewInnerHook() *InnerHook {
	return &InnerHook{}
}

func (c *InnerHook) SetRegistry(legacyRESTStorage corerest.LegacyRESTStorage) {
	c.UserRegistry = legacyRESTStorage.UserRegistry
	c.UserTokenRegistry = legacyRESTStorage.TokenRegistry
}

func (c *InnerHook) AuthenticateTokenInnerHook(value string) (freezeruser.Info, bool, error) {
	ctx := api.NewContext()
	token, err := c.UserTokenRegistry.GetUserTokenByToken(ctx, value)
	if err != nil {
		return nil, false, err
	}

	info := &freezeruser.DefaultInfo{
		Name: token.Spec.Name,
	}
	return info, true, nil
}
