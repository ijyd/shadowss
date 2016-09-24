package apiserver

import (
	"net/http"
	"strconv"

	"cloud-keeper/pkg/api"
	apierr "cloud-keeper/pkg/api/errors"
	"cloud-keeper/pkg/api/vps"
	"cloud-keeper/pkg/api/vps/etcdrest"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
)

func installControllerSrv(container *restful.Container) {
	// ws := new(restful.WebService)
	// ws.
	// 	Path("/api/v1/").
	// 	Consumes(restful.MIME_JSON).
	// 	Produces(restful.MIME_JSON) // you can specify this per route as well
	// container.Add(ws)
	// route := ws.POST("").To(controller.StartConfigNodeUser).
	// 	Doc("reinit node config").
	// 	Returns(http.StatusOK, http.StatusText(http.StatusOK), apierr.Status{}).
	// 	Operation("StartConfigNodeUser")
	// ws.Route(route)
	//
	// route = ws.GET("").To(controller.AddNodeConfig).
	// 	Doc("add node config").
	// 	Returns(http.StatusOK, http.StatusText(http.StatusOK), apierr.Status{}).
	// 	Operation("AddNodeConfig")
	// ws.Route(route)
}

func installUserSrv(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/api/v1/users").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) // you can specify this per route as well
	container.Add(ws)

	route := ws.GET("").To(vps.GetUsers).
		Doc("get api server").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), api.UserList{}).
		Operation("GetUsers")
	ws.Route(route)

	route = ws.POST("").To(vps.PostUser).
		Doc("post user").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), api.User{}).
		Param(ws.BodyParameter("body", "identifier of the login").DataType("api.User")).
		Operation("PostUser")
	ws.Route(route)
	route = ws.DELETE("{name}").To(vps.DeleteUser).
		Doc("delete user by name").
		Param(ws.QueryParameter("pretty", "If 'true', then the output is pretty printed.")).
		Returns(http.StatusOK, http.StatusText(http.StatusOK), apierr.Status{}).
		Operation("DeleteUser")
	ws.Route(route)

	route = ws.GET("{name}/bindingnodes").To(etcdrest.GetBindingNodes).
		Doc("get user binding nodes").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), api.UserServiceList{}).
		Operation("GetBindingNodes")
	ws.Route(route)
}

func installNodeSrv(container *restful.Container) {
	wsNode := new(restful.WebService)
	wsNode.
		Path("/api/v1/nodes").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) // you can specify this per route as well
	container.Add(wsNode)
	route := wsNode.GET("").To(vps.GetNodes).
		Doc("get nodes").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), api.NodeList{}).
		Operation("GetNodes")
	wsNode.Route(route)

	route = wsNode.POST("").To(vps.PostNode).
		Doc("post api server").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), api.Node{}).
		Param(wsNode.BodyParameter("body", "identifier of the login").DataType("api.Node")).
		Operation("PostNode")
	wsNode.Route(route)
	route = wsNode.DELETE("{name}").To(vps.DeleteNode).
		Doc("delete api server by name").
		Param(wsNode.QueryParameter("pretty", "If 'true', then the output is pretty printed.")).
		Returns(http.StatusOK, http.StatusText(http.StatusOK), apierr.Status{}).
		Operation("DeleteNode")
	wsNode.Route(route)

	route = wsNode.GET("{name}/bindingusers").To(etcdrest.GetBindingUsers).
		Doc("get nodes binding users").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), api.NodeUserList{}).
		Operation("GetBindingUsers")
	wsNode.Route(route)

}

func installAPIServerSrv(container *restful.Container) {
	wsApiServer := new(restful.WebService)
	wsApiServer.
		Path("/api/v1/apiservers").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) // you can specify this per route as well
	container.Add(wsApiServer)

	route := wsApiServer.GET("").To(vps.GetAPIServers).
		Doc("get api server").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), api.APIServerList{}).
		Operation("GetAPIServers")
	wsApiServer.Route(route)

	// route = wsApiServer.POST("").To(vps.PostAPIServer).
	// 	Doc("post api server").
	// 	Returns(http.StatusOK, http.StatusText(http.StatusOK), api.APIServer{}).
	// 	Param(wsApiServer.BodyParameter("body", "identifier of the login").DataType("api.APIServer")).
	// 	Operation("PostAPIServer")
	// wsApiServer.Route(route)
	// route = wsApiServer.DELETE("{name}").To(vps.DeleteAPIServer).
	// 	Doc("delete api server by name").
	// 	Param(wsApiServer.QueryParameter("pretty", "If 'true', then the output is pretty printed.")).
	// 	Returns(http.StatusOK, http.StatusText(http.StatusOK), apierr.Status{}).
	// 	Operation("DeleteAPIServer")
	// wsApiServer.Route(route)
}

