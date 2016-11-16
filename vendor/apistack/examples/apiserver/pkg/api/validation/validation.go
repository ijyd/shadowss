package validation

import (
	"apistack/examples/apiserver/pkg/api"
	apivalidation "gofreezer/pkg/api/validation"

	"gofreezer/pkg/util/validation/field"
)

func ValidateLogin(login *api.Login) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMeta(&login.ObjectMeta, false, apivalidation.NameIsDNSSubdomain, field.NewPath("metadata"))

	fldPath := field.NewPath("spec")

	if !(len(login.Spec.AuthName) > 0) {
		allErrs = append(allErrs, field.Required(fldPath.Child("authName"), ""))
	}

	if login.Spec.AuthName != login.Name {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("authName"), login.Spec.AuthName, "authName and metadta.name must be equal"))
	}

	if !(len(login.Spec.Auth) > 0) {
		allErrs = append(allErrs, field.Required(fldPath.Child("auth"), ""))
	}

	return allErrs
}

func ValidateLoginUpdate(login, oldlogin *api.Login) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMetaUpdate(&login.ObjectMeta, &oldlogin.ObjectMeta, field.NewPath("metadata"))
	return allErrs
}

func ValidateUserToken(token *api.UserToken) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMeta(&token.ObjectMeta, false, apivalidation.NameIsDNSSubdomain, field.NewPath("metadata"))

	fldPath := field.NewPath("spec")

	if !(len(token.Spec.Token) > 0) {
		allErrs = append(allErrs, field.Required(fldPath.Child("token"), ""))
	}
	return allErrs
}

func ValidateUserTokenUpdate(login, oldlogin *api.UserToken) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMetaUpdate(&login.ObjectMeta, &oldlogin.ObjectMeta, field.NewPath("metadata"))
	return allErrs
}

func ValidateUser(user *api.User) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMeta(&user.ObjectMeta, false, apivalidation.NameIsDNSSubdomain, field.NewPath("metadata"))

	fldPath := field.NewPath("detailInfo")
	if !(len(user.Spec.DetailInfo.Email) > 0) {
		allErrs = append(allErrs, field.Required(fldPath.Child("email"), ""))
	}

	if !(len(user.Spec.DetailInfo.Name) > 0) {
		allErrs = append(allErrs, field.Required(fldPath.Child("name"), ""))
	}

	if !(len(user.Spec.DetailInfo.Passwd) > 0) {
		allErrs = append(allErrs, field.Required(fldPath.Child("passwd"), ""))
	}

	if !(len(user.Spec.DetailInfo.ManagePasswd) > 0) {
		allErrs = append(allErrs, field.Required(fldPath.Child("managePasswd"), ""))
	}

	if user.Spec.DetailInfo.Name != user.Name {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("name"), user.Spec.DetailInfo.Name, "metadta.name and detailinfo.name must be equal"))
	}

	return allErrs
}

func ValidateUserUpdate(user, olduser *api.User) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMetaUpdate(&user.ObjectMeta, &olduser.ObjectMeta, field.NewPath("metadata"))
	return allErrs
}
