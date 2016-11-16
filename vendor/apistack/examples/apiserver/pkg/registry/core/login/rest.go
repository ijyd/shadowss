package login

import (
	"apistack/examples/apiserver/pkg/api"
	"apistack/examples/apiserver/pkg/registry/core/user"
	"apistack/examples/apiserver/pkg/registry/core/usertoken"
	"crypto/rand"
	"fmt"
	freezerapi "gofreezer/pkg/api"
	apierrors "gofreezer/pkg/api/errors"
	"gofreezer/pkg/api/unversioned"
	"gofreezer/pkg/runtime"
	"time"
)

var groupResource = api.Resource("logins")

// REST implements the REST endpoint for login
type REST struct {
	user      user.Registry
	usertoken usertoken.Registry
}

func NewREST(user user.Registry, token usertoken.Registry) *REST {
	return &REST{
		user:      user,
		usertoken: token,
	}
}

func (r *REST) New() runtime.Object {
	return &api.Login{}
}

func (r *REST) Create(ctx freezerapi.Context, obj runtime.Object) (runtime.Object, error) {
	login, ok := obj.(*api.Login)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("not a Login: %#v", obj))
	}
	user, err := r.user.GetUser(ctx, login.Spec.AuthName)
	if err != nil {
		return nil, err
	}

	if login.Spec.Auth != user.Spec.DetailInfo.ManagePasswd {
		return nil, apierrors.NewForbidden(groupResource, user.Name, fmt.Errorf("fobidden"))
	}

	token, err := randBearerToken()
	if err != nil {
		return nil, err
	}

	usertoken := &api.UserToken{
		ObjectMeta: freezerapi.ObjectMeta{
			Name: user.Spec.DetailInfo.Name,
		},
		Spec: api.UserTokenSpec{
			Token:      token,
			UserID:     user.Spec.DetailInfo.ID,
			Name:       user.Spec.DetailInfo.Name,
			CreateTime: unversioned.NewTime(time.Now()),
			ExpireTime: unversioned.NewTime(time.Now().Add(time.Duration(1) * time.Hour)),
		},
	}

	err = r.usertoken.UpdateUserToken(ctx, usertoken)

	return login, err
}

func randBearerToken() (string, error) {
	token := make([]byte, 16)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", token), err
}
