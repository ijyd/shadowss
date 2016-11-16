package apiext

import (
	"gofreezer/pkg/conversion"
	"gofreezer/pkg/runtime"

	"github.com/golang/glog"
)

// This is a "fast-path" that avoids reflection for common types. It focuses on the objects that are
// converted the most in the cluster.
// TODO: generate one of these for every external API group - this is to prove the impact
func addFastPathConversionFuncs(scheme *runtime.Scheme) error {
	scheme.AddGenericConversionFunc(func(objA, objB interface{}, s conversion.Scope) (bool, error) {
		glog.V(5).Infof("generic coversions %v to %v", objA, objB)
		//debug.PrintStack()
		switch a := objA.(type) {
		case []byte:
			switch b := objB.(type) {
			case *Login:
				return true, Convert_DBResult_To_v1_Login(a, b, s)
			}
		}
		return false, nil
	})
	return nil
}

func Convert_DBResult_To_v1_Login(in []byte, out *Login, s conversion.Scope) error {
	glog.Infof("call byte login v1")
	return nil
}

func addConversionFuncs(scheme *runtime.Scheme) error {
	//return scheme.AddConversionFuncs()
	return nil
}
