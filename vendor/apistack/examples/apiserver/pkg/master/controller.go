package master

import (
	"apistack/pkg/genericapiserver"
	"apistack/pkg/util/async"

	corerest "apistack/examples/apiserver/pkg/registry/core/rest"
	"apistack/examples/apiserver/pkg/registry/core/user"
)

// // Controller is the controller manager for the core bootstrap cloud-keeper controller
// // loops
type Controller struct {
	UserRegistry user.Registry

	runner *async.Runner
}

//
// // NewBootstrapController returns a controller for watching the core capabilities of the master
func (c *Config) NewBootstrapController(legacyRESTStorage corerest.LegacyRESTStorage) *Controller {
	return &Controller{
		UserRegistry: legacyRESTStorage.UserRegistry,
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

	//not have any runner
	// c.runner = async.NewRunner(c.RunKubernetesNamespaces, c.RunKubernetesService, repairClusterIPs.RunUntil, repairNodePorts.RunUntil)
	// c.runner.Start()
}
