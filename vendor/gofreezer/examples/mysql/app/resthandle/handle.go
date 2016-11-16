package resthandle

import (
	"encoding/json"

	"gofreezer/examples/common/apiext"
	apierr "gofreezer/examples/common/apiext/errors"
	"gofreezer/pkg/genericstoragecodec"
	"gofreezer/pkg/runtime"
	selection "gofreezer/pkg/storage"
	storage "gofreezer/pkg/storage/mysqls"

	apiComm "gofreezer/pkg/api"

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

	output, err = json.Marshal(login)
	statusCode = 200

	ctx := apiComm.NewContext()
	outLogin := new(apiext.Login)
	err = StorageCodec.Create(ctx, login.Name, login, outLogin)

	glog.Infof("Got return err %v out %v\r\n", err, outLogin)

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

	ctx := apiComm.NewContext()
	outLogin := new(apiext.Login)
	p := selection.SelectionPredicate{}
	p.Query = "name = ?"
	p.QueryArgs = name

	err := StorageCodec.Get(ctx, name, p, outLogin, true)

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

	ctx := apiComm.NewContext()
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

	ctx := apiComm.NewContext()
	outLogin := new(apiext.Login)
	p := selection.SelectionPredicate{}

	p.Query = "name = ?"
	p.QueryArgs = name

	err := StorageCodec.Delete(ctx, name, p, outLogin)
	glog.Infof("result %v\r\n", err)
	if err != nil {
		newErr := apierr.NewNotFound("invalid request name", name)
		output = encodeError(newErr)
		statusCode = 404
		return
	}

	output = apierr.NewSuccess().Encode()
	statusCode = 200

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

	ctx := apiComm.NewContext()
	p := selection.SelectionPredicate{}
	p.Query = "name = ?"
	p.QueryArgs = name

	err = StorageCodec.GuaranteedUpdate(ctx, name, p, login, true, func(input runtime.Object) (output runtime.Object, fields []string, err error) {
		login.Spec.User.Count = 5
		return login, []string{"Count"}, nil
	})

	if err != nil {
		newErr := apierr.NewNotFound("invalid request name", name)
		output = encodeError(newErr)
		statusCode = 404
	}

	output, err = runtime.Encode(GenericStorage.Codecs, login)
	if err != nil {
		newErr := apierr.NewInternalError("marshal nodes resource failure")
		output = encodeError(newErr)
		statusCode = 500
		return
	}

	output = apierr.NewSuccess().Encode()
	statusCode = 200

}
