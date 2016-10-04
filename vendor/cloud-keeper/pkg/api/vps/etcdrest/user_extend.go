package etcdrest

import (
	"cloud-keeper/pkg/api"
	apierr "cloud-keeper/pkg/api/errors"
	"cloud-keeper/pkg/api/validation"
	. "cloud-keeper/pkg/api/vps/common"
	"cloud-keeper/pkg/controller"
	"cloud-keeper/pkg/controller/userctl"

	"gofreezer/pkg/runtime"

	restful "github.com/emicklei/go-restful"
	"github.com/golang/glog"
)

//GetRouters ... get router list
func GetBindingNodes(request *restful.Request, response *restful.Response) {
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

	user, err := CheckUserLevelToken(encoded)
	if err != nil || user == nil {
		glog.Errorln("Unauth request ", err)
		newErr := apierr.NewUnauthorized("invalid token")
		output = EncodeError(newErr)
		statusCode = 401
		return
	}

	// page, err := api.PageParse(request)
	// if err != nil {
	// 	glog.Errorln("Unauth request ", err)
	// 	newErr := apierr.NewBadRequestError("invalid pagination")
	// 	output = EncodeError(newErr)
	// 	statusCode = 400
	// 	return
	// }

	obj, err := userctl.GetUserService(EtcdStorage, name)

	output, err = runtime.Encode(EtcdStorage.StorageCodec.Codecs, obj)
	if err != nil {
		newErr := apierr.NewInternalError("marshal nodes resource failure")
		output = EncodeError(newErr)
		statusCode = 500
		return
	}

	// baseLink := request.SelectedRoutePath()
	// api.SetPageLink(baseLink, response, page)
}

//PutProperties ... update user properties
func PutProperties(request *restful.Request, response *restful.Response) {
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

	user, err := CheckUserLevelToken(encoded)
	if err != nil || user == nil {
		glog.Errorln("Unauth request ", err)
		newErr := apierr.NewUnauthorized("invalid token")
		output = EncodeError(newErr)
		statusCode = 401
		return
	}

	annotations := new(map[string]string)
	err = request.ReadEntity(annotations)
	if err != nil {
		glog.Errorf("invalid request body:%v", err)
		newErr := apierr.NewBadRequestError("request body invalid")
		output = EncodeError(newErr)
		statusCode = 400
		return
	}

	obj, err := userctl.UpdateUserAnnotations(EtcdStorage, name, *annotations)
	if err != nil {
		glog.Errorf("update user annotations failure %v\r\n", err)
		newErr := apierr.NewInternalError(err.Error())
		output = EncodeError(newErr)
		statusCode = 500
		return
	}

	go controller.ReallocUserNodeByProperties(name, *annotations)
	// err = controller.ReallocUserNodeByProperties(name, *annotations)
	// if err != nil {
	// 	glog.Errorf("update user node by properties failure %v\r\n", err)
	// 	newErr := apierr.NewInternalError(err.Error())
	// 	output = EncodeError(newErr)
	// 	statusCode = 500
	// 	return
	// }

	output, err = runtime.Encode(EtcdStorage.StorageCodec.Codecs, obj)
	if err != nil {
		newErr := apierr.NewInternalError("marshal nodes resource failure")
		output = EncodeError(newErr)
		statusCode = 500
		return
	}

	output = apierr.NewSuccess().Encode()
	statusCode = 200

}

func PutUserToNode(request *restful.Request, response *restful.Response) {
	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")

	statusCode := 200
	output := apierr.NewSuccess().Encode()
	encoded := request.Request.Header.Get("Authorization")

	defer func() {
		w.WriteHeader(statusCode)
		w.Write(output)
	}()

	user, err := CheckToken(encoded)
	if err != nil || user == nil {
		glog.Errorln("Unauth request ", err)
		newErr := apierr.NewUnauthorized("invalid token")
		output = EncodeError(newErr)
		statusCode = 401
		return
	}

	nodeRefer := new(map[string]api.UserReferences)
	err = request.ReadEntity(nodeRefer)
	if err != nil {
		glog.Errorf("invalid request body:%v", err)
		newErr := apierr.NewBadRequestError("request body invalid")
		output = EncodeError(newErr)
		statusCode = 400
		return
	}

	for _, userRefer := range *nodeRefer {
		err = validation.ValidateUserReference(userRefer)
		if err != nil {
			newErr := apierr.NewBadRequestError(err.Error())
			output = EncodeError(newErr)
			statusCode = 400
			return
		}
	}

	err = controller.BindUserToNode(*nodeRefer)
	if err != nil {
		newErr := apierr.NewInternalError(err.Error())
		output = EncodeError(newErr)
		statusCode = 500
		return
	}

	output = apierr.NewSuccess().Encode()
	statusCode = 200

	return
}

func DeleteUserNode(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")
	nodeName := request.PathParameter("nodename")

	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	statusCode := 200
	var output []byte
	encoded := request.Request.Header.Get("Authorization")

	defer func() {
		w.WriteHeader(statusCode)
		w.Write(output)
	}()

	user, err := CheckToken(encoded)
	if err != nil || user == nil {
		glog.Errorln("Unauth request ", err)
		newErr := apierr.NewUnauthorized("invalid token")
		output = EncodeError(newErr)
		statusCode = 401
		return
	}

	err = controller.DeleteUserServiceNode(nodeName, name)
	if err != nil {
		glog.Errorf("delete user from node error %v\r\n", err)
		newErr := apierr.NewInternalError(err.Error())
		output = EncodeError(newErr)
		statusCode = 500
	}

	output = apierr.NewSuccess().Encode()
	statusCode = 200

}
