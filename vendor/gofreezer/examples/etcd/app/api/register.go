/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/api/unversioned"
	"gofreezer/pkg/runtime"
)

var Unversioned = unversioned.GroupVersion{Group: "", Version: "v1"}

var Scheme = prototype.Scheme

const GroupName = prototype.GroupName

var Codecs = prototype.Codecs

var SchemeGroupVersion = prototype.SchemeGroupVersion

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes, addDefaultingFuncs, addConversionFuncs)
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
	)
	return nil
}
