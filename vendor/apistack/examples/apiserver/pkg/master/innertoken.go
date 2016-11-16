package master

import (
	"gofreezer/pkg/api"
	freezeruser "gofreezer/pkg/auth/user"

	"github.com/golang/glog"

	corerest "apistack/examples/apiserver/pkg/registry/core/rest"
	"apistack/examples/apiserver/pkg/registry/core/user"
	"apistack/examples/apiserver/pkg/registry/core/usertoken"
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
		glog.Infof("authenticate with innerHook  %v and token %+v err:%v", value, token, err)
		return nil, false, nil
	}

	info := &freezeruser.DefaultInfo{
		Name: token.Spec.Name,
	}
	glog.V(5).Infof("authenticate with innerHook %v and token %+v", value, token)

	return info, true, nil
}