func installLoginSrv(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/api/v1/logins").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) // you can specify this per route as well
	container.Add(ws)
	route := ws.POST("").To(vps.PostLogin).
		Doc("get token").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), api.Login{}).
		Param(ws.BodyParameter("body", "identifier of the login").DataType("api.Login")).
		Operation("PostLogin")
	ws.Route(route)
}

func installAccountResource(container *restful.Container) {

	wsAccount := new(restful.WebService)
	wsAccount.
		Path("/api/v1/accounts").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) // you can specify this per route as well
	container.Add(wsAccount)
	route := wsAccount.GET("").To(vps.GetAccounts).
		Doc("get account list").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), api.AccountList{}).
		Operation("GetAccounts")
	wsAccount.Route(route)
	route = wsAccount.GET("{name}/info").To(vps.GetAccountInfo).
		Doc("get account information").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), api.AccountInfo{}).
		Operation("GetAccountInfo")
	wsAccount.Route(route)

	route = wsAccount.POST("").To(vps.PostAccount).
		Doc("post account").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), api.Account{}).
		Param(wsAccount.BodyParameter("body", "identifier of the api key").DataType("api.Account")).
		Operation("PostAccount")
	wsAccount.Route(route)
	route = wsAccount.DELETE("{name}").To(vps.DeleteAccount).
		Doc("delete account by name").
		Param(wsAccount.QueryParameter("pretty", "If 'true', then the output is pretty printed.")).
		Returns(http.StatusOK, http.StatusText(http.StatusOK), apierr.Status{}).
		Operation("DeleteAccount")
	wsAccount.Route(route)

	//post account server

	route = wsAccount.GET("{name}/servers").To(vps.GetAccServers).
		Doc("get account servers").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), api.AccServerList{}).
		Operation("GetAccServers")
	wsAccount.Route(route)

	route = wsAccount.POST("{name}/servers").To(vps.PostAccServer).
		Doc("post server with account").
		Param(wsAccount.QueryParameter("pretty", "If 'true', then the output is pretty printed.")).
		Operation("PostAccServer").
		Produces(restful.MIME_JSON).
		Returns(http.StatusOK, http.StatusText(http.StatusOK), apierr.Status{}).
		Reads(api.AccServer{})
	wsAccount.Route(route)

	route = wsAccount.DELETE("{name}/servers/{id}").To(vps.DeleteAccServer).
		Doc("delete server by id").
		Param(wsAccount.QueryParameter("pretty", "If 'true', then the output is pretty printed.")).
		Returns(http.StatusOK, http.StatusText(http.StatusOK), apierr.Status{}).
		Operation("DeleteAccServer")
	wsAccount.Route(route)

	route = wsAccount.POST("{name}/servers/exec").To(vps.PostAccServerExec).
		Doc("exec command on server").
		Param(wsAccount.QueryParameter("pretty", "If 'true', then the output is pretty printed.")).
		Operation("PostAccServerExec").
		Produces(restful.MIME_JSON).
		Returns(http.StatusOK, http.StatusText(http.StatusOK), apierr.Status{}).
		Reads(api.AccServerCommand{})
	wsAccount.Route(route)

	route = wsAccount.GET("{name}/servers/sshkeys").To(vps.GetAccServersSSHkey).
		Doc("get servers ssh key").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), api.AccServerSSHKey{}).
		Operation("GetAccServersSSHkey")
	wsAccount.Route(route)

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
		// SchemaFormatHandler: func(typeName string) string {
		// 	switch typeName {
		// 	case "unversioned.Time", "*unversioned.Time":
		// 		return "date-time"
		// 	}
		// 	return ""
		// },
	}
	swagger.RegisterSwaggerService(swaggerConfig, container)
}

func (apis *APIServer) install(container *restful.Container) error {

	installAccountResource(container)
	installControllerSrv(container)
	installAPIServerSrv(container)
	installLoginSrv(container)
	installNodeSrv(container)
	installUserSrv(container)

	if len(apis.SwaggerPath) > 0 {
		apis.installSwaggerAPI(container)
	}

	return nil
}
