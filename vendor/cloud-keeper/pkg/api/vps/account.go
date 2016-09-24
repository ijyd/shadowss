package vps

import (
	"encoding/json"
	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/api/unversioned"
	"golib/pkg/util/timewrap"

	"cloud-keeper/pkg/api"
	apierr "cloud-keeper/pkg/api/errors"
	"cloud-keeper/pkg/pagination"

	"github.com/golang/glog"

	restful "github.com/emicklei/go-restful"
)

func getAccount(page pagination.Pager) ([]byte, int) {
	statusCode := 500
	var output []byte
	var err error

	var acclist api.AccountList
	accdetails, err := Storage.GetAccounts(page)
	if err != nil {
		if isNotfoundErr(err) == true {
			acclist = api.AccountList{
				TypeMeta: unversioned.TypeMeta{
					Kind:       "AccountList",
					APIVersion: "v1",
				},
				ListMeta: unversioned.ListMeta{
					SelfLink: "/api/v1/accounts",
				},
			}
		} else {
			glog.Errorf("Get account failure %v \r\n", err)
			newErr := apierr.NewInternalError(err.Error())
			output = encodeError(newErr)
			return output, statusCode
		}

	} else {

		var accs []api.Account
		for _, v := range accdetails {
			item := api.Account{
				TypeMeta: unversioned.TypeMeta{
					Kind:       "Account",
					APIVersion: "v1",
				},
				ObjectMeta: prototype.ObjectMeta{
					Name: v.Name,
				},
				Spec: api.AccountSpec{
					AccDetail: api.AccountDetail{
						Name:           v.Name,
						Operators:      v.Operators,
						Key:            v.Key,
						CreditCeilings: v.CreditCeilings,
						Lables:         v.Lables,
						Descryption:    v.Descryption,
						CreateTime:     timewrap.NewTime(v.CreateDBTime),
						ExpireTime:     timewrap.NewTime(v.ExpireDBTime),
					},
				},
			}
			accs = append(accs, item)
		}

		acclist = api.AccountList{
			TypeMeta: unversioned.TypeMeta{
				Kind:       "AccountList",
				APIVersion: "v1",
			},
			ListMeta: unversioned.ListMeta{
				SelfLink: "/api/v1/accounts",
			},
			Items: accs,
		}

	}

	output, err = json.Marshal(acclist)
	if err != nil {
		glog.Errorln("Marshal router err", err)
		newErr := apierr.NewInternalError("marshal apiserver list resource failure")
		output = encodeError(newErr)

	} else {
		statusCode = 200
	}
	return output, statusCode
}

//GetAccounts ... get account list
func GetAccounts(request *restful.Request, response *restful.Response) {

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

	glog.V(5).Infof("Got page %+v \r\n", page)
	output, statusCode = getAccount(page)
	baseLink := request.SelectedRoutePath()
	api.SetPageLink(baseLink, response, page)

	return
}

func PostAccount(request *restful.Request, response *restful.Response) {
	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	statusCode := 200

	output := apierr.NewSuccess().Encode()

	acc := new(api.Account)
	err := request.ReadEntity(acc)
	if err != nil {
		glog.Errorf("invalid request body:%v", err)
		newErr := apierr.NewBadRequestError("request body invalid")
		output = encodeError(newErr)
		statusCode = 400
	} else {
		glog.Infof("Got Post account: %+v\n", acc)
		err := Storage.CreateAccount(acc.Spec.AccDetail)
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
			output, err = json.Marshal(acc)
			statusCode = 200
		}
	}

	w.WriteHeader(statusCode)
	w.Write(output)
}

func DeleteAccount(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")

	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	statusCode := 200
	var output []byte

	err := Storage.DeleteAccount(name)
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
