package etcdrest

import (
	apierr "cloud-keeper/pkg/api/errors"
	. "cloud-keeper/pkg/api/vps/common"
	"cloud-keeper/pkg/controller/nodectl"
	"gofreezer/pkg/runtime"

	restful "github.com/emicklei/go-restful"
	"github.com/golang/glog"
)

//GetRouters ... get router list
func GetBindingUsers(request *restful.Request, response *restful.Response) {
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

	obj, err := nodectl.GetNodeAllUsers(EtcdStorage, name)

	output, err = runtime.Encode(EtcdStorage.StorageCodec.Codecs, obj)
	if err != nil {
		newErr := apierr.NewInternalError("marshal nodes resource failure")
		output = EncodeError(newErr)
		statusCode = 500
		return
	}

	// baseLink := request.SelectedRoutePath()
	// api.SetPageLink(baseLink, response, page)
	return
}
