package apiserver

import (
	"net/http"
	"strconv"

	"gofreezer/examples/etcd/app/api"
	"gofreezer/examples/etcd/app/api/resthandle"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
)

func installWebServuce(container *restful.Container) {

	ws := new(restful.WebService)
	ws.
		Path("/api/v1/logins").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) // you can specify this per route as well
	container.Add(ws)
	route := ws.POST("").To(resthandle.PostLogin).
		Doc("get token").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), api.Login{}).
		Param(ws.BodyParameter("body", "identifier of the login").DataType("api.Login")).
		Operation("PostLogin")
	ws.Route(route)

	route = ws.GET("").To(resthandle.GetLoginList).
		Doc("get login list").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), api.Login{}).
		Operation("GetLoginList")
	ws.Route(route)
}

func (apis *APIServer) installSwaggerAPI(container *restful.Container) {
	hostAndPort := apis.Host + string(":") + strconv.Itoa(apis.Port)
	//protocol := "https://"
	protocol := "http://"
	webServicesUrl := protocol + hostAndPort

	// Enable swagger UI and discovery API
	swaggerConfig := swagger.Config{
		WebServicesUrl:  webServicesUrl,
		WebServices:     container.RegisteredWebServices(),
		ApiPath:         "/swaggerapi/",
		SwaggerPath:     "/swaggerui/",
		SwaggerFilePath: apis.SwaggerPath,
		SchemaFormatHandler: func(typeName string) string {
			switch typeName {
			case "unversioned.Time", "*unversioned.Time":
				return "date-time"
			}
			return ""
		},
	}
	swagger.RegisterSwaggerService(swaggerConfig, container)
}

func (apis *APIServer) install(container *restful.Container) error {

	installWebServuce(container)
	if len(apis.SwaggerPath) > 0 {
		apis.installSwaggerAPI(container)
	}

	return nil
}
