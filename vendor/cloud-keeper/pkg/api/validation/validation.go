package validation

import (
	apivalidation "gofreezer/pkg/api/validation"
	"gofreezer/pkg/util/validation/field"

	"cloud-keeper/pkg/api"
)

// ValidateAccServer requires that api.ObjectMeta has a Label with key and expectedValue
func ValidateAccServer(server *api.AccServer) field.ErrorList {

	allErrs := apivalidation.ValidateObjectMeta(&server.ObjectMeta, false, apivalidation.NameIsDNSSubdomain, field.NewPath("metadata"))
	fldPath := field.NewPath("spec")

	if !(len(server.Spec.Region) > 0) {
		allErrs = append(allErrs, field.Required(fldPath.Child("region"), ""))
	}

	if !(len(server.Spec.Size) > 0) {
		allErrs = append(allErrs, field.Required(fldPath.Child("Size"), ""))
	}

	return allErrs
}

func ValidateAccServerUpdate(srv, oldsrv *api.AccServer) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMetaUpdate(&srv.ObjectMeta, &oldsrv.ObjectMeta, field.NewPath("metadata"))
	return allErrs
}

// ValidateUserPublicFile requires that api.ObjectMeta has a Label with key and expectedValue
func ValidateUserPublicFile(server *api.UserPublicFile) field.ErrorList {

	allErrs := apivalidation.ValidateObjectMeta(&server.ObjectMeta, false, apivalidation.NameIsDNSSubdomain, field.NewPath("metadata"))

	//fldPath := field.NewPath("spec")

	return allErrs
}

func ValidateUserPublicFileUpdate(file, oldfile *api.UserPublicFile) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMetaUpdate(&file.ObjectMeta, &oldfile.ObjectMeta, field.NewPath("metadata"))
	return allErrs
}

// ValidateAPIServer requires that api.ObjectMeta has a Label with key and expectedValue
func ValidateAPIServer(server *api.APIServer) field.ErrorList {

	allErrs := apivalidation.ValidateObjectMeta(&server.ObjectMeta, false, apivalidation.NameIsDNSSubdomain, field.NewPath("metadata"))

	fldPath := field.NewPath("spec")

	if server.Spec.Server.Port == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("port"), ""))
	}

	if !(len(server.Spec.Server.Host) > 0) {
		allErrs = append(allErrs, field.Required(fldPath.Child("host"), ""))
	}

	return allErrs
}

func ValidateAPIServerUpdate(node, oldnode *api.APIServer) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMetaUpdate(&node.ObjectMeta, &oldnode.ObjectMeta, field.NewPath("metadata"))
	return allErrs
}

func ValidateNodeUser(node *api.NodeUser) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMeta(&node.ObjectMeta, false, apivalidation.NameIsDNSSubdomain, field.NewPath("metadata"))
	return allErrs
}

func ValidateNodeUserUpdate(node, oldnode *api.NodeUser) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMetaUpdate(&node.ObjectMeta, &oldnode.ObjectMeta, field.NewPath("metadata"))
	return allErrs
}

func ValidateNode(node *api.Node) field.ErrorList {

	allErrs := apivalidation.ValidateObjectMeta(&node.ObjectMeta, false, apivalidation.NameIsDNSSubdomain, field.NewPath("metadata"))

	fldPath := field.NewPath("spec")

	if !(len(node.Spec.Server.Host) > 0) {
		allErrs = append(allErrs, field.Required(fldPath.Child("host"), ""))
	}

	if !(len(node.Spec.Server.Name) > 0) {
		allErrs = append(allErrs, field.Required(fldPath.Child("name"), ""))
	}
	return allErrs
}

func ValidateNodeUpdate(node, oldnode *api.Node) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMetaUpdate(&node.ObjectMeta, &oldnode.ObjectMeta, field.NewPath("metadata"))
	return allErrs
}

func ValidateAccount(acc *api.Account) field.ErrorList {

	allErrs := apivalidation.ValidateObjectMeta(&acc.ObjectMeta, false, apivalidation.NameIsDNSSubdomain, field.NewPath("metadata"))

	fldPath := field.NewPath("spec")

	if acc.Spec.AccDetail.ExpireTime.IsZero() {
		allErrs = append(allErrs, field.Required(fldPath.Child("expireTiem"), ""))
	}

	if !(len(acc.Spec.AccDetail.Name) > 0) {
		allErrs = append(allErrs, field.Required(fldPath.Child("name"), ""))
	}

	if !(len(acc.Spec.AccDetail.Key) > 0) {
		allErrs = append(allErrs, field.Required(fldPath.Child("key"), ""))
	}

	if !(len(acc.Spec.AccDetail.Descryption) > 0) {
		allErrs = append(allErrs, field.Required(fldPath.Child("description"), ""))
	}

	if !(len(acc.Spec.AccDetail.Lables) > 0) {
		allErrs = append(allErrs, field.Required(fldPath.Child("lables"), ""))
	}

	if acc.Spec.AccDetail.CreditCeilings == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("creditCeilings"), ""))
	}
	return allErrs
}

func ValidateAccountUpdate(login, oldlogin *api.Account) field.ErrorList {
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

func ValidateUser(user *api.User) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMeta(&user.ObjectMeta, false, nil, field.NewPath("metadata"))

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

func ValidateUserService(user *api.UserService) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMeta(&user.ObjectMeta, false, nil, field.NewPath("metadata"))
	// fldPath := field.NewPath("spec")
	//
	// refer := user.Spec.NodeUserReference
	// if !(len(refer.Name) > 0) {
	// 	return fmt.Errorf("invalid name")
	// }
	//
	// if refer.ID == 0 {
	// 	return fmt.Errorf("invalid user id")
	// }
	//
	// if !(len(refer.Password) > 0) {
	// 	return fmt.Errorf("invalid passwd")
	// }
	//
	// if !(len(refer.Method) > 0) {
	// 	return fmt.Errorf("invalid method")
	// }

	return allErrs
}

func ValidateUserServiceUpdate(user, olduser *api.UserService) field.ErrorList {
	allErrs := apivalidation.ValidateObjectMetaUpdate(&user.ObjectMeta, &olduser.ObjectMeta, field.NewPath("metadata"))
	return allErrs
}
