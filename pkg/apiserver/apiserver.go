package apiserver

import (
	"fmt"
	"net/http"
	"shadowsocks-go/pkg/api/shadowssapi"
	"shadowsocks-go/pkg/backend"

	"github.com/emicklei/go-restful"
)

//Config ...http server configure
type Config struct {
	Host          string
	Port          int
	StorageClient *backend.Backend
}

//APIServer ... http server configure
type APIServer struct {
	Host        string
	Port        int
	wsContainer *restful.Container
}

//New ...new a apiserver
func NewApiServer(config Config) *APIServer {
	shadowssapi.Storage = config.StorageClient
	err := shadowssapi.Storage.CreateStorage()
	if err != nil {
		return nil
	}

	return &APIServer{
		Host: config.Host,
		Port: config.Port,
	}
}

//Run ...start http server run
func (apis *APIServer) Run() error {
	apis.wsContainer = restful.NewContainer()
	apis.wsContainer.Router(restful.CurlyRouter{})
	install(apis.wsContainer)

	// Add container filter to enable CORS
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST"},
		CookiesAllowed: false,
		Container:      apis.wsContainer}
	apis.wsContainer.Filter(cors.Filter)

	addr := apis.Host + ":" + fmt.Sprintf("%d", apis.Port)
	server := &http.Server{Addr: addr, Handler: apis.wsContainer}

	return server.ListenAndServe()
}
