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
	"github.com/golang/glog"
)

//Config ...http server configure
type Config struct {
	SecurePort         int
	InsecurePort       int
	SwaggerPath        string
	TLSCertFile        string
	TLSPrivateKeyFile  string
	StorageClient      *backend.Backend
	EtcdStorageOptions *storageoptions.StorageOptions
}

//APIServer ... http server configure
type APIServer struct {
	SecurePort        int
	InsecurePort      int
	SwaggerPath       string
	TLSCertFile       string
	TLSPrivateKeyFile string
	wsContainer       *restful.Container
}

//New ...new a apiserver
func NewApiServer(config Config) *APIServer {
	comm.Storage = config.StorageClient

	err := comm.Storage.CreateStorage()
	if err != nil {
		glog.Errorf("CreateStorage failure %v \r\n", err)
		return nil
	}

	comm.EtcdStorage = etcdhelper.NewEtcdHelper(config.EtcdStorageOptions)
	if comm.EtcdStorage == nil {
		return nil
	}

	return &APIServer{
		SecurePort:        config.SecurePort,
		InsecurePort:      config.InsecurePort,
		SwaggerPath:       config.SwaggerPath,
		TLSCertFile:       config.TLSCertFile,
		TLSPrivateKeyFile: config.TLSPrivateKeyFile,
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

	err := controller.ControllerStart(comm.EtcdStorage, comm.Storage, port)
	if err != nil {
		return err
	}

	if tls {
		if len(apis.SwaggerPath) > 0 {
			apis.installSwaggerAPI(apis.wsContainer, true, port)
		}
		return server.ListenAndServeTLS(apis.TLSCertFile, apis.TLSPrivateKeyFile)
	} else {
		if len(apis.SwaggerPath) > 0 {
			apis.installSwaggerAPI(apis.wsContainer, false, port)
		}
		return server.ListenAndServe()
	}

}
