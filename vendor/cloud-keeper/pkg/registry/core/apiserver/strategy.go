package apiserver

import (
	"fmt"

	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/fields"
	"gofreezer/pkg/labels"
	"gofreezer/pkg/pages"
	"gofreezer/pkg/runtime"
	apistorage "gofreezer/pkg/storage"
	"gofreezer/pkg/util/validation/field"

	"apistack/pkg/registry/generic"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/api/validation"
)

// apiserverStrategy implements behavior for Login objects
type apiserverStrategy struct {
	runtime.ObjectTyper
	freezerapi.NameGenerator
}

// Strategy is the default logic that applies when creating and updating
// StorageClass objects via the REST API.
var Strategy = apiserverStrategy{api.Scheme, freezerapi.SimpleNameGenerator}

func (apiserverStrategy) NamespaceScoped() bool {
	return false
}

// ResetBeforeCreate clears the Status field which is not allowed to be set by end users on creation.
func (apiserverStrategy) PrepareForCreate(ctx freezerapi.Context, obj runtime.Object) {
	_ = obj.(*api.APIServer)
}

func (apiserverStrategy) Validate(ctx freezerapi.Context, obj runtime.Object) field.ErrorList {
	server := obj.(*api.APIServer)
	return validation.ValidateAPIServer(server)
}

// Canonicalize normalizes the object after validation.
func (apiserverStrategy) Canonicalize(obj runtime.Object) {
}

func (apiserverStrategy) AllowCreateOnUpdate() bool {
	return false
}

// PrepareForUpdate sets the Status fields which is not allowed to be set by an end user updating a PV
func (apiserverStrategy) PrepareForUpdate(ctx freezerapi.Context, obj, old runtime.Object) {
	_ = obj.(*api.Node)
	_ = old.(*api.Node)
	// PadObj(obj)
	// PadObj(old)
}

func (apiserverStrategy) ValidateUpdate(ctx freezerapi.Context, obj, old runtime.Object) field.ErrorList {
	errorList := validation.ValidateAPIServer(obj.(*api.APIServer))
	return append(errorList, validation.ValidateAPIServerUpdate(obj.(*api.APIServer), old.(*api.APIServer))...)
}

func (apiserverStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// MatchAPIServer returns a generic matcher for a given label and field selector.
func MatchAPIServer(label labels.Selector, field fields.Selector, page pages.Selector) apistorage.SelectionPredicate {
	return apistorage.SelectionPredicate{
		Label: label,
		Field: field,
		Page:  page,
		GetAttrs: func(obj runtime.Object) (labels.Set, fields.Set, error) {
			cls, ok := obj.(*api.APIServer)
			if !ok {
				return nil, nil, fmt.Errorf("given object is not of type TestType")
			}

			return labels.Set(cls.ObjectMeta.Labels), StorageClassToSelectableFields(cls), nil
		},
	}
}

// StorageClassToSelectableFields returns a label set that represents the object
func StorageClassToSelectableFields(data *api.APIServer) fields.Set {
	return generic.ObjectMetaFieldsSet(&data.ObjectMeta, false)
}

// func PadObj(obj runtime.Object) error {
// 	node := obj.(*api.Node)
// 	node.Name = node.Spec.Server.Name
// 	node.ResourceVersion = "1"
// 	return nil
// }
