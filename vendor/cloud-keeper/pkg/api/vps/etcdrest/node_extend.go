package etcdrest

import (
	"cloud-keeper/pkg/api"
	apierr "cloud-keeper/pkg/api/errors"
	. "cloud-keeper/pkg/api/vps/common"
	"cloud-keeper/pkg/controller/userctl"
	"cloud-keeper/pkg/pagination"

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

	page, err := api.PageParse(request)
	if err != nil {
		newErr := apierr.NewBadRequestError("invalid pagination")
		output = EncodeError(newErr)
		statusCode = 400
		return
	}

	objlist, err := userctl.GetUserServicesByNodeName(EtcdStorage, name)
	if err != nil {
		newErr := apierr.NewInternalError(err.Error())
		output = EncodeError(newErr)
		statusCode = 500
		return
	}

	for k, user := range objlist.Items {
		nodeRefer, ok := user.Spec.NodeUserReference[name]
		if ok {
			nodeUserRefer := map[string]api.NodeReferences{name: nodeRefer}
			objlist.Items[k].Spec.NodeUserReference = nodeUserRefer
		}
	}

	obj, err := userlistPage(objlist, page)
	if err != nil {
		newErr := apierr.NewInternalError(err.Error())
		output = EncodeError(newErr)
		statusCode = 500
		return
	}

	output, err = runtime.Encode(EtcdStorage.StorageCodec.Codecs, obj)
	if err != nil {
		newErr := apierr.NewInternalError("marshal nodes resource failure")
		output = EncodeError(newErr)
		statusCode = 500
		return
	}

	baseLink := request.SelectedRoutePath()
	api.SetPageLink(baseLink, response, page)
	return
}

func userlistPage(list *api.UserServiceList, page pagination.Pager) (*api.UserServiceList, error) {

	listLen := len(list.Items)

	pageList := list

	var notPage bool
	if page == nil {
		notPage = true
	} else {
		notPage = page.Empty()
	}

	if !notPage {
		hasPage, perPage, skip := api.PagerToCondition(page, uint64(listLen))
		glog.V(5).Infof("Got page has %v  perpage %v skip %v\r\n", hasPage, perPage, skip)
		if hasPage {
			pageList.Items = list.Items[skip:perPage]
		}
	}

	return pageList, nil

}
