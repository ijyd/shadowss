package batchserver

import (
	"fmt"

	batch "cloud-keeper/pkg/apis/batch"
	batchvalidation "cloud-keeper/pkg/apis/batch/validation"

	"gofreezer/pkg/api"
	apierrors "gofreezer/pkg/api/errors"
	"gofreezer/pkg/runtime"
)

type REST struct {
}

func NewREST() *REST {
	return &REST{}
}

func (r *REST) New() runtime.Object {
	return &batch.BatchAccServer{}
}

func (r *REST) Create(ctx api.Context, obj runtime.Object) (runtime.Object, error) {
	accsrv, ok := obj.(*batch.BatchAccServer)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("not a SelfSubjectAccessReview: %#v", obj))
	}
	if errs := batchvalidation.ValidateBatchAccServer(accsrv); len(errs) > 0 {
		return nil, apierrors.NewInvalid(batch.Kind(accsrv.Kind), "", errs)
	}

	return accsrv, nil
}
