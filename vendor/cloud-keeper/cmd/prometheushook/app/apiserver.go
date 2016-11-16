package app

import (
	"fmt"
	"net/http"

	apierr "cloud-keeper/cmd/prometheushook/app/errors"
	"cloud-keeper/cmd/prometheushook/app/rest"

	restful "github.com/emicklei/go-restful"
	"github.com/golang/glog"
)

//APIServer ... http server configure
type APIServer struct {
	SecurePort        int
	InsecurePort      int
	TLSCertFile       string
	TLSPrivateKeyFile string
	wsContainer       *restful.Container
}

func (apis *APIServer) install(container *restful.Container) error {

	ws := new(restful.WebService)
	ws.
		Path("/api/v1/alerthook").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) // you can specify this per route as well
	container.Add(ws)
	route := ws.POST("").To(rest.PostAlert).
		Doc("prometheus web hook").
		Returns(http.StatusOK, http.StatusText(http.StatusOK), apierr.Status{}).
		Param(ws.BodyParameter("body", "identifier of the alert").DataType("api.PostAlert")).
		Operation("PostAlert")
	ws.Route(route)
	return nil
}

//Run ...start http server run
func (apis *APIServer) Run() error {

	apis.wsContainer = restful.NewContainer()
	apis.wsContainer.Router(restful.CurlyRouter{})
	apis.install(apis.wsContainer)

	// Add container filter to enable CORS
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST"},
		CookiesAllowed: false,
		Container:      apis.wsContainer}
	apis.wsContainer.Filter(cors.Filter)

	port := apis.InsecurePort
	var tls bool
	if apis.SecurePort != 0 {
		tls = true
		port = apis.SecurePort
		if apis.TLSCertFile == "" || apis.TLSPrivateKeyFile == "" {
			return fmt.Errorf("must give cert and private key file")
		}
	}

	addr := ":" + fmt.Sprintf("%d", port)
	glog.V(5).Infof("server on (%v %v)", apis.InsecurePort, apis.SecurePort)
	server := &http.Server{Addr: addr, Handler: apis.wsContainer}

	if tls {
		return server.ListenAndServeTLS(apis.TLSCertFile, apis.TLSPrivateKeyFile)
	} else {
		return server.ListenAndServe()
	}

}
