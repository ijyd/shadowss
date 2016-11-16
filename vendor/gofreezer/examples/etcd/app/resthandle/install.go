package resthandle

import (
	"gofreezer/examples/common/apiext"
	"net/http"

	restful "github.com/emicklei/go-restful"
)

func InstallWS(container *restful.Container) (err error) {
	ws := new(restful.WebService)
	ws.
		Path("/api/v1/logins").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) // you can specify this per route as well
	container.Add(ws)

	route := ws.POST("").To(PostLogin).
		Doc("get token").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), apiext.Login{}).
		Param(ws.BodyParameter("body", "identifier of the login").DataType("api.Login")).
		Operation("PostLogin")
	ws.Route(route)

	route = ws.GET("").To(GetLoginList).
		Doc("get login list").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), apiext.Login{}).
		Operation("GetLoginList")
	ws.Route(route)

	return nil
}
