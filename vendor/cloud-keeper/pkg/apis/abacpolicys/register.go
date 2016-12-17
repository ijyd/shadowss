package abacpolicys

import (
	api "gofreezer/pkg/api"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/runtime/schema"
)

const GroupName = "abac.keeper"

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: runtime.APIVersionInternal}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Kind takes an unqualified kind and returns a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// Adds the list of known types to api.Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Policy{},
		&PolicyList{},
	)

	scheme.AddKnownTypes(SchemeGroupVersion,
		&api.ListOptions{},
	)
	return nil
}
