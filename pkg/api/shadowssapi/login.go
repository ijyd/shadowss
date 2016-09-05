package shadowssapi

import (
	"encoding/json"
	"fmt"

	"shadowsocks-go/pkg/api"
	apierr "shadowsocks-go/pkg/api/errors"

	restful "github.com/emicklei/go-restful"

	"github.com/golang/glog"
)

const (
	maxLoginCacheSize = 32 * 4
)

func NewToken(login *api.Login) error {

	if len(login.Spec.AuthName) == 0 || len(login.Spec.Auth) == 0 {
		return fmt.Errorf("invalid request")
	}

	glog.V(5).Infof("Got Users by db")
	user, err := Storage.GetUserByName(login.Spec.AuthName)
	if err != nil {
		return err
	}
	glog.V(5).Infof("Got Users %+v", user)

	if login.Spec.Auth != user.ManagePasswd {
		return fmt.Errorf("auth failure")
	}

	token, err := addToken(user)
	login.Spec.Token = token

	return err
}

func PostLogin(request *restful.Request, response *restful.Response) {
	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	statusCode := 200

	output := apierr.NewSuccess().Encode()

	login := new(api.Login)
	err := request.ReadEntity(login)
	if err != nil {
		glog.Errorf("invalid request body:%v", err)
		newErr := apierr.NewBadRequestError("request body invalid")
		internalErr, ok := newErr.(*apierr.StatusError)
		if ok {
			output = internalErr.ErrStatus.Encode()
		} else {
			glog.Errorln("status type error")
		}
		statusCode = 400
	} else {
		glog.Infof("Got Post logins:%+v\n", login)
		err := NewToken(login)
		if err != nil {
			newErr := apierr.NewBadRequestError(err.Error())
			badReq, ok := newErr.(*apierr.StatusError)
			if ok {
				output = badReq.ErrStatus.Encode()
			} else {
				glog.Errorln("status type error")
			}
			statusCode = 400
		} else {
			output, err = json.Marshal(login)
			statusCode = 200
		}
	}

	w.WriteHeader(statusCode)
	w.Write(output)
}
