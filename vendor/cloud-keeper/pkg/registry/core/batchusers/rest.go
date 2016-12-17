package batchusers

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/user"

	"github.com/golang/glog"

	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/runtime"
)

type BatchUsersREST struct {
	user user.Registry
}

func NewREST(user user.Registry) *BatchUsersREST {
	return &BatchUsersREST{
		user: user,
	}
}

func (*BatchUsersREST) New() runtime.Object {
	return &api.BatchUsers{}
}

func (r *BatchUsersREST) Create(ctx freezerapi.Context, obj runtime.Object) (runtime.Object, error) {

	users := obj.(*api.BatchUsers)

	glog.V(5).Infof("post batch users %+v\r\n", *users)

	if users.Spec.Resume {
		err := r.user.ResumeUsers(ctx, users.Spec.TargetUser)
		if err != nil {
			return nil, err
		}
	}

	return users, nil
}
