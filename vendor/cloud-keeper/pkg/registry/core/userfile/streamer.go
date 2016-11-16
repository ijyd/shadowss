package userfile

import (
	"gofreezer/pkg/api"
	"gofreezer/pkg/api/errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/golang/glog"
)

type FileStreamer struct {
	Responder ErrorResponder
}

// ErrorResponder abstracts error reporting to the proxy handler to remove the need to hardcode a particular
// error format.
type ErrorResponder interface {
	Error(err error)
}

func NewFileStreamer(responder ErrorResponder) *FileStreamer {
	return &FileStreamer{
		Responder: responder,
	}
}

// ServeHTTP handles the streamer request
func (h *FileStreamer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	fields := strings.Split(path, "/")
	if len(fields) < 3 {
		h.Responder.Error(errors.NewBadRequest("need userfiles/{name}/stream as a path"))
		return
	}

	filename := fields[0]

	switch req.Method {
	case "GET":
		h.DownLoadFile(w, req, filename)
	case "POST":
		h.UploadFile(w, req)
	default:
		h.Responder.Error(errors.NewMethodNotSupported(api.Resource("userfiles"), req.Method))
	}

	return
}

func (h *FileStreamer) UploadFile(w http.ResponseWriter, req *http.Request) {

	httpReq := req

	httpReq.ParseMultipartForm(32 << 20)
	glog.V(5).Infof("got req %+v \r\n", httpReq)
	file, handler, err := httpReq.FormFile("file")
	if err != nil {
		h.Responder.Error(errors.NewBadRequest("Form boyd required"))
		return
	}
	defer file.Close()

	writeFile := rootPath + "/" + handler.Filename
	glog.V(5).Infof("open file %v\r\n", writeFile)
	os.Remove(writeFile)

	f, err := os.OpenFile(writeFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		h.Responder.Error(errors.NewInternalError(err))
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
		h.Responder.Error(errors.NewInternalError(err))
	}
}

func (h *FileStreamer) DownLoadFile(w http.ResponseWriter, req *http.Request, fileName string) {

	w.Header().Set("Content-Type", "application/x-tgz")
	w.Header().Set("Content-Disposition:", `"attachment;filename="`+fileName+`"`)

	file := rootPath + "/" + fileName
	if _, err := os.Stat(file); os.IsNotExist(err) {
		h.Responder.Error(errors.NewNotFound(api.Resource("userfiles"), fileName))
		return
	}

	http.ServeFile(w, req, file)
	return
}
