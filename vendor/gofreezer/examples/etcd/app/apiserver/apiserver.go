package apiserver

import (
	"fmt"
	"gofreezer/examples/etcd/app/api/resthandle"
	"gofreezer/pkg/genericstoragecodec"
	"net/http"

	"github.com/emicklei/go-restful"
)

//Config ...http server configure
type Config struct {
	Host         string
	Port         int
	SwaggerPath  string
	StorageCodec *genericstoragecodec.GenericStorageCodec
}

//APIServer ... http server configure
type APIServer struct {
	Host        string
	Port        int
	SwaggerPath string
	wsContainer *restful.Container
}

//New ...new a apiserver
func NewApiServer(config Config) *APIServer {
	resthandle.StorageCodec = config.StorageCodec

	return &APIServer{
		Host:        config.Host,
		Port:        config.Port,
		SwaggerPath: config.SwaggerPath,
	}
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

	addr := apis.Host + ":" + fmt.Sprintf("%d", apis.Port)
	server := &http.Server{Addr: addr, Handler: apis.wsContainer}

	return server.ListenAndServe()
}