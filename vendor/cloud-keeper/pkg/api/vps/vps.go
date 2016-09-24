package vps

import (
	"cloud-keeper/pkg/api"
	apierr "cloud-keeper/pkg/api/errors"
	"cloud-keeper/pkg/api/validation"
	"cloud-keeper/pkg/backend"
	"cloud-keeper/pkg/collector"
	"cloud-keeper/pkg/collector/collectorbackend"
	"cloud-keeper/pkg/collector/collectorbackend/factory"
	"cloud-keeper/pkg/etcdhelper"
	"fmt"
	"strconv"

	restful "github.com/emicklei/go-restful"

	"github.com/golang/glog"
)

var Storage *backend.Backend
var EtcdStorage *etcdhelper.EtcdHelper

var accountCollector = make(map[string]collector.Collector, 2)

func updateCollector(v *api.AccountDetail) collector.Collector {

	cfg := collectorbackend.Config{
		Type:   string(v.Operators),
		APIKey: v.Key,
	}
	collectorHandle, err := factory.Create(cfg)
	if err != nil {
		glog.Errorf("create collector failure %v \r\n", err)
	} else {
		accountCollector[v.Name] = collectorHandle
	}
	return collectorHandle
}

func getCollector(name string) (collector.Collector, error) {
	collectorHandle, ok := accountCollector[name]
	if !ok {
		acc, err := Storage.GetAccountByname(name)
		if err != nil {
			return nil, fmt.Errorf("not found vps")
		} else {
			collectorHandle = updateCollector(acc)
			if collectorHandle == nil {
				return nil, fmt.Errorf("not found vps")
			}
		}
	}

	return collectorHandle, nil

}

func GetAccountInfo(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")

	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	encoded := request.Request.Header.Get("Authorization")
	statusCode := 200
	var output []byte

	user, err := CheckToken(encoded)
	if err != nil || user == nil {
		glog.Errorln("Unauth request ", err)
		newErr := apierr.NewUnauthorized("invalid token")
		output = encodeError(newErr)
		statusCode = 401
	} else {
		collectorHandle, err := getCollector(name)
		if err == nil {
			output, err = collectorHandle.GetAccount()
			if err != nil {
				newErr := apierr.NewInternalError(err.Error())
				output = encodeError(newErr)
				statusCode = 401
			}
		} else {
			newErr := apierr.NewNotFound("not support vps", name)
			output = encodeError(newErr)
			statusCode = 404
		}

	}

	w.WriteHeader(statusCode)
	w.Write(output)
}

func GetAccServers(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")

	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	encoded := request.Request.Header.Get("Authorization")
	statusCode := 200
	var output []byte

	defer func() {
		w.WriteHeader(statusCode)
		w.Write(output)
	}()

	user, err := CheckToken(encoded)
	if err != nil || user == nil {
		glog.Errorln("Unauth request ", err)
		newErr := apierr.NewUnauthorized("invalid token")
		output = encodeError(newErr)
		statusCode = 401
		return
	}

	page, err := api.PageParse(request)
	if err != nil {
		newErr := apierr.NewBadRequestError("invalid pagination")
		output = encodeError(newErr)
		statusCode = 400
		return
	}

	collectorHandle, err := getCollector(name)
	if err == nil {
		output, err = collectorHandle.GetServers(page)
		if err != nil {
			newErr := apierr.NewInternalError(err.Error())
			output = encodeError(newErr)
			statusCode = 500
		}
	} else {
		newErr := apierr.NewNotFound("not support vps", name)
		output = encodeError(newErr)
		statusCode = 404
	}

	baseLink := request.SelectedRoutePath()
	api.SetPageLink(baseLink, response, page)

	return
}

func PostAccServer(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")
	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")

	statusCode := 200
	output := apierr.NewSuccess().Encode()

	defer func() {
		w.WriteHeader(statusCode)
		w.Write(output)
	}()

	acc := new(api.AccServer)
	err := request.ReadEntity(acc)

	if err != nil {
		newErr := apierr.NewBadRequestError("request body invalid")
		output = encodeError(newErr)
		statusCode = 400
		return
	}

	glog.V(5).Infof("check request ")
	err = validation.ValidateAccServer(*acc)
	if err != nil {
		newErr := apierr.NewBadRequestError(err.Error())
		output = encodeError(newErr)
		statusCode = 400
		return
	}

	collectorHandle, err := getCollector(name)
	if err != nil {
		newErr := apierr.NewNotFound(err.Error(), name)
		output = encodeError(newErr)
		statusCode = 404
		return
	}

	err = collectorHandle.CreateServer(acc)
	if err != nil {
		newErr := apierr.NewInternalError(err.Error())
		output = encodeError(newErr)
		statusCode = 401
	} else {
		output = apierr.NewSuccess().Encode()
		statusCode = 200
	}

	return
}

func DeleteAccServer(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")
	idStr := request.PathParameter("id")

	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	statusCode := 200
	var output []byte

	id, err := strconv.Atoi(idStr)
	if err != nil {
		newErr := apierr.NewInternalError(err.Error())
		output = encodeError(newErr)
		statusCode = 500
	} else {

		collectorHandle, err := getCollector(name)
		if err == nil {
			err := collectorHandle.DeleteServer(int64(id))
			if err != nil {
				newErr := apierr.NewInternalError(err.Error())
				output = encodeError(newErr)
				statusCode = 500
			} else {
				output = apierr.NewSuccess().Encode()
				statusCode = 200
			}
		} else {
			newErr := apierr.NewNotFound("not support vps", name)
			output = encodeError(newErr)
			statusCode = 404
		}
	}

	w.WriteHeader(statusCode)
	w.Write(output)
}

func PostAccServerExec(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")
	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")

	statusCode := 200
	output := apierr.NewSuccess().Encode()

	defer func() {
		w.WriteHeader(statusCode)
		w.Write(output)
	}()

	command := new(api.AccServerCommand)
	err := request.ReadEntity(command)

	if err != nil {
		newErr := apierr.NewBadRequestError("request body invalid")
		output = encodeError(newErr)
		statusCode = 400
		return
	}

	collectorHandle, err := getCollector(name)
	if err != nil {
		newErr := apierr.NewNotFound(err.Error(), name)
		output = encodeError(newErr)
		statusCode = 404
		return
	}

	err = collectorHandle.Exec(command)
	if err != nil {
		newErr := apierr.NewInternalError(err.Error())
		output = encodeError(newErr)
		statusCode = 500
	} else {
		output = apierr.NewSuccess().Encode()
		statusCode = 200
	}

	return
}

func GetAccServersSSHkey(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")

	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	encoded := request.Request.Header.Get("Authorization")
	statusCode := 200
	var output []byte

	defer func() {
		w.WriteHeader(statusCode)
		w.Write(output)
	}()

	user, err := CheckToken(encoded)
	if err != nil || user == nil {
		glog.Errorln("Unauth request ", err)
		newErr := apierr.NewUnauthorized("invalid token")
		output = encodeError(newErr)
		statusCode = 401
		return
	}

	collectorHandle, err := getCollector(name)
	if err == nil {
		output, err = collectorHandle.GetSSHKey()
		if err != nil {
			newErr := apierr.NewInternalError(err.Error())
			output = encodeError(newErr)
			statusCode = 500
		}
	} else {
		newErr := apierr.NewNotFound("not support vps", name)
		output = encodeError(newErr)
		statusCode = 404
	}

}
