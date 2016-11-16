package nodeuser

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
type shadowuserStrategy struct {
	runtime.ObjectTyper
	freezerapi.NameGenerator
}

// Strategy is the default logic that applies when creating and updating
// StorageClass objects via the REST API.
var Strategy = shadowuserStrategy{api.Scheme, freezerapi.SimpleNameGenerator}

func (shadowuserStrategy) NamespaceScoped() bool {
	return false
}

// ResetBeforeCreate clears the Status field which is not allowed to be set by end users on creation.
func (shadowuserStrategy) PrepareForCreate(ctx freezerapi.Context, obj runtime.Object) {
	_ = obj.(*api.NodeUser)
}

func (shadowuserStrategy) Validate(ctx freezerapi.Context, obj runtime.Object) field.ErrorList {
	node := obj.(*api.NodeUser)
	return validation.ValidateNodeUser(node)
}

// Canonicalize normalizes the object after validation.
func (shadowuserStrategy) Canonicalize(obj runtime.Object) {
}

func (shadowuserStrategy) AllowCreateOnUpdate() bool {
	return true
}

// PrepareForUpdate sets the Status fields which is not allowed to be set by an end user updating a PV
func (shadowuserStrategy) PrepareForUpdate(ctx freezerapi.Context, obj, old runtime.Object) {
	_ = obj.(*api.NodeUser)
	_ = old.(*api.NodeUser)
	// PadObj(obj)
	// PadObj(old)
}

func (shadowuserStrategy) ValidateUpdate(ctx freezerapi.Context, obj, old runtime.Object) field.ErrorList {
	errorList := validation.ValidateNodeUser(obj.(*api.NodeUser))
	return append(errorList, validation.ValidateNodeUserUpdate(obj.(*api.NodeUser), old.(*api.NodeUser))...)
}

func (shadowuserStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// MatchNode returns a generic matcher for a given label and field selector.
func MatchNodeUser(label labels.Selector, field fields.Selector) apistorage.SelectionPredicate {
	return apistorage.SelectionPredicate{
		Label: label,
		Field: field,
		GetAttrs: func(obj runtime.Object) (labels.Set, fields.Set, error) {
			cls, ok := obj.(*api.NodeUser)
			if !ok {
				return nil, nil, fmt.Errorf("given object is not of type TestType")
			}

			return labels.Set(cls.ObjectMeta.Labels), StorageClassToSelectableFields(cls), nil
		},
	}
}

// StorageClassToSelectableFields returns a label set that represents the object
func StorageClassToSelectableFields(data *api.NodeUser) fields.Set {
	return generic.ObjectMetaFieldsSet(&data.ObjectMeta, false)
}
