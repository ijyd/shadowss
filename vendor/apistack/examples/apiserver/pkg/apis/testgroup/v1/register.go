package v1

import (
	"apistack/examples/apiserver/pkg/api/v1"

	"gofreezer/pkg/api/unversioned"
	"gofreezer/pkg/runtime"

	versionedwatch "gofreezer/pkg/watch/versioned"
)

const GroupName = "testgroup"

var SchemeGroupVersion = unversioned.GroupVersion{Group: GroupName, Version: "v1"}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Adds the list of known types to api.Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&TestType{},
		&TestTypeList{},
	)

	scheme.AddKnownTypes(SchemeGroupVersion,
		&v1.ListOptions{},
		&v1.DeleteOptions{},
		&unversioned.Status{},
		&v1.ExportOptions{},
	)
	versionedwatch.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

func (obj *TestType) GetObjectKind() unversioned.ObjectKind     { return &obj.TypeMeta }
func (obj *TestTypeList) GetObjectKind() unversioned.ObjectKind { return &obj.TypeMeta }
