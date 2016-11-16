package userfile

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

// loginStuserfileStrategyrategy implements behavior for userfile objects
type userfileStrategy struct {
	runtime.ObjectTyper
	freezerapi.NameGenerator
}

// Strategy is the default logic that applies when creating and updating
// StorageClass objects via the REST API.
var Strategy = userfileStrategy{api.Scheme, freezerapi.SimpleNameGenerator}

func (userfileStrategy) NamespaceScoped() bool {
	return false
}

// ResetBeforeCreate clears the Status field which is not allowed to be set by end users on creation.
func (userfileStrategy) PrepareForCreate(ctx freezerapi.Context, obj runtime.Object) {
	_ = obj.(*api.UserPublicFile)
}

func (userfileStrategy) Validate(ctx freezerapi.Context, obj runtime.Object) field.ErrorList {
	file := obj.(*api.UserPublicFile)
	return validation.ValidateUserPublicFile(file)
}

// Canonicalize normalizes the object after validation.
func (userfileStrategy) Canonicalize(obj runtime.Object) {
}

func (userfileStrategy) AllowCreateOnUpdate() bool {
	return false
}

// PrepareForUpdate sets the Status fields which is not allowed to be set by an end user updating a PV
func (userfileStrategy) PrepareForUpdate(ctx freezerapi.Context, obj, old runtime.Object) {
	_ = obj.(*api.UserPublicFile)
	_ = old.(*api.UserPublicFile)
	// PadObj(obj)
	// PadObj(old)
}

func (userfileStrategy) ValidateUpdate(ctx freezerapi.Context, obj, old runtime.Object) field.ErrorList {
	errorList := validation.ValidateUserPublicFile(obj.(*api.UserPublicFile))
	return append(errorList, validation.ValidateUserPublicFileUpdate(obj.(*api.UserPublicFile), old.(*api.UserPublicFile))...)
}

func (userfileStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// MatchLogin returns a generic matcher for a given label and field selector.
func MatchUserToken(label labels.Selector, field fields.Selector) apistorage.SelectionPredicate {
	return apistorage.SelectionPredicate{
		Label: label,
		Field: field,
		GetAttrs: func(obj runtime.Object) (labels.Set, fields.Set, error) {
			cls, ok := obj.(*api.UserPublicFile)
			if !ok {
				return nil, nil, fmt.Errorf("given object is not of type TestType")
			}

			return labels.Set(cls.ObjectMeta.Labels), StorageClassToSelectableFields(cls), nil
		},
	}
}

// StorageClassToSelectableFields returns a label set that represents the object
func StorageClassToSelectableFields(data *api.UserPublicFile) fields.Set {
	return generic.ObjectMetaFieldsSet(&data.ObjectMeta, false)
}

// func PadObj(obj runtime.Object) error {
// 	token := obj.(*api.UserToken)
// 	token.Name = token.Spec.Name
// 	token.ResourceVersion = "1"
// 	return nil
// }

// PortForwardLocation returns the port-forward URL for a pod.
// func UserFileLocation(name string) (*url.URL, http.RoundTripper, error) {
// 	loc := &url.URL{
// 		Scheme: "",
// 		Host:   net.JoinHostPort("192.168.60.128", 20008),
// 		Path:   fmt.Sprintf("/userdata/%s", name),
// 	}
// 	t := &http.Transport{}
// 	t.RegisterProtocol("file", http.NewFileTransport(http.Dir(fmt.Sprintf("/userdata/%s", name))))
// 	return loc, t, nil
// }
