package masterhook

import (
	"gofreezer/pkg/api"
	freezeruser "gofreezer/pkg/auth/user"

	"github.com/golang/glog"

	corerest "cloud-keeper/pkg/registry/core/rest"
	"cloud-keeper/pkg/registry/core/user"
	"cloud-keeper/pkg/registry/core/usertoken"
)

const (
	GeneralGroup = "general"
	OPSGroup     = "ops"
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
	//token, err := c.UserTokenRegistry.GetUserTokenByToken(ctx, value)
	token, err := c.UserTokenRegistry.GetUserToken(ctx, value)
	if err != nil {
		return nil, false, err
	}

	info := &freezeruser.DefaultInfo{
		Name: token.Spec.Name,
	}

	switch token.Spec.Name {
	case "admin":
		info.Groups = append(info.Groups, freezeruser.SystemPrivilegedGroup)
	case "ops":
		info.Groups = append(info.Groups, OPSGroup)
	default:
		info.Groups = append(info.Groups, GeneralGroup)
	}

	glog.V(5).Infof("authenticate with innerHook(%v) result %+v", value, info)

	return info, true, nil
}
