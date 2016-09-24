package vps

import (
	"encoding/json"
	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/api/unversioned"

	"cloud-keeper/pkg/api"
	apierr "cloud-keeper/pkg/api/errors"
	"cloud-keeper/pkg/api/validation"
	"cloud-keeper/pkg/pagination"

	"github.com/golang/glog"

	restful "github.com/emicklei/go-restful"
)

func getServers(page pagination.Pager) ([]byte, int) {
	statusCode := 500
	var output []byte
	var err error

	var apisrvList api.APIServerList
	servers, err := Storage.GetAPIServer(page)
	glog.V(5).Infof("Get servers %v \r\n", servers)
	if err != nil {
		if isNotfoundErr(err) == true {
			apisrvList = api.APIServerList{
				TypeMeta: unversioned.TypeMeta{
					Kind:       "APIServerList",
					APIVersion: "v1",
				},
				ListMeta: unversioned.ListMeta{
					SelfLink: "/api/v1/apiservers",
				},
			}
		} else {
			glog.Errorf("Get apiserver failure %v \r\n", err)
			newErr := apierr.NewInternalError(err.Error())
			output = encodeError(newErr)

			return output, statusCode
		}

	} else {

		var apiservers []api.APIServer
		for _, v := range servers {
			item := api.APIServer{
				TypeMeta: unversioned.TypeMeta{
					Kind:       "APIServer",
					APIVersion: "v1",
				},
				ObjectMeta: prototype.ObjectMeta{
					Name: v.Name,
				},
				Spec: api.APIServerSpec{
					Server: api.APIServerInfor{
						ID:     v.ID,
						Host:   v.Host,
						Port:   v.Port,
						Status: v.Status,
					},
				},
			}
			apiservers = append(apiservers, item)
		}

		apisrvList = api.APIServerList{
			TypeMeta: unversioned.TypeMeta{
				Kind:       "APIServerList",
				APIVersion: "v1",
			},
			ListMeta: unversioned.ListMeta{
				SelfLink: "/api/v1/apiservers",
			},
			Items: apiservers,
		}

	}

	output, err = json.Marshal(apisrvList)
	if err != nil {
		glog.Errorln("Marshal router err", err)
		newErr := apierr.NewInternalError("marshal apiserver list resource failure")
		output = encodeError(newErr)
		statusCode = 500

	} else {
		statusCode = 200
	}
	return output, statusCode
}

//GetAPIServers ... get router list
func GetAPIServers(request *restful.Request, response *restful.Response) {

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
		glog.Errorln("Unauth request ", err)
		newErr := apierr.NewBadRequestError("invalid pagination")
		output = encodeError(newErr)
		statusCode = 400
		return
	}

	output, statusCode = getServers(page)
	baseLink := request.SelectedRoutePath()
	api.SetPageLink(baseLink, response, page)

}

func PostAPIServer(request *restful.Request, response *restful.Response) {
	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")

	statusCode := 200
	output := apierr.NewSuccess().Encode()

	defer func() {
		w.WriteHeader(statusCode)
		w.Write(output)
	}()

	server := new(api.APIServer)
	err := request.ReadEntity(server)
	if err != nil {
		glog.Errorf("invalid request body:%v", err)
		newErr := apierr.NewBadRequestError("request body invalid")
		output = encodeError(newErr)
		statusCode = 400
		return
	}

	err = validation.ValidateAPIServer(*server)
	if err != nil {
		newErr := apierr.NewBadRequestError(err.Error())
		output = encodeError(newErr)
		statusCode = 400
		return
	}

	glog.Infof("Got Post api server:%+v\n", server)
	err = Storage.CreateAPIServer(server.Spec.Server)
	if err != nil {
		newErr := apierr.NewBadRequestError(err.Error())
		output = encodeError(newErr)
		statusCode = 400
	} else {
		output, err = json.Marshal(server)
		statusCode = 200
	}

	return
}

func DeleteAPIServer(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")

	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	statusCode := 200
	var output []byte

	err := Storage.DeleteAPIServerByName(name)
	if err == nil {
		output = apierr.NewSuccess().Encode()
		statusCode = 200
	} else {
		newErr := apierr.NewNotFound("invalid request name", name)
		output = encodeError(newErr)
		statusCode = 404
	}

	w.WriteHeader(statusCode)
	w.Write(output)
}
