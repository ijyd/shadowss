package batch

import (
	api "gofreezer/pkg/api"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/runtime/schema"
)

const GroupName = "batch.keeper"

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
		&BatchAccServer{},
		&BatchResumeUsers{},
	)

	scheme.AddKnownTypes(SchemeGroupVersion,
		&api.ListOptions{},
		&api.DeleteOptions{},
	)
	return nil
}

func (obj *BatchAccServer) GetObjectKind() schema.ObjectKind { return &obj.TypeMeta }
