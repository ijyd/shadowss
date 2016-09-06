package shadowssapi

import (
	"encoding/json"
	"shadowsocks-go/pkg/api"
	apierr "shadowsocks-go/pkg/api/errors"
	"shadowsocks-go/pkg/backend/db"

	restful "github.com/emicklei/go-restful"
	"github.com/golang/glog"
)

func getServers(user *db.User) ([]byte, int) {
	statusCode := 500
	var output []byte
	var err error

	servers, err := Storage.GetAPIServer()
	if err != nil {
		glog.Errorf("Get apiserver failure %v \r\n", err)
		newErr := apierr.NewInternalError("marshal router resource failure")
		internalErr, _ := newErr.(*apierr.StatusError)

		output = internalErr.ErrStatus.Encode()
		return output, statusCode
	}

	var apiserver []api.APIServerInfor
	for _, v := range servers {
		apisrvInfo := api.APIServerInfor{
			Host: v.Host,
			Port: v.Port,
		}
		apiserver = append(apiserver, apisrvInfo)
	}

	apiServers := api.APIServer{
		TypeMeta: api.TypeMeta{
			Kind:       "ShadowAPIServer",
			APIVersion: "v1",
		},
		Spec: api.APIServerSpec{
			Server: apiserver,
		},
	}

	output, err = json.Marshal(apiServers)
	if err != nil {
		glog.Errorln("Marshal router err", err)
		newErr := apierr.NewInternalError("marshal router resource failure")
		internalErr, _ := newErr.(*apierr.StatusError)

		output = internalErr.ErrStatus.Encode()

	} else {
		statusCode = 200
	}
	return output, statusCode
}

//GetRouters ... get router list
func GetAPIServers(request *restful.Request, response *restful.Response) {

	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	encoded := request.Request.Header.Get("Authorization")
	statusCode := 200
	var output []byte

	user, err := CheckToken(encoded)
	if err != nil || user == nil {
		glog.Errorln("Unauth request ", err)
		newErr := apierr.NewUnauthorized("invalid token")
		internalErr, ok := newErr.(*apierr.StatusError)
		if ok {
			output = internalErr.ErrStatus.Encode()
		} else {
			glog.Errorln("status type error")
		}
		statusCode = 401
	} else {
		output, statusCode = getServers(user)
	}

	w.WriteHeader(statusCode)
	w.Write(output)
}

func PostAPIServer(request *restful.Request, response *restful.Response) {
	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	statusCode := 200

	output := apierr.NewSuccess().Encode()

	server := new(api.APIServer)
	err := request.ReadEntity(server)
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
		glog.Infof("Got Post api server:%+v\n", server)
		err := Storage.CreateAPIServer(server.Spec.Server[0].Host, server.Spec.Server[0].Port, true)
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
			output, err = json.Marshal(server)
			statusCode = 200
		}
	}

	w.WriteHeader(statusCode)
	w.Write(output)
}
