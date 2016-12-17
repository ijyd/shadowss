package role

import (
	"errors"

	freezerapi "gofreezer/pkg/api"
	apierrors "gofreezer/pkg/api/errors"
	"gofreezer/pkg/runtime"

	"apistack/pkg/apis/abac"
	authorizes_abac "apistack/pkg/auth/authorizer/abac"
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/apis/abacpolicys"
)

const (
	wildcardsStar = "*"
)

type REST struct {
	policyFilePath string
}

func NewREST(path string) *REST {
	return &REST{path}
}

func (r *REST) New() runtime.Object {
	return &abacpolicys.Policy{}
}

func (r *REST) NewList() runtime.Object {
	return &abacpolicys.PolicyList{}
}

func subjectMatch(p abac.Policy, username string, groups []string) (matched bool) {

	// If the policy specified a user, ensure it matches
	if len(p.Spec.User) > 0 {
		if p.Spec.User == wildcardsStar {
			matched = true
		} else {
			matched = p.Spec.User == username
			if !matched {
				return false
			}
		}
	}

	// If the policy specified a group, ensure it matches
	if len(p.Spec.Group) > 0 {
		if p.Spec.Group == wildcardsStar {
			matched = true
		} else {
			matched = false
			for _, group := range groups {
				if p.Spec.Group == group {
					matched = true
				}
			}
			if !matched {
				return false
			}
		}
	}

	return matched
}

func (r *REST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {
	user, ok := freezerapi.UserFrom(ctx)
	if !ok {
		return nil, apierrors.NewForbidden(api.Resource("policys"), "", errors.New("invalid user"))
	}

	policyList := authorizes_abac.GetPolicys()

	username := user.GetName()
	groups := user.GetGroups()
	policys := &abacpolicys.PolicyList{}
	for _, v := range policyList {
		if subjectMatch(*v, username, groups) {
			policys.Items = append(policys.Items, *v)
		}
	}

	return policys, nil
}
