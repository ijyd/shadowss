package vps

import (
	"encoding/json"
	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/api/unversioned"

	"cloud-keeper/pkg/api"
	apierr "cloud-keeper/pkg/api/errors"
	"cloud-keeper/pkg/api/validation"
	"cloud-keeper/pkg/controller"
	"cloud-keeper/pkg/pagination"

	restful "github.com/emicklei/go-restful"
	"github.com/golang/glog"
)

func getUsers(page pagination.Pager) ([]byte, int) {
	statusCode := 500
	var output []byte
	var err error

	var userList api.UserList
	userInfoArray, err := Storage.GetUserList(page)
	if err != nil {
		if isNotfoundErr(err) == true {
			userList = api.UserList{
				TypeMeta: unversioned.TypeMeta{
					Kind:       "UserList",
					APIVersion: "v1",
				},
				ListMeta: unversioned.ListMeta{
					SelfLink: "/api/v1/users",
				},
			}
		} else {
			glog.Errorf("Get user failure %v \r\n", err)
			newErr := apierr.NewInternalError(err.Error())
			output = encodeError(newErr)

			return output, statusCode
		}

	} else {

		var items []api.User
		for _, v := range userInfoArray {
			//get node from etcd fix our actual

			item := api.User{
				TypeMeta: unversioned.TypeMeta{
					Kind:       "User",
					APIVersion: "v1",
				},
				ObjectMeta: prototype.ObjectMeta{
					Name: v.Name,
				},
				Spec: api.UserSpec{
					DetailInfo: api.UserInfo{
						ID:              v.ID,
						Passwd:          v.Passwd,
						EnableOTA:       v.EnableOTA,
						TrafficLimit:    v.TrafficLimit,
						UploadTraffic:   v.UploadTraffic,
						DownloadTraffic: v.DownloadTraffic,
						Name:            v.Name,
						Email:           v.Email,
						ManagePasswd:    v.ManagePasswd,
						ExpireTime:      v.ExpireTime,
						RegIPAddr:       v.RegIPAddr,
						RegDBTime:       v.RegDBTime,
						Description:     v.Description,
						TrafficRate:     v.TrafficRate,
						IsAdmin:         v.IsAdmin,
					},
				},
			}
			items = append(items, item)
		}

		userList = api.UserList{
			TypeMeta: unversioned.TypeMeta{
				Kind:       "UserList",
				APIVersion: "v1",
			},
			ListMeta: unversioned.ListMeta{
				SelfLink: "/api/v1/users",
			},
			Items: items,
		}

		output, err = json.Marshal(userList)
		if err != nil {
			glog.Errorln("Marshal router err", err)
			newErr := apierr.NewInternalError("marshal user list resource failure")
			output = encodeError(newErr)
			statusCode = 500

		} else {
			statusCode = 200
		}
		return output, statusCode
	}

	return output, statusCode
}

//GetRouters ... get router list
func GetUsers(request *restful.Request, response *restful.Response) {

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

	output, statusCode = getUsers(page)
	baseLink := request.SelectedRoutePath()
	api.SetPageLink(baseLink, response, page)
}

func PostUser(request *restful.Request, response *restful.Response) {
	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")

	statusCode := 200
	output := apierr.NewSuccess().Encode()

	defer func() {
		w.WriteHeader(statusCode)
		w.Write(output)
	}()

	user := new(api.User)
	err := request.ReadEntity(user)
	if err != nil {
		glog.Errorf("invalid request body:%v", err)
		newErr := apierr.NewBadRequestError("request body invalid")
		output = encodeError(newErr)
		statusCode = 400
		return
	}

	err = validation.ValidateUser(*user)
	if err != nil {
		newErr := apierr.NewBadRequestError(err.Error())
		output = encodeError(newErr)
		statusCode = 400
		return
	}

	glog.Infof("Got Post api user:%+v\n", user)
	err = Storage.CreateUser(user.Spec.DetailInfo)
	if err != nil {
		newErr := apierr.NewBadRequestError(err.Error())
		output = encodeError(newErr)
		statusCode = 400
	} else {
		output, err = json.Marshal(user)
		statusCode = 200
	}

	err = controller.AllocNode(user)
	if err != nil {
		glog.Warningf("alloc node to user failure %v\r\n", err)
	}

	return
}

func DeleteUser(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")

	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	statusCode := 200
	var output []byte

	err := Storage.DeleteUserByName(name)
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
