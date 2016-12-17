package masterhook

import (
	"apistack/pkg/genericapiserver"
	"apistack/pkg/master"
	"apistack/pkg/util/async"
	freezerapi "gofreezer/pkg/api"
	apierrs "gofreezer/pkg/api/errors"
	"gofreezer/pkg/api/unversioned"
	"gofreezer/pkg/watch"
	"golib/pkg/util/network"
	"golib/pkg/util/wait"
	"strings"
	"time"

	"github.com/golang/glog"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/apiserver"
	"cloud-keeper/pkg/registry/core/node"
	corerest "cloud-keeper/pkg/registry/core/rest"
	"cloud-keeper/pkg/registry/core/user"
)

// // Controller is the controller manager for the core bootstrap cloud-keeper controller
// // loops
type Controller struct {
	UserRegistry      user.Registry
	APIServerRegistry apiserver.Registry
	NodeRegistry      node.Registry
	runner            *async.Runner
	port              int
}

// NewBootstrapController returns a controller for watching the core capabilities of the master
func NewBootstrapController(c *master.Config, legacyRESTStorage corerest.LegacyRESTStorage) *Controller {
	InnerHookHandler.SetRegistry(legacyRESTStorage)
	return &Controller{
		UserRegistry:      legacyRESTStorage.UserRegistry,
		APIServerRegistry: legacyRESTStorage.APIServerRegistry,
		NodeRegistry:      legacyRESTStorage.NodeRegistry,
		port:              c.GenericConfig.ReadWritePort,
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

	c.runner = async.NewRunner(c.MaintainNode)
	c.runner.Start()
}

func (c *Controller) MaintainNode(ch chan struct{}) {
	ctx := freezerapi.NewContext()
	options := &freezerapi.ListOptions{}

	wait.Until(func() {
		for {
			w, err := c.NodeRegistry.WatchNodes(ctx, options)
			if err != nil {
				glog.Errorf("cant watch node\r\n")
				return
			}

			condition := func(event watch.Event) (bool, error) {
				return event.Type == watch.Deleted, nil
			}

			event, err := watch.Until(0, w, condition)
			if err != nil {
				glog.Errorf("watch node event(%v) error:%v", event, err)
				time.Sleep(time.Second * 1)
				continue
			}
			node, ok := event.Object.(*api.Node)
			if !ok {
				glog.Errorf("%#v is not a node event", event)
				continue
			}
			nodeName := node.Name

			userlist, err := c.UserRegistry.ListUserByNodeName(ctx, nodeName, nil)
			if err != nil {
				glog.Errorf("list user by node name(%v) error:%v\r\n", nodeName, err)
				continue
			}

			for _, v := range userlist.Items {
				err = c.UserRegistry.DelNodeFromUser(ctx, &v, nodeName, true, false)
				if err != nil {
					glog.Errorf("delete node(%v) from user(%v)", nodeName, v.Name)
					continue
				}
			}
		}
	}, time.Duration(time.Second*1), ch)

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
