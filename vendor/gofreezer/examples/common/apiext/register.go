package apiext

import (
	"gofreezer/pkg/api"
	"gofreezer/pkg/api/unversioned"
	"gofreezer/pkg/runtime"
)

var Unversioned = unversioned.GroupVersion{Group: "", Version: "v1"}

var Scheme = api.Scheme

const GroupName = api.GroupName

var Codecs = api.Codecs

var SchemeGroupVersion = api.SchemeGroupVersion

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes, addDefaultingFuncs, addConversionFuncs, addFastPathConversionFuncs)
	AddToScheme   = AddALLToScheme //SchemeBuilder.AddToScheme
)

func init() {
	api.InitInternalAPI(Unversioned)
}

// AddToScheme applies all the stored functions to the scheme.
// this schem in api package
func AddALLToScheme() error {
	//add api scheme
	if err := api.AddToScheme(api.Scheme); err != nil {
		// Programmer error, detect immediately
		panic(err)
	}

	//add customes scheme
	if err := SchemeBuilder.AddToScheme(api.Scheme); err != nil {
		// Programmer error, detect immediately
		panic(err)
	}

	return nil
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(api.SchemeGroupVersion,
		&Login{},
		&LoginList{},
	)
	return nil
}
