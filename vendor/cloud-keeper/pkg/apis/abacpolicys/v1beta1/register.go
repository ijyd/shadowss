package v1beta1

import (
	"gofreezer/pkg/api/v1"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/runtime/schema"

	versionedwatch "gofreezer/pkg/watch/versioned"
)

const GroupName = "abac.keeper"

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1beta1"}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Adds the list of known types to api.Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Policy{},
		&PolicyList{},
	)

	scheme.AddKnownTypes(SchemeGroupVersion,
		&v1.ListOptions{},
	)
	versionedwatch.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
