package rest

import (
	"fmt"
	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/errors"
	"gofreezer/pkg/runtime"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/account"
)

type AccExtendREST struct {
	Exec    *ExecREST
	SSHKeys *SSHKeysREST
	AccInfo *AccInfoREST
}

func NewExtendREST(accRegistry account.Registry) *AccExtendREST {
	return &AccExtendREST{
		Exec:    &ExecREST{accRegistry},
		SSHKeys: &SSHKeysREST{accRegistry},
		AccInfo: &AccInfoREST{accRegistry},
	}
}

type AccInfoREST struct {
	accRegistry account.Registry
}

func (r *AccInfoREST) New() runtime.Object {
	return &api.AccountInfo{}
}

func (r *AccInfoREST) Get(ctx freezerapi.Context, name string) (runtime.Object, error) {
	collectorHandler, err := r.accRegistry.GetCloudProvider(ctx, name)
	if err == nil {
		infoSpec, err := collectorHandler.GetAccount()
		if err != nil {
			return nil, errors.NewBadRequest(err.Error())
		}

		accinfo := &api.AccountInfo{
			Spec: *infoSpec,
		}

		return accinfo, nil
	}
	return nil, errors.NewBadRequest(fmt.Sprintf("not found cloud provider by account(%v)", name))
}

type ExecREST struct {
	accRegistry account.Registry
}

func (r *ExecREST) New() runtime.Object {
	return &api.AccExec{}
}

func (r *ExecREST) Create(ctx freezerapi.Context, obj runtime.Object) (runtime.Object, error) {
	exec := obj.(*api.AccExec)
	targetName := exec.Spec.AccName
	collectorHandler, err := r.accRegistry.GetCloudProvider(ctx, targetName)
	if err == nil {
		err := collectorHandler.Exec(exec)
		if err != nil {
			return nil, errors.NewBadRequest(err.Error())
		}

		return exec, nil
	}
	return nil, errors.NewBadRequest(fmt.Sprintf("not found cloud provider by account(%v)", targetName))
}

type SSHKeysREST struct {
	accRegistry account.Registry
}

func (r *SSHKeysREST) New() runtime.Object {
	return &api.AccSSHKey{}
}

func (r *SSHKeysREST) Get(ctx freezerapi.Context, name string) (runtime.Object, error) {
	collectorHandler, err := r.accRegistry.GetCloudProvider(ctx, name)
	if err == nil {
		sshkeys, err := collectorHandler.GetSSHKey()
		if err != nil {
			return nil, errors.NewBadRequest(err.Error())
		}

		return sshkeys, nil
	}
	return nil, errors.NewBadRequest(fmt.Sprintf("not found cloud provider by account(%v)", name))
}
