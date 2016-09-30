package vps

import (
	"encoding/json"
	"fmt"
	"strconv"

	"cloud-keeper/pkg/api"
	apierr "cloud-keeper/pkg/api/errors"
	"cloud-keeper/pkg/api/validation"
	. "cloud-keeper/pkg/api/vps/common"

	restful "github.com/emicklei/go-restful"

	"github.com/golang/glog"
)

func NewToken(login *api.Login) error {

	if len(login.Spec.AuthName) == 0 || len(login.Spec.Auth) == 0 {
		return fmt.Errorf("invalid request")
	}

	glog.V(5).Infof("Got Users by name %v", login.Spec.AuthName)
	user, err := Storage.GetUserByName(login.Spec.AuthName)
	if err != nil {
		return err
	}
	glog.V(5).Infof("Got Users %+v", user)

	if login.Spec.Auth != user.ManagePasswd {
		return fmt.Errorf("auth failure")
	}

	token, err := AddToken(user)
	login.Spec.Token = token
	login.Spec.AuthID = strconv.FormatInt(user.ID, 10)
	//empty user password
	login.Spec.Auth = string("")

	return err
}

func PostLogin(request *restful.Request, response *restful.Response) {
	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")

	statusCode := 200
	output := apierr.NewSuccess().Encode()

	defer func() {
		w.WriteHeader(statusCode)
		w.Write(output)
	}()

	login := new(api.Login)
	err := request.ReadEntity(login)
	if err != nil {
		glog.Errorf("invalid request body:%v", err)
		newErr := apierr.NewBadRequestError("request body invalid")
		output = EncodeError(newErr)
		statusCode = 400
		return
	}

	err = validation.ValidateLogin(*login)
	if err != nil {
		newErr := apierr.NewBadRequestError(err.Error())
		output = EncodeError(newErr)
		statusCode = 400
		return
	}

	glog.Infof("Got Post logins:%+v\n", login)
	err = NewToken(login)
	if err != nil {
		newErr := apierr.NewBadRequestError(err.Error())
		output = EncodeError(newErr)
		statusCode = 400
	} else {
		glog.V(5).Infof("Got login len %+v", len(login.Spec.Auth))
		output, err = json.Marshal(login)
		statusCode = 200
	}

	return

}
