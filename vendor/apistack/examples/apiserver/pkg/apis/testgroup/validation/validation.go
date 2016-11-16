package validation

import (
	"apistack/examples/apiserver/pkg/apis/testgroup"

	apivalidation "gofreezer/pkg/api/validation"
	"gofreezer/pkg/util/validation/field"
)

// ValidateTestType validates a StorageClass.
func ValidateTestType(test *testgroup.TestType) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMeta(&test.ObjectMeta, false, apivalidation.NameIsDNSSubdomain, field.NewPath("metadata"))
	return allErrs
}

// ValidateTestTypeUpdate tests if an update to StorageClass is valid.
func ValidateTestTypeUpdate(test, oldtest *testgroup.TestType) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMetaUpdate(&test.ObjectMeta, &oldtest.ObjectMeta, field.NewPath("metadata"))
	return allErrs
}
