package innerhook

import "gofreezer/pkg/auth/user"

type InnerHookFunc func(value string) (user.Info, bool, error)

type InnerHookAuthenticator struct {
	HookFunc InnerHookFunc
}

func (hook InnerHookAuthenticator) AuthenticateToken(value string) (user.Info, bool, error) {
	return hook.HookFunc(value)
}
