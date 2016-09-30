package validation

import (
	"fmt"
	"gofreezer/pkg/api/prototype"

	"cloud-keeper/pkg/api"
)

// ValidateObjMeta requires that api.ObjectMeta has a name
func ValidateObjMeta(meta prototype.ObjectMeta) error {

	if !(len(meta.Name) > 0) {
		return fmt.Errorf("invalid obj name")
	}

	return nil
}

// ValidateAccServer requires that api.ObjectMeta has a Label with key and expectedValue
func ValidateAccServer(server api.AccServer) error {

	if !(len(server.Region) > 0) {
		return fmt.Errorf("invalid region")
	}

	if !(len(server.Size) > 0) {
		return fmt.Errorf("invalid size")
	}

	return nil
}

// ValidateAPIServer requires that api.ObjectMeta has a Label with key and expectedValue
func ValidateAPIServer(server api.APIServer) error {

	err := ValidateObjMeta(server.ObjectMeta)
	if err != nil {
		return fmt.Errorf("invalid ObjectMeta")
	}

	if server.Spec.Server.Port == 0 {
		return fmt.Errorf("invalid port")
	}

	if !(len(server.Spec.Server.Host) > 0) {
		return fmt.Errorf("invalid host")
	}

	return nil
}

func ValidateAccount(acc api.Account) error {
	err := ValidateObjMeta(acc.ObjectMeta)
	if err != nil {
		return fmt.Errorf("invalid ObjectMeta")
	}

	if acc.Spec.AccDetail.ExpireTime.IsZero() {
		return fmt.Errorf("invalid ExpireTime")
	}

	if !(len(acc.Spec.AccDetail.Name) > 0) {
		return fmt.Errorf("invalid Name")
	}

	if !(len(acc.Spec.AccDetail.Key) > 0) {
		return fmt.Errorf("invalid Key")
	}

	if !(len(acc.Spec.AccDetail.Descryption) > 0) {
		return fmt.Errorf("invalid Descryption")
	}

	if !(len(acc.Spec.AccDetail.Lables) > 0) {
		return fmt.Errorf("invalid Lables")
	}

	if acc.Spec.AccDetail.CreditCeilings == 0 {
		return fmt.Errorf("invalid CreditCeilings")
	}

	return nil
}

func ValidateLogin(login api.Login) error {
	err := ValidateObjMeta(login.ObjectMeta)
	if err != nil {
		return fmt.Errorf("invalid ObjectMeta")
	}

	if !(len(login.Spec.AuthName) > 0) {
		return fmt.Errorf("invalid authname")
	}

	if login.Spec.AuthName != login.Name {
		return fmt.Errorf("not a same user")
	}

	if !(len(login.Spec.Auth) > 0) {
		return fmt.Errorf("invalid auth")
	}

	return nil
}

func ValidateNode(node api.Node) error {
	err := ValidateObjMeta(node.ObjectMeta)
	if err != nil {
		return fmt.Errorf("invalid ObjectMeta")
	}

	if !(len(node.Spec.Server.Host) > 0) {
		return fmt.Errorf("invalid authname")
	}

	if !(len(node.Spec.Server.Name) > 0) {
		return fmt.Errorf("invalid auth")
	}

	return nil
}

func ValidateUser(user api.User) error {
	err := ValidateObjMeta(user.ObjectMeta)
	if err != nil {
		return fmt.Errorf("invalid ObjectMeta")
	}

	if !(len(user.Spec.DetailInfo.Email) > 0) {
		return fmt.Errorf("invalid email")
	}

	if !(len(user.Spec.DetailInfo.Name) > 0) {
		return fmt.Errorf("invalid name")
	}

	if !(len(user.Spec.DetailInfo.Passwd) > 0) {
		return fmt.Errorf("invalid passwd")
	}

	if !(len(user.Spec.DetailInfo.ManagePasswd) > 0) {
		return fmt.Errorf("invalid manage passwd")
	}

	if user.Spec.DetailInfo.Name != user.Name {
		return fmt.Errorf("not a same name")
	}

	return nil
}

func ValidateUserReference(refer api.UserReferences) error {
	if !(len(refer.Name) > 0) {
		return fmt.Errorf("invalid name")
	}

	if refer.ID == 0 {
		return fmt.Errorf("invalid user id")
	}

	if !(len(refer.Password) > 0) {
		return fmt.Errorf("invalid passwd")
	}

	if !(len(refer.Method) > 0) {
		return fmt.Errorf("invalid method")
	}

	return nil
}
