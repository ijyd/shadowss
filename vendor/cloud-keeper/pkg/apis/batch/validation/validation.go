package validation

import (
	"cloud-keeper/pkg/apis/batch"

	apivalidation "gofreezer/pkg/api/validation"
	"gofreezer/pkg/util/validation/field"
)

// ValidateBatchAccServer validates a StorageClass.
func ValidateBatchAccServer(test *batch.BatchAccServer) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMeta(&test.ObjectMeta, false, apivalidation.NameIsDNSSubdomain, field.NewPath("metadata"))
	return allErrs
}

// ValidateBatchAccServerUpdate tests if an update to StorageClass is valid.
func ValidateBatchAccServerUpdate(accsrv, oldaccsrv *batch.BatchAccServer) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMetaUpdate(&accsrv.ObjectMeta, &oldaccsrv.ObjectMeta, field.NewPath("metadata"))
	return allErrs
}
