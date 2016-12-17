package v1alpha1

import (
	//"cloud-keeper/pkg/api/v1"

	"gofreezer/pkg/api/v1"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/runtime/schema"

	versionedwatch "gofreezer/pkg/watch/versioned"
)

const GroupName = "batch.keeper"

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha1"}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Adds the list of known types to api.Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&BatchAccServer{},
		&BatchResumeUsers{},
	)

	scheme.AddKnownTypes(SchemeGroupVersion,
		&v1.ListOptions{},
		&v1.DeleteOptions{},
	)
	versionedwatch.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

func (obj *BatchAccServer) GetObjectKind() schema.ObjectKind { return &obj.TypeMeta }
