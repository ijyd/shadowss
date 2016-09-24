package apiserver

import (
	"fmt"
	"net/http"

	comm "cloud-keeper/pkg/api/vps/common"
	"cloud-keeper/pkg/backend"
	"cloud-keeper/pkg/controller"
	"cloud-keeper/pkg/etcdhelper"

	storageoptions "gofreezer/pkg/genericstoragecodec/options"

	"github.com/emicklei/go-restful"
)

//Config ...http server configure
type Config struct {
	Host               string
	Port               int
	SwaggerPath        string
	StorageClient      *backend.Backend
	EtcdStorageOptions *storageoptions.StorageOptions
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
	comm.Storage = config.StorageClient

	err := comm.Storage.CreateStorage()
	if err != nil {
		return nil
	}

	comm.EtcdStorage = etcdhelper.NewEtcdHelper(config.EtcdStorageOptions)
	if comm.EtcdStorage == nil {
		return nil
	}

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

	controller.ControllerStart(comm.EtcdStorage, comm.Storage, apis.Host, apis.Port)

	return server.ListenAndServe()
}
