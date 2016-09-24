package resthandle

import (
	"encoding/json"
	"gofreezer/examples/etcd/app/api"
	apierr "gofreezer/examples/etcd/app/api/errors"
	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/genericstoragecodec"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/storage"

	apiComm "gofreezer/pkg/api"

	restful "github.com/emicklei/go-restful"
	"github.com/golang/glog"
)

// var StorageHandle storage.Interface
// var RunTimeCodecs runtime.Codec

var StorageCodec *genericstoragecodec.GenericStorageCodec

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

	login := new(api.Login)
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
	ctx := apiComm.NewContext()
	outLogin := new(api.Login)
	err = StorageCodec.Storage.Create(ctx, prefix+"/"+login.Name, login, outLogin, 0)

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

	ctx := apiComm.NewContext()
	outLogin := new(api.LoginList)

	options := &prototype.ListOptions{ResourceVersion: "0"}
	err := StorageCodec.Storage.List(ctx, prefix, options.ResourceVersion, storage.Everything, outLogin)

	glog.Errorf("Get to list have error %v data %+v\r\n", err, outLogin)

	output, err = runtime.Encode(StorageCodec.Codecs, outLogin)
	if err != nil {
		glog.Errorf("encode failure %v \r\n", err)
	}

	return
}
