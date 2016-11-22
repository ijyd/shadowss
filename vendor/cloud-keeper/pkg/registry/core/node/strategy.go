package node

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
type nodeStrategy struct {
	runtime.ObjectTyper
	freezerapi.NameGenerator
}

// Strategy is the default logic that applies when creating and updating
// StorageClass objects via the REST API.
var Strategy = nodeStrategy{api.Scheme, freezerapi.SimpleNameGenerator}

func (nodeStrategy) NamespaceScoped() bool {
	return false
}

// ResetBeforeCreate clears the Status field which is not allowed to be set by end users on creation.
func (nodeStrategy) PrepareForCreate(ctx freezerapi.Context, obj runtime.Object) {
	_ = obj.(*api.Node)
}

func (nodeStrategy) Validate(ctx freezerapi.Context, obj runtime.Object) field.ErrorList {
	node := obj.(*api.Node)
	return validation.ValidateNode(node)
}

// Canonicalize normalizes the object after validation.
func (nodeStrategy) Canonicalize(obj runtime.Object) {
}

func (nodeStrategy) AllowCreateOnUpdate() bool {
	return true
}

// PrepareForUpdate sets the Status fields which is not allowed to be set by an end user updating a PV
func (nodeStrategy) PrepareForUpdate(ctx freezerapi.Context, obj, old runtime.Object) {
	PadObj(obj)
	PadObj(old)

	node := obj.(*api.Node)
	oldnode := old.(*api.Node)

	node.UID = oldnode.UID
}

func (nodeStrategy) ValidateUpdate(ctx freezerapi.Context, obj, old runtime.Object) field.ErrorList {
	errorList := validation.ValidateNode(obj.(*api.Node))
	return append(errorList, validation.ValidateNodeUpdate(obj.(*api.Node), old.(*api.Node))...)
}

func (nodeStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// MatchNode returns a generic matcher for a given label and field selector.
func MatchNode(label labels.Selector, field fields.Selector) apistorage.SelectionPredicate {
	return apistorage.SelectionPredicate{
		Label: label,
		Field: field,
		GetAttrs: func(obj runtime.Object) (labels.Set, fields.Set, error) {
			cls, ok := obj.(*api.Node)
			if !ok {
				return nil, nil, fmt.Errorf("given object is not of type TestType")
			}

			return labels.Set(cls.ObjectMeta.Labels), StorageClassToSelectableFields(cls), nil
		},
	}
}

// StorageClassToSelectableFields returns a label set that represents the object
func StorageClassToSelectableFields(data *api.Node) fields.Set {
	return generic.ObjectMetaFieldsSet(&data.ObjectMeta, false)
}

func PadObj(obj runtime.Object) error {
	node := obj.(*api.Node)
	node.Name = node.Spec.Server.Name
	node.ResourceVersion = "1"
	return nil
}
