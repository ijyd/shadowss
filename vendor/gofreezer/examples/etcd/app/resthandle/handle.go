package resthandle

import (
	"encoding/json"
	"gofreezer/examples/common/apiext"
	apierr "gofreezer/examples/common/apiext/errors"
	"gofreezer/pkg/genericstoragecodec"
	"gofreezer/pkg/runtime"
	storageinterface "gofreezer/pkg/storage"
	storage "gofreezer/pkg/storage/etcds"

	"gofreezer/pkg/api"

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

func PostLogin(request *restful.Request, response *restful.Response) {
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
	// err = NewToken(login)
	// if err != nil {
	// 	newErr := apierr.NewBadRequestError(err.Error())
	// 	output = encodeError(newErr)
	// 	statusCode = 400
	// } else {
	output, err = json.Marshal(login)
	statusCode = 200
	//	}
	ctx := api.NewContext()
	outLogin := new(apiext.Login)
	err = StorageCodec.Create(ctx, prefix+"/"+login.Name, login, outLogin, 0)

	glog.Infof("Got return err %v out %v\r\n", err, outLogin)

	return

}

//GetLoginList ... get login list
func GetLoginList(request *restful.Request, response *restful.Response) {

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

	options := &api.ListOptions{ResourceVersion: "0"}
	err := StorageCodec.List(ctx, prefix, options.ResourceVersion, storageinterface.Everything, outLogin)

	glog.Errorf("Get to list have error %v data %+v\r\n", err, outLogin)

	output, err = runtime.Encode(GenericStorage.Codecs, outLogin)
	if err != nil {
		glog.Errorf("encode failure %v \r\n", err)
	}

	return
}
