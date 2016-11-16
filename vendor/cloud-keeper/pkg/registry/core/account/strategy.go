package account

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

// loginStrategy implements behavior for Login objects
type accountStrategy struct {
	runtime.ObjectTyper
	freezerapi.NameGenerator
}

// Strategy is the default logic that applies when creating and updating
// StorageClass objects via the REST API.
var Strategy = accountStrategy{api.Scheme, freezerapi.SimpleNameGenerator}

func (accountStrategy) NamespaceScoped() bool {
	return false
}

// ResetBeforeCreate clears the Status field which is not allowed to be set by end users on creation.
func (accountStrategy) PrepareForCreate(ctx freezerapi.Context, obj runtime.Object) {
	_ = obj.(*api.Account)
}

func (accountStrategy) Validate(ctx freezerapi.Context, obj runtime.Object) field.ErrorList {
	acc := obj.(*api.Account)
	return validation.ValidateAccount(acc)
}

// Canonicalize normalizes the object after validation.
func (accountStrategy) Canonicalize(obj runtime.Object) {
}

func (accountStrategy) AllowCreateOnUpdate() bool {
	return true
}

// PrepareForUpdate sets the Status fields which is not allowed to be set by an end user updating a PV
func (accountStrategy) PrepareForUpdate(ctx freezerapi.Context, obj, old runtime.Object) {
	PadObj(obj)
	PadObj(old)
}

func (accountStrategy) ValidateUpdate(ctx freezerapi.Context, obj, old runtime.Object) field.ErrorList {
	errorList := validation.ValidateAccount(obj.(*api.Account))
	return append(errorList, validation.ValidateAccountUpdate(obj.(*api.Account), old.(*api.Account))...)
}

func (accountStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// MatchLogin returns a generic matcher for a given label and field selector.
func MatchAccount(label labels.Selector, field fields.Selector) apistorage.SelectionPredicate {
	return apistorage.SelectionPredicate{
		Label: label,
		Field: field,
		GetAttrs: func(obj runtime.Object) (labels.Set, fields.Set, error) {
			cls, ok := obj.(*api.Account)
			if !ok {
				return nil, nil, fmt.Errorf("given object is not of type TestType")
			}

			return labels.Set(cls.ObjectMeta.Labels), StorageClassToSelectableFields(cls), nil
		},
	}
}

// StorageClassToSelectableFields returns a label set that represents the object
func StorageClassToSelectableFields(data *api.Account) fields.Set {
	return generic.ObjectMetaFieldsSet(&data.ObjectMeta, false)
}

func PadObj(obj runtime.Object) error {
	acc := obj.(*api.Account)
	acc.Name = acc.Spec.AccDetail.Name
	acc.ResourceVersion = "1"
	return nil
}
