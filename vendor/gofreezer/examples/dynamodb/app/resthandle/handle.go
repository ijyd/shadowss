package resthandle

import (
	"time"

	"gofreezer/examples/common/apiext"
	apierr "gofreezer/examples/common/apiext/errors"
	"gofreezer/pkg/api"
	"gofreezer/pkg/api/unversioned"
	"gofreezer/pkg/genericstoragecodec"
	"gofreezer/pkg/runtime"
	selection "gofreezer/pkg/storage"
	storage "gofreezer/pkg/storage/awsdynamodb"
	"gofreezer/pkg/util/uuid"

	restful "github.com/emicklei/go-restful"
	"github.com/golang/glog"
)

// var StorageHandle storage.Interface
// var RunTimeCodecs runtime.Codec

var GenericStorage *genericstoragecodec.GenericStorageCodec
var StorageCodec storage.Interface

const (
	prefix = "/" + "Login"
)

func PostLoginUser(request *restful.Request, response *restful.Response) {
	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")

	statusCode := 200
	output := apierr.NewSuccess().Encode()

	defer func() {
		w.WriteHeader(statusCode)
		w.Write(output)
	}()

	login := new(apiext.Login)
	err := request.ReadEntity(login)
	if err != nil {
		glog.Errorf("invalid request body:%v", err)
		newErr := apierr.NewBadRequestError("request body invalid")
		output = encodeError(newErr)
		statusCode = 400
		return
	}

	glog.Infof("Got Post logins:%+v\n", login)

	ctx := api.NewContext()
	outLogin := new(apiext.Login)
	login.CreationTimestamp = unversioned.NewTime(time.Now())

	login.DeletionTimestamp = nil
	login.DeletionGracePeriodSeconds = nil
	login.UID = uuid.NewUUID()

	err = StorageCodec.Create(ctx, login.Name, login, outLogin, 0)
	if err != nil {
		newErr := apierr.NewInternalError(err.Error())
		output = encodeError(newErr)
		statusCode = 500
		return
	}
	output, err = runtime.Encode(GenericStorage.Codecs, outLogin)
	if err != nil {
		newErr := apierr.NewInternalError("marshal nodes resource failure")
		output = encodeError(newErr)
		statusCode = 500
		return
	}

	return
}

//GetLoginUser ... get user by name
func GetLoginUser(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")

	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")

	statusCode := 200
	var output []byte

	defer func() {
		w.WriteHeader(statusCode)
		w.Write(output)
	}()

	ctx := api.NewContext()
	outLogin := new(apiext.Login)
	// p := storage.SelectionPredicate{
	// 	Query:     "name = ?",
	// 	QueryArgs: name,
	// }

	err := StorageCodec.Get(ctx, name, outLogin, true)

	glog.V(5).Infof("Get result %v\r\n", err)
	output, err = runtime.Encode(GenericStorage.Codecs, outLogin)
	if err != nil {
		newErr := apierr.NewInternalError("marshal nodes resource failure")
		output = encodeError(newErr)
		statusCode = 500
		return
	}
}

//GetLoginList ... get login list
func GetLoginUserList(request *restful.Request, response *restful.Response) {

	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	statusCode := 200
	var output []byte

	defer func() {
		w.WriteHeader(statusCode)
		w.Write(output)
	}()

	ctx := api.NewContext()
	outLogin := new(apiext.LoginList)

	p := selection.SelectionPredicate{}
	err := StorageCodec.GetToList(ctx, "LoginUserList", p, outLogin)

	glog.Errorf("Get to list have error %v data %+v\r\n", err, outLogin)

	output, err = runtime.Encode(GenericStorage.Codecs, outLogin)
	if err != nil {
		glog.Errorf("encode failure %v \r\n", err)
	}

	return
}

func DeleteLoginUser(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")

	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	statusCode := 200
	var output []byte
	//encoded := request.Request.Header.Get("Authorization")

	defer func() {
		w.WriteHeader(statusCode)
		w.Write(output)
	}()

	ctx := api.NewContext()
	outLogin := new(apiext.Login)
	// p := storage.SelectionPredicate{
	// 	Query:     "name = ?",
	// 	QueryArgs: name,
	// }

	//precondtion := storage.NewUIDPreconditions(string(""))

	err := StorageCodec.Delete(ctx, name, outLogin, nil)
	glog.Infof("result %v\r\n", err)
	if err != nil {
		newErr := apierr.NewNotFound("invalid request name", name)
		output = encodeError(newErr)
		statusCode = 404
		return
	}

	output, err = runtime.Encode(GenericStorage.Codecs, outLogin)
	if err != nil {
		newErr := apierr.NewInternalError("marshal nodes resource failure")
		output = encodeError(newErr)
		statusCode = 500
	}

}

//PutLoginUser ... update user field
func PutLoginUser(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")

	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	statusCode := 200
	var output []byte

	defer func() {
		w.WriteHeader(statusCode)
		w.Write(output)
	}()

	login := new(apiext.Login)
	err := request.ReadEntity(login)
	if err != nil {
		glog.Errorf("invalid request body:%v", err)
		newErr := apierr.NewBadRequestError("request body invalid")
		output = encodeError(newErr)
		statusCode = 400
		return
	}

	ctx := api.NewContext()
	// p := storage.SelectionPredicate{
	// 	Query:     "name = ?",
	// 	QueryArgs: name,
	// }
	err = StorageCodec.GuaranteedUpdate(ctx, name, login, true, nil, func(input runtime.Object, attributeValues map[string]interface{}) (output runtime.Object, ttl *uint64, err error) {
		newLogin := input.(*apiext.Login)

		userName := string("test here")
		//count := 5
		newLogin.Spec.User.Count = 5
		newLogin.Spec.User.UserName = userName
		expire := uint64(0)
		attributeValues["spec.#user.#userName"] = userName
		//attributeValues["spec.#user.#count"] = count

		return newLogin, &expire, nil
	})

	if err != nil {
		newErr := apierr.NewNotFound(err.Error(), name)
		output = encodeError(newErr)
		statusCode = 404
		return
	}

	output, err = runtime.Encode(GenericStorage.Codecs, login)
	if err != nil {
		newErr := apierr.NewInternalError(err.Error())
		output = encodeError(newErr)
		statusCode = 500
		return
	}

}
