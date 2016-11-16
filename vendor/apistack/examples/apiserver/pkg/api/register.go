package api

import (
	prototype "gofreezer/pkg/api"
	"gofreezer/pkg/api/unversioned"
	"gofreezer/pkg/runtime"
)

var Unversioned = unversioned.GroupVersion{Group: "", Version: "v1"}

var Scheme = prototype.Scheme

const GroupName = prototype.GroupName

var Codecs = prototype.Codecs

var SchemeGroupVersion = prototype.SchemeGroupVersion

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = AddALLToScheme //SchemeBuilder.AddToScheme
)

func init() {
	prototype.InitInternalAPI(Unversioned)
}

// AddToScheme applies all the stored functions to the scheme.
// this schem in prototype package
func AddALLToScheme() error {
	//add prototype scheme
	if err := prototype.AddToScheme(prototype.Scheme); err != nil {
		// Programmer error, detect immediately
		panic(err)
	}

	//add customes scheme
	if err := SchemeBuilder.AddToScheme(prototype.Scheme); err != nil {
		// Programmer error, detect immediately
		panic(err)
	}

	return nil
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(prototype.SchemeGroupVersion,
		&Login{},
		&LoginList{},
		&User{},
		&UserList{},
		&UserToken{},
		&UserTokenList{},
	)
	return nil
}
