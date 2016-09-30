package vps

import (
	"encoding/json"
	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/api/unversioned"
	"gofreezer/pkg/runtime"

	"cloud-keeper/pkg/api"
	apierr "cloud-keeper/pkg/api/errors"
	"cloud-keeper/pkg/api/validation"
	. "cloud-keeper/pkg/api/vps/common"
	"cloud-keeper/pkg/controller/nodectl"
	"cloud-keeper/pkg/pagination"

	"github.com/golang/glog"

	restful "github.com/emicklei/go-restful"
)

func getNodes(page pagination.Pager) ([]byte, int) {
	statusCode := 500
	var output []byte
	var err error

	var nodeList api.NodeList
	nodes, err := Storage.GetNodes(page)
	if err != nil {
		if IsNotfoundErr(err) == true {
			nodeList = api.NodeList{
				TypeMeta: unversioned.TypeMeta{
					Kind:       "NodeList",
					APIVersion: "v1",
				},
				ListMeta: unversioned.ListMeta{
					SelfLink: "/api/v1/nodes",
				},
			}
		} else {
			newErr := apierr.NewInternalError("marshal nodes resource failure")
			internalErr, _ := newErr.(*apierr.StatusError)

			output = internalErr.ErrStatus.Encode()
			return output, statusCode
		}
	} else {
		var items []api.Node
		for _, v := range nodes {
			item := api.Node{
				TypeMeta: unversioned.TypeMeta{
					Kind:       "Node",
					APIVersion: "v1",
				},
				ObjectMeta: prototype.ObjectMeta{
					Name: v.Name,
				},
				Spec: api.NodeSpec{
					Server: api.NodeServer{
						ID:            v.ID,
						Host:          v.Host,
						Name:          v.Name,
						Status:        v.Status,
						EnableOTA:     v.EnableOTA,
						Method:        v.Method,
						TrafficRate:   v.TrafficRate,
						Description:   v.Description,
						TrafficLimit:  v.TrafficLimit,
						Upload:        v.Upload,
						Download:      v.Download,
						Location:      v.Location,
						AccServerID:   v.AccServerID,
						AccServerName: v.AccServerName,
					},
				},
			}
			items = append(items, item)
		}

		nodeList = api.NodeList{
			TypeMeta: unversioned.TypeMeta{
				Kind:       "NodeList",
				APIVersion: "v1",
			},
			ListMeta: unversioned.ListMeta{
				SelfLink: "/api/v1/nodes",
			},
			Items: items,
		}
	}

	output, err = json.Marshal(nodeList)
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
func GetNodes(request *restful.Request, response *restful.Response) {

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
		glog.Errorln("Unauth request %v", err)
		newErr := apierr.NewUnauthorized("invalid token")
		internalErr, ok := newErr.(*apierr.StatusError)
		if ok {
			output = internalErr.ErrStatus.Encode()
		} else {
			glog.Errorln("status type error")
		}
		statusCode = 401
	}

	page, err := api.PageParse(request)
	if err != nil {
		glog.Errorln("Unauth request ", err)
		newErr := apierr.NewBadRequestError("invalid pagination")
		output = EncodeError(newErr)
		statusCode = 400
		return
	}

	output, statusCode = getNodes(page)
	baseLink := request.SelectedRoutePath()
	api.SetPageLink(baseLink, response, page)
	return
}

func PostNode(request *restful.Request, response *restful.Response) {
	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")

	statusCode := 200
	output := apierr.NewSuccess().Encode()

	defer func() {
		w.WriteHeader(statusCode)
		w.Write(output)
	}()

	item := new(api.Node)
	err := request.ReadEntity(item)
	if err != nil {
		glog.Errorf("invalid request body:%v", err)
		newErr := apierr.NewBadRequestError("request body invalid")
		output = EncodeError(newErr)
		statusCode = 400
		return
	}

	err = validation.ValidateNode(*item)
	if err != nil {
		newErr := apierr.NewBadRequestError(err.Error())
		output = EncodeError(newErr)
		statusCode = 400
		return
	}

	outItem, err := nodectl.AddNode(Storage, EtcdStorage, item, 0, true, false)
	if err != nil {
		newErr := apierr.NewBadRequestError(err.Error())
		output = EncodeError(newErr)
		statusCode = 400
		return
	}

	output, err = runtime.Encode(EtcdStorage.StorageCodec.Codecs, outItem)
	if err != nil {
		newErr := apierr.NewInternalError("marshal nodes resource failure")
		output = EncodeError(newErr)
		statusCode = 500
		return
	}

	return
}

func DeleteNode(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")

	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	statusCode := 200
	var output []byte

	defer func() {
		w.WriteHeader(statusCode)
		w.Write(output)
	}()

	err := nodectl.DelNode(Storage, EtcdStorage, name, true, true)
	if err != nil {
		newErr := apierr.NewNotFound("invalid request name", name)
		output = EncodeError(newErr)
		statusCode = 404
		return
	}

	output = apierr.NewSuccess().Encode()
	statusCode = 200
	return
}

func PutNode(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")

	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	statusCode := 200
	var output []byte

	defer func() {
		w.WriteHeader(statusCode)
		w.Write(output)
	}()

	item := new(api.Node)
	err := request.ReadEntity(item)
	if err != nil {
		glog.Errorf("invalid request body:%v", err)
		newErr := apierr.NewBadRequestError("request body invalid")
		output = EncodeError(newErr)
		statusCode = 400
		return
	}

	err = validation.ValidateNode(*item)
	if err != nil {
		newErr := apierr.NewBadRequestError(err.Error())
		output = EncodeError(newErr)
		statusCode = 400
		return
	}

	outItem, err := nodectl.UpdateNode(Storage, EtcdStorage, item, true, true)
	if err != nil {
		newErr := apierr.NewNotFound("invalid request name", name)
		output = EncodeError(newErr)
		statusCode = 404
		return
	}

	output, err = runtime.Encode(EtcdStorage.StorageCodec.Codecs, outItem)
	if err != nil {
		newErr := apierr.NewInternalError("marshal nodes resource failure")
		output = EncodeError(newErr)
		statusCode = 500
		return
	}

	return

}
