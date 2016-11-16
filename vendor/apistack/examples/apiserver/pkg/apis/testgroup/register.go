package testgroup

import (
	api "gofreezer/pkg/api"

	"gofreezer/pkg/api/unversioned"
	"gofreezer/pkg/runtime"
)

const GroupName = "testgroup"

var SchemeGroupVersion = unversioned.GroupVersion{Group: GroupName, Version: runtime.APIVersionInternal}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Kind takes an unqualified kind and returns a Group qualified GroupKind
func Kind(kind string) unversioned.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) unversioned.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// Adds the list of known types to api.Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&TestType{},
		&TestTypeList{},
	)

	scheme.AddKnownTypes(SchemeGroupVersion,
		&api.ListOptions{},
	)
	return nil
}

func (obj *TestType) GetObjectKind() unversioned.ObjectKind     { return &obj.TypeMeta }
func (obj *TestTypeList) GetObjectKind() unversioned.ObjectKind { return &obj.TypeMeta }
