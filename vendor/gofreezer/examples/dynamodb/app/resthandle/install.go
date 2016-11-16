package resthandle

import (
	"net/http"

	"gofreezer/examples/common/apiext"
	apierr "gofreezer/examples/common/apiext/errors"

	restful "github.com/emicklei/go-restful"
)

func InstallWS(container *restful.Container) (err error) {
	ws := new(restful.WebService)
	ws.
		Path("/api/v1/logins").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) // you can specify this per route as well
	container.Add(ws)

	route := ws.POST("").To(PostLoginUser).
		Doc("get token").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), apiext.Login{}).
		Param(ws.BodyParameter("body", "identifier of the login").DataType("api.Login")).
		Operation("PostLoginUser")
	ws.Route(route)

	route = ws.GET("").To(GetLoginUserList).
		Doc("get login list").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), apiext.Login{}).
		Operation("GetLoginUserList")
	ws.Route(route)

	route = ws.GET("{name}").To(GetLoginUser).
		Doc("get login user").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), apiext.Login{}).
		Operation("GetLoginUser")
	ws.Route(route)

	route = ws.DELETE("{name}").To(DeleteLoginUser).
		Doc("delete login user").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), apierr.Status{}).
		Operation("DeleteLoginUser")
	ws.Route(route)

	route = ws.PUT("{name}").To(PutLoginUser).
		Doc("put user's properties").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), apiext.Login{}).
		Param(ws.BodyParameter("body", "identifier of the login").DataType("api.UserService.Annotations")).
		Operation("PutLoginUser")
	ws.Route(route)

	return nil
}
