package shadowssapi

import (
	"encoding/json"
	"strconv"

	"shadowss/pkg/api"
	apierr "shadowss/pkg/api/errors"
	"shadowss/pkg/backend/db"

	restful "github.com/emicklei/go-restful"
	"github.com/golang/glog"
)

func getServers(user *db.User) ([]byte, int) {
	statusCode := 500
	var output []byte
	var err error

	var apisrvList api.APIServerList
	servers, err := Storage.GetAPIServer()
	if err != nil {
		if err.Error() == "not found" {
			apisrvList = api.APIServerList{
				TypeMeta: api.TypeMeta{
					Kind:       "APIServerList",
					APIVersion: "v1",
				},
				ListMeta: api.ListMeta{
					SelfLink: "/api/v1/apiservers",
				},
			}
		} else {
			glog.Errorf("Get apiserver failure %v \r\n", err)
			newErr := apierr.NewInternalError(err.Error())
			internalErr, _ := newErr.(*apierr.StatusError)

			output = internalErr.ErrStatus.Encode()
			return output, statusCode
		}

	} else {

		var apiservers []api.APIServer
		for _, v := range servers {
			item := api.APIServer{
				TypeMeta: api.TypeMeta{
					Kind:       "APIServer",
					APIVersion: "v1",
				},
				ObjectMeta: api.ObjectMeta{
					Name: v.Name,
				},
				Spec: api.APIServerSpec{
					Server: api.APIServerInfor{
						ID:   v.ID,
						Host: v.Host,
						Port: v.Port,
					},
				},
			}
			apiservers = append(apiservers, item)
		}

		apisrvList = api.APIServerList{
			TypeMeta: api.TypeMeta{
				Kind:       "APIServerList",
				APIVersion: "v1",
			},
			ListMeta: api.ListMeta{
				SelfLink: "/api/v1/apiservers",
			},
			Items: apiservers,
		}

	}

	output, err = json.Marshal(apisrvList)
	if err != nil {
		glog.Errorln("Marshal router err", err)
		newErr := apierr.NewInternalError("marshal apiserver list resource failure")
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
		err := Storage.CreateAPIServer(server.ObjectMeta.Name, server.Spec.Server.Host, server.Spec.Server.Port, true)
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

func DeleteAPIServer(request *restful.Request, response *restful.Response) {
	idStr := request.PathParameter("id")

	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	statusCode := 200
	var output []byte

	id, err := strconv.Atoi(idStr)
	glog.V(5).Infoln("Get api server:", id)

	err = Storage.DeleteAPIServerByID(int64(id))
	if err == nil {
		output = apierr.NewSuccess().Encode()
		statusCode = 200
	} else {
		newErr := apierr.NewNotFound("invalid request name", idStr)
		internalErr, ok := newErr.(*apierr.StatusError)
		if ok {
			output = internalErr.ErrStatus.Encode()
		} else {
			glog.Errorln("status type error")
		}
		statusCode = 404
	}

	w.WriteHeader(statusCode)
	w.Write(output)
}
