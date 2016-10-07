package api

import (
	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/conversion"
	"gofreezer/pkg/runtime"
	"reflect"
)

func init() {
	SchemeBuilder.Register(RegisterDeepCopies)
}

func RegisterDeepCopies(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedDeepCopyFuncs(
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_api_Login, InType: reflect.TypeOf(&Login{})},
	)
}

func DeepCopy_api_LoginSpec(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*LoginSpec)
		out := out.(*LoginSpec)
		out.Auth = in.Auth
		out.AuthName = in.AuthName
		out.Token = in.Token
		return nil
	}
}

func DeepCopy_api_Login(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*Login)
		out := out.(*Login)
		out.TypeMeta = in.TypeMeta
		if err := prototype.DeepCopy_api_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, c); err != nil {
			return err
		}
		if err := DeepCopy_api_LoginSpec(&in.Spec, &out.Spec, c); err != nil {
			return err
		}
		return nil
	}
}
