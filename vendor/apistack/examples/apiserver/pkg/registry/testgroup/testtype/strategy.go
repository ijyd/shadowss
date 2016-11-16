package testtype

import (
	"fmt"

	"gofreezer/pkg/api"
	"gofreezer/pkg/fields"
	"gofreezer/pkg/labels"
	"gofreezer/pkg/runtime"
	apistorage "gofreezer/pkg/storage"
	"gofreezer/pkg/util/validation/field"

	"apistack/examples/apiserver/pkg/apis/testgroup"
	"apistack/examples/apiserver/pkg/apis/testgroup/validation"
	"apistack/pkg/registry/generic"
)

// storageClassStrategy implements behavior for StorageClass objects
type testTypeStrategy struct {
	runtime.ObjectTyper
	api.NameGenerator
}

// Strategy is the default logic that applies when creating and updating
// StorageClass objects via the REST API.
var Strategy = testTypeStrategy{api.Scheme, api.SimpleNameGenerator}

func (testTypeStrategy) NamespaceScoped() bool {
	return false
}

// ResetBeforeCreate clears the Status field which is not allowed to be set by end users on creation.
func (testTypeStrategy) PrepareForCreate(ctx api.Context, obj runtime.Object) {
	_ = obj.(*testgroup.TestType)
}

func (testTypeStrategy) Validate(ctx api.Context, obj runtime.Object) field.ErrorList {
	testType := obj.(*testgroup.TestType)
	return validation.ValidateTestType(testType)
}

// Canonicalize normalizes the object after validation.
func (testTypeStrategy) Canonicalize(obj runtime.Object) {
}

func (testTypeStrategy) AllowCreateOnUpdate() bool {
	return false
}

// PrepareForUpdate sets the Status fields which is not allowed to be set by an end user updating a PV
func (testTypeStrategy) PrepareForUpdate(ctx api.Context, obj, old runtime.Object) {
	_ = obj.(*testgroup.TestType)
	_ = old.(*testgroup.TestType)
}

func (testTypeStrategy) ValidateUpdate(ctx api.Context, obj, old runtime.Object) field.ErrorList {
	errorList := validation.ValidateTestType(obj.(*testgroup.TestType))
	return append(errorList, validation.ValidateTestTypeUpdate(obj.(*testgroup.TestType), old.(*testgroup.TestType))...)
}

func (testTypeStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// MatchStorageClass returns a generic matcher for a given label and field selector.
func MatchStorageClasses(label labels.Selector, field fields.Selector) apistorage.SelectionPredicate {
	return apistorage.SelectionPredicate{
		Label: label,
		Field: field,
		GetAttrs: func(obj runtime.Object) (labels.Set, fields.Set, error) {
			cls, ok := obj.(*testgroup.TestType)
			if !ok {
				return nil, nil, fmt.Errorf("given object is not of type TestType")
			}

			return labels.Set(cls.ObjectMeta.Labels), StorageClassToSelectableFields(cls), nil
		},
	}
}

// StorageClassToSelectableFields returns a label set that represents the object
func StorageClassToSelectableFields(test *testgroup.TestType) fields.Set {
	return generic.ObjectMetaFieldsSet(&test.ObjectMeta, false)
}
