package v1

import "gofreezer/pkg/runtime"

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return scheme.AddDefaultingFuncs(SetDefaults_Login)
}

func SetDefaults_Login(obj *Login) {
	obj.Spec.Auth = ""
	obj.Spec.Token = ""
}
