package usertoken

import (
	"fmt"

	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/fields"
	"gofreezer/pkg/labels"
	"gofreezer/pkg/pagination"
	"gofreezer/pkg/runtime"
	apistorage "gofreezer/pkg/storage"
	"gofreezer/pkg/util/validation/field"

	"apistack/pkg/registry/generic"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/api/validation"
)

// loginStrategy implements behavior for Login objects
type tokenStrategy struct {
	runtime.ObjectTyper
	freezerapi.NameGenerator
}

// Strategy is the default logic that applies when creating and updating
// StorageClass objects via the REST API.
var Strategy = tokenStrategy{api.Scheme, freezerapi.SimpleNameGenerator}

func (tokenStrategy) NamespaceScoped() bool {
	return false
}

// ResetBeforeCreate clears the Status field which is not allowed to be set by end users on creation.
func (tokenStrategy) PrepareForCreate(ctx freezerapi.Context, obj runtime.Object) {
	_ = obj.(*api.UserToken)
}

func (tokenStrategy) Validate(ctx freezerapi.Context, obj runtime.Object) field.ErrorList {
	token := obj.(*api.UserToken)
	return validation.ValidateUserToken(token)
}

// Canonicalize normalizes the object after validation.
func (tokenStrategy) Canonicalize(obj runtime.Object) {
}

func (tokenStrategy) AllowCreateOnUpdate() bool {
	return true
}

// PrepareForUpdate sets the Status fields which is not allowed to be set by an end user updating a PV
func (tokenStrategy) PrepareForUpdate(ctx freezerapi.Context, obj, old runtime.Object) {
	PadObj(obj)
	PadObj(old)
}

func (tokenStrategy) ValidateUpdate(ctx freezerapi.Context, obj, old runtime.Object) field.ErrorList {
	errorList := validation.ValidateUserToken(obj.(*api.UserToken))
	return append(errorList, validation.ValidateUserTokenUpdate(obj.(*api.UserToken), old.(*api.UserToken))...)
}

func (tokenStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// MatchLogin returns a generic matcher for a given label and field selector.
func MatchUserToken(label labels.Selector, field fields.Selector, page pagination.Pager) apistorage.SelectionPredicate {
	return apistorage.SelectionPredicate{
		Label: label,
		Field: field,
		Pager: page,
		GetAttrs: func(obj runtime.Object) (labels.Set, fields.Set, error) {
			cls, ok := obj.(*api.UserToken)
			if !ok {
				return nil, nil, fmt.Errorf("given object is not of type TestType")
			}

			return labels.Set(cls.ObjectMeta.Labels), StorageClassToSelectableFields(cls), nil
		},
	}
}

// StorageClassToSelectableFields returns a label set that represents the object
func StorageClassToSelectableFields(data *api.UserToken) fields.Set {
	return generic.ObjectMetaFieldsSet(&data.ObjectMeta, false)
}

func PadObj(obj runtime.Object) error {
	token := obj.(*api.UserToken)
	token.Name = token.Spec.Name
	token.ResourceVersion = "1"
	return nil
}
