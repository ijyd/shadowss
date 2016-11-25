package master

import (
	"apistack/pkg/genericapiserver"
	"apistack/pkg/util/async"
	freezerapi "gofreezer/pkg/api"
	apierrs "gofreezer/pkg/api/errors"
	"gofreezer/pkg/api/unversioned"
	"golib/pkg/util/network"
	"strings"
	"time"

	"github.com/golang/glog"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/apiserver"
	corerest "cloud-keeper/pkg/registry/core/rest"
	"cloud-keeper/pkg/registry/core/user"
)

// // Controller is the controller manager for the core bootstrap cloud-keeper controller
// // loops
type Controller struct {
	UserRegistry      user.Registry
	APIServerRegistry apiserver.Registry
	runner            *async.Runner
	port              int
}

// NewBootstrapController returns a controller for watching the core capabilities of the master
func (c *Config) NewBootstrapController(legacyRESTStorage corerest.LegacyRESTStorage, port int) *Controller {
	InnerHookHandler.SetRegistry(legacyRESTStorage)
	return &Controller{
		UserRegistry:      legacyRESTStorage.UserRegistry,
		APIServerRegistry: legacyRESTStorage.APIServerRegistry,
		port:              port,
	}
}

func (c *Controller) PostStartHook(hookContext genericapiserver.PostStartHookContext) error {
	c.Start()
	return nil
}

// Start begins the core controller loops that must exist for bootstrapping
// a cluster.
func (c *Controller) Start() {
	if c.runner != nil {
		return
	}

	//publish api server
	c.PublishAPIServer()
	// glog.V(5).Infof("begin migrate userservier")
	// c.UserRegistry.MigrateUser()

	//not have any runner
	// c.runner = async.NewRunner(c.RunKubernetesNamespaces, c.RunKubernetesService, repairClusterIPs.RunUntil, repairNodePorts.RunUntil)
	// c.runner.Start()
}

func (c *Controller) PublishAPIServer() {
	//add apiserver node
	apiserverName, err := network.ExternalMAC()
	if err != nil {
		glog.Fatalf("Publish api server error:%v\r\n", err.Error())
		return
	}
	apiserverName = strings.Replace(apiserverName, ":", "", -1)

	ctx := freezerapi.NewContext()
	_, err = c.APIServerRegistry.GetAPIServer(ctx, apiserverName)
	if err != nil && !apierrs.IsNotFound(err) {
		glog.Fatalf("Publish api server error:%v\r\n", err.Error())
		return
	}

	if apierrs.IsNotFound(err) {
		var hostList []string
		localExternalHost, err := network.ExternalIP()
		if err != nil {
			glog.Fatalf("Publish api server error:%v\r\n", err.Error())
			return
		}
		hostList = append(hostList, localExternalHost)

		internetIP, err := network.ExternalInternetIP()
		if err != nil {
			glog.Fatalf("Publish api server error:%v\r\n", err.Error())
			return
		}
		internetIP = strings.Replace(internetIP, "\n", "", -1)
		hostList = append(hostList, internetIP)

		spec := api.APIServerSpec{
			Server: api.APIServerInfor{
				Name:       apiserverName,
				Host:       localExternalHost,
				Status:     true,
				Port:       int64(c.port),
				CreateTime: time.Now(),
			},
			HostList: hostList,
		}

		apisrv := &api.APIServer{
			TypeMeta: unversioned.TypeMeta{
				Kind: "APIServer",
			},
			ObjectMeta: freezerapi.ObjectMeta{
				Name: spec.Server.Name,
			},
			Spec: spec,
		}

		err = c.APIServerRegistry.CreateAPIServer(ctx, apisrv)
		if err != nil {
			glog.Fatalf("Publish api server error:%v\r\n", err.Error())
			return
		}
	}
}
