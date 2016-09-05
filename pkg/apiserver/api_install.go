package apiserver

import (
	"shadowsocks-go/pkg/api/shadowssapi"

	"github.com/emicklei/go-restful"
)

func installDeviceResource(container *restful.Container) {

	ws := new(restful.WebService)
	ws.
		Path("/api/v1/logins").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) // you can specify this per route as well

	ws.Route(ws.POST("").To(shadowssapi.PostLogin).
		Doc("get token").
		Operation("PostLogin"))

	wsApiServer := new(restful.WebService)
	wsApiServer.
		Path("/api/v1/apiservers").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) // you can specify this per route as well

	wsApiServer.Route(wsApiServer.GET("").To(shadowssapi.GetAPIServers).
		Doc("get api server").
		Operation("GetAPIServers"))

	//container.Add(wsApiServer)

	wsNode := new(restful.WebService)
	wsNode.
		Path("/api/v1/nodes").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) // you can specify this per route as well

	wsNode.Route(wsNode.GET("").To(shadowssapi.GetNodes).
		Doc("get nodes").
		Operation("GetNodes"))

	container.Add(ws)
	container.Add(wsApiServer)
	container.Add(wsNode)
}

func install(container *restful.Container) error {

	installDeviceResource(container)
	return nil
}
