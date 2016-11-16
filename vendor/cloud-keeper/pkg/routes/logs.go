package routes

import (
	"net/http"
	"path"

	"github.com/emicklei/go-restful"

	"apistack/pkg/genericapiserver/mux"
)

// Logs adds handlers for the /logs path serving log files from /var/log.
type Logs struct{}

func (l Logs) Install(c *mux.APIContainer) {
	// use restful: ws.Route(ws.GET("/logs/{logpath:*}").To(fileHandler))
	// See github.com/emicklei/go-restful/blob/master/examples/restful-serve-static.go
	ws := new(restful.WebService)
	ws.Path("/logs")
	ws.Doc("get log files")
	ws.Route(ws.GET("/{logpath:*}").To(logFileHandler).Param(ws.PathParameter("logpath", "path to the log").DataType("string")))
	ws.Route(ws.GET("/").To(logFileListHandler))

	c.Add(ws)
}

func logFileHandler(req *restful.Request, resp *restful.Response) {
	logdir := "/var/log"
	actual := path.Join(logdir, req.PathParameter("logpath"))
	http.ServeFile(resp.ResponseWriter, req.Request, actual)
}

func logFileListHandler(req *restful.Request, resp *restful.Response) {
	logdir := "/var/log"
	http.ServeFile(resp.ResponseWriter, req.Request, logdir)
}
