package userfile

import (
	"fmt"
	freezerapi "gofreezer/pkg/api"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/runtime"
	"os"

	"cloud-keeper/pkg/api"

	"io/ioutil"
	"net/http"
)

const (
	rootPath = "/userdata"
)

type REST struct {
	FileDesc   *FileDescREST
	FileStream *FileStreamREST
	File       *FileREST
}

type FileREST struct {
}

func NewREST() *REST {
	return &REST{
		FileDesc: &FileDescREST{},
		File:     &FileREST{},
	}
}

func (r *FileREST) New() runtime.Object {
	return &api.UserPublicFile{}
}

func (r *FileREST) NewList() runtime.Object {
	return &api.UserPublicFileList{}
}

func (rs *FileREST) List(ctx freezerapi.Context, options *freezerapi.ListOptions) (runtime.Object, error) {

	files, err := ioutil.ReadDir(rootPath)
	if err != nil {
		return nil, fmt.Errorf("search dir %v", err.Error())
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("Not found files")
	}

	userfileList := &api.UserPublicFileList{}

	for _, f := range files {
		fmt.Println(f.Name())
		name := fmt.Sprintf("%s\r\n", f.Name())
		userfile := api.UserPublicFile{
			Spec: api.UserPublicFileSpec{
				FileName: name,
			},
		}
		userfileList.Items = append(userfileList.Items, userfile)
	}

	return userfileList, nil
}

type FileStreamREST struct {
}

// Implement Connecter
var _ = rest.Connecter(&FileStreamREST{})

var userFileMethods = []string{"GET", "POST"}

// New creates a new UserPublicFiles  object
func (r *FileStreamREST) New() runtime.Object {
	// TODO - return a resource that represents a log
	return &api.UserPublicFile{}
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (r *FileStreamREST) ConnectMethods() []string {
	return userFileMethods
}

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *FileStreamREST) NewConnectOptions() (runtime.Object, bool, string) {
	return nil, false, ""
}

// Connect returns a handler for the pod proxy
func (r *FileStreamREST) Connect(ctx freezerapi.Context, id string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	return NewFileStreamer(responder), nil
}

type FileDescREST struct {
}

// New creates a new UserPublicFiles  object
func (r *FileDescREST) New() runtime.Object {
	// TODO - return a resource that represents a log
	return &api.UserPublicFile{}
}

func (rs *FileDescREST) Get(ctx freezerapi.Context, name string) (runtime.Object, error) {
	contentDescFile := rootPath + "/" + name + "_desc"
	if _, err := os.Stat(contentDescFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("not found %v", name)
	}

	output, err := ioutil.ReadFile(contentDescFile)
	if err != nil {
		return nil, err
	}

	userfile := &api.UserPublicFile{
		Spec: api.UserPublicFileSpec{
			FileName:    name,
			Description: string(output),
		},
	}

	return userfile, nil
}
