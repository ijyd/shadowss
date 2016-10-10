package vps

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	apierr "cloud-keeper/pkg/api/errors"
	. "cloud-keeper/pkg/api/vps/common"

	restful "github.com/emicklei/go-restful"
	"github.com/golang/glog"
)

const (
	rootPath = "/userdata"
)

func GetFile(request *restful.Request, response *restful.Response) {

	fileName := request.PathParameter("name")
	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	encoded := request.Request.Header.Get("Authorization")
	statusCode := 200

	defer func() {
		w.WriteHeader(statusCode)
	}()

	user, err := CheckUserLevelToken(encoded)
	if err != nil || user == nil {
		glog.Errorln("Unauth request ", err)
		newErr := apierr.NewUnauthorized("invalid token")
		output := EncodeError(newErr)
		w.Write(output)
		statusCode = 401
		return
	}
	w.Header().Set("Content-Type", "application/x-tgz")
	w.Header().Set("Content-Disposition:", `"attachment;filename="`+fileName+`"`)

	file := rootPath + "/" + fileName
	if _, err := os.Stat(file); os.IsNotExist(err) {
		newErr := apierr.NewNotFound("not found", fileName)
		output := EncodeError(newErr)
		w.Write(output)
		statusCode = 404
		return
	}

	http.ServeFile(response.ResponseWriter, request.Request, file)
}

func GetFileDesc(request *restful.Request, response *restful.Response) {
	fileName := request.PathParameter("name")
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
		w.Write(output)
		statusCode = 401
		return
	}

	contentDescFile := rootPath + "/" + fileName + "_desc"
	if _, err = os.Stat(contentDescFile); os.IsNotExist(err) {
		newErr := apierr.NewNotFound("not found", fileName)
		output = EncodeError(newErr)
		statusCode = 404
		return
	}

	output, err = ioutil.ReadFile(contentDescFile)
	if err != nil {
		newErr := apierr.NewInternalError(err.Error())
		output = EncodeError(newErr)
		statusCode = 500
		return
	}

}

func GetFileList(request *restful.Request, response *restful.Response) {
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
		w.Write(output)
		statusCode = 401
		return
	}

	files, err := ioutil.ReadDir(rootPath)
	if err != nil {
		newErr := apierr.NewInternalError(err.Error())
		output = EncodeError(newErr)
		statusCode = 500
	}

	if len(files) == 0 {
		newErr := apierr.NewNotFound("not found", "")
		output = EncodeError(newErr)
		statusCode = 404
		return
	}

	for _, f := range files {
		fmt.Println(f.Name())
		name := fmt.Sprintf("%s\r\n", f.Name())
		output = append(output, []byte(name)...)
	}

}

func PostPublicFiles(request *restful.Request, response *restful.Response) {
	w := response.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	encoded := request.Request.Header.Get("Authorization")

	statusCode := 200
	output := apierr.NewSuccess().Encode()

	defer func() {
		w.WriteHeader(statusCode)
		w.Write(output)
	}()

	glog.V(5).Infof("test here for 1")
	tokenUser, err := CheckToken(encoded)
	if err != nil || tokenUser == nil {
		glog.Errorln("Unauth request ", err)
		newErr := apierr.NewUnauthorized("invalid token")
		output = EncodeError(newErr)
		statusCode = 401
		return
	}

	httpReq := request.Request

	httpReq.ParseMultipartForm(32 << 20)
	glog.V(5).Infof("got req %+v \r\n", httpReq)
	file, handler, err := httpReq.FormFile("file")
	if err != nil {
		glog.V(5).Infof("got file error %v \r\n", err)
		newErr := apierr.NewInternalError(err.Error())
		output = EncodeError(newErr)
		statusCode = 500
		return
	}
	defer file.Close()

	writeFile := rootPath + "/" + handler.Filename
	glog.V(5).Infof("open file %v\r\n", writeFile)
	os.Remove(writeFile)

	f, err := os.OpenFile(writeFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		newErr := apierr.NewInternalError(err.Error())
		output = EncodeError(newErr)
		statusCode = 500
		return
	}
	defer f.Close()
	io.Copy(f, file)

	content := httpReq.FormValue("desc")
	contentDescFile := rootPath + "/" + handler.Filename + "_desc"

	glog.V(5).Infof("write desc file %v \r\n", contentDescFile)
	os.Remove(contentDescFile)

	err = ioutil.WriteFile(contentDescFile, []byte(content), 0666)
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
