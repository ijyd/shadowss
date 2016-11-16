package userservice

import (
	"fmt"

	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/fields"
	"gofreezer/pkg/labels"
	"gofreezer/pkg/runtime"
	apistorage "gofreezer/pkg/storage"
	"gofreezer/pkg/util/validation/field"

	"apistack/pkg/registry/generic"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/api/validation"
)

// storageClassStrategy implements behavior for StorageClass objects
type userserviceStrategy struct {
	runtime.ObjectTyper
	freezerapi.NameGenerator
}

// Strategy is the default logic that applies when creating and updating
// StorageClass objects via the REST API.
var Strategy = userserviceStrategy{api.Scheme, freezerapi.SimpleNameGenerator}

func (userserviceStrategy) NamespaceScoped() bool {
	return false
}

// ResetBeforeCreate clears the Status field which is not allowed to be set by end users on creation.
func (userserviceStrategy) PrepareForCreate(ctx freezerapi.Context, obj runtime.Object) {
	_ = obj.(*api.UserService)
}

func (userserviceStrategy) Validate(ctx freezerapi.Context, obj runtime.Object) field.ErrorList {
	user := obj.(*api.UserService)
	return validation.ValidateUserService(user)
}

// Canonicalize normalizes the object after validation.
func (userserviceStrategy) Canonicalize(obj runtime.Object) {
}

func (userserviceStrategy) AllowCreateOnUpdate() bool {
	return false
}

// PrepareForUpdate sets the Status fields which is not allowed to be set by an end user updating a PV
func (userserviceStrategy) PrepareForUpdate(ctx freezerapi.Context, obj, old runtime.Object) {
	_ = obj.(*api.UserService)
	_ = old.(*api.UserService)
}

func (userserviceStrategy) ValidateUpdate(ctx freezerapi.Context, obj, old runtime.Object) field.ErrorList {
	errorList := validation.ValidateUserService(obj.(*api.UserService))
	return append(errorList, validation.ValidateUserServiceUpdate(obj.(*api.UserService), old.(*api.UserService))...)
}

func (userserviceStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// MatchUserService returns a generic matcher for a given label and field selector.
func MatchUserService(label labels.Selector, field fields.Selector) apistorage.SelectionPredicate {
	return apistorage.SelectionPredicate{
		Label: label,
		Field: field,
		GetAttrs: func(obj runtime.Object) (labels.Set, fields.Set, error) {
			cls, ok := obj.(*api.UserService)
			if !ok {
				return nil, nil, fmt.Errorf("given object is not of type TestType")
			}

			return labels.Set(cls.ObjectMeta.Labels), StorageClassToSelectableFields(cls), nil
		},
	}
}

// StorageClassToSelectableFields returns a label set that represents the object
func StorageClassToSelectableFields(user *api.UserService) fields.Set {
	return generic.ObjectMetaFieldsSet(&user.ObjectMeta, false)
}
