package apiserver

import (
	"net/http"
	"strconv"

	"shadowss/pkg/api"
	apierr "shadowss/pkg/api/errors"
	"shadowss/pkg/api/shadowssapi"

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
	ws.Route(ws.POST("").To(shadowssapi.PostLogin).
		Doc("get token").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), api.Login{}).
		Param(ws.BodyParameter("body", "identifier of the login").DataType("api.Login")).
		Operation("PostLogin"))

	wsApiServer := new(restful.WebService)
	wsApiServer.
		Path("/api/v1/apiservers").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) // you can specify this per route as well
	container.Add(wsApiServer)
	wsApiServer.Route(wsApiServer.GET("").To(shadowssapi.GetAPIServers).
		Doc("get api server").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), api.APIServerList{}).
		Operation("GetAPIServers"))
	wsApiServer.Route(wsApiServer.POST("").To(shadowssapi.PostAPIServer).
		Doc("post api server").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), api.APIServer{}).
		Param(ws.BodyParameter("body", "identifier of the login").DataType("api.APIServer")).
		Operation("PostAPIServer"))
	wsApiServer.Route(wsApiServer.DELETE("{id}").To(shadowssapi.DeleteAPIServer).
		Doc("delete api server by id").
		Param(ws.QueryParameter("pretty", "If 'true', then the output is pretty printed.")).
		Returns(http.StatusOK, http.StatusText(http.StatusOK), apierr.Status{}).
		Operation("DeleteAPIServer"))

	wsNode := new(restful.WebService)
	wsNode.
		Path("/api/v1/nodes").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) // you can specify this per route as well
	container.Add(wsNode)
	wsNode.Route(wsNode.GET("").To(shadowssapi.GetNodes).
		Doc("get nodes").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), api.NodeList{}).
		Operation("GetNodes"))
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
