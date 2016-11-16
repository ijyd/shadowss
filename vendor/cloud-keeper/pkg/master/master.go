package master

import (
	"net/http"
	"time"

	"gofreezer/pkg/api/unversioned"
	"gofreezer/pkg/util/sets"

	"apistack/pkg/genericapiserver"
	"apistack/pkg/registry/generic"

	apiv1 "cloud-keeper/pkg/api/v1"
	batchapiv1alpha1 "cloud-keeper/pkg/apis/batch/v1alpha1"

	// RESTStorage installers
	batchrest "cloud-keeper/pkg/registry/batch/rest"
	corerest "cloud-keeper/pkg/registry/core/rest"

	"github.com/golang/glog"
)

const (
	// DefaultEndpointReconcilerInterval is the default amount of time for how often the endpoints for
	// the kubernetes Service are reconciled.
	DefaultEndpointReconcilerInterval = 10 * time.Second
)

type Config struct {
	GenericConfig *genericapiserver.Config

	StorageFactory          genericapiserver.StorageFactory
	DeleteCollectionWorkers int
	EnableCoreControllers   bool
	EnableWatchCache        bool
	// genericapiserver.RESTStorageProviders provides RESTStorage building methods keyed by groupName
	RESTStorageProviders map[string]genericapiserver.RESTStorageProvider
	// Used to start and monitor tunneling
	Tunneler          genericapiserver.Tunneler
	EnableUISupport   bool
	EnableLogsSupport bool
	ProxyTransport    http.RoundTripper
}

// Master contains state for a Kubernetes cluster master/api server.
type Master struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

type completedConfig struct {
	*Config
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (c *Config) Complete() completedConfig {
	c.GenericConfig.Complete()

	// enable swagger UI only if general UI support is on
	c.GenericConfig.EnableSwaggerUI = c.GenericConfig.EnableSwaggerUI && c.EnableUISupport

	return completedConfig{c}
}

// SkipComplete provides a way to construct a server instance without config completion.
func (c *Config) SkipComplete() completedConfig {
	return completedConfig{c}
}

// New returns a new instance of Master from the given config.
// Certain config fields will be set to a default value if unset.
// Certain config fields must be specified, including:
//   KubeletClientConfig
func (c completedConfig) New() (*Master, error) {

	s, err := c.Config.GenericConfig.SkipComplete().New() // completion is done in Complete, no need for a second time
	if err != nil {
		return nil, err
	}

	// if c.EnableUISupport {
	// 	routes.UIRedirect{}.Install(s.HandlerContainer)
	// }

	m := &Master{
		GenericAPIServer: s,
	}

	restOptionsFactory := restOptionsFactory{
		deleteCollectionWorkers: c.DeleteCollectionWorkers,
		enableGarbageCollection: c.GenericConfig.EnableGarbageCollection,
		storageFactory:          c.StorageFactory,
	}

	c.EnableWatchCache = false
	if c.EnableWatchCache {
		//restOptionsFactory.storageDecorator = registry.StorageWithCacher
	} else {
		restOptionsFactory.storageDecorator = generic.UndecoratedStorage
	}

	//install legacy rest storage
	if c.GenericConfig.APIResourceConfigSource.AnyResourcesForVersionEnabled(apiv1.SchemeGroupVersion) {
		legacyRESTStorageProvider := corerest.LegacyRESTStorageProvider{
			StorageFactory: c.StorageFactory,
		}
		m.InstallLegacyAPI(c.Config, restOptionsFactory.NewFor, legacyRESTStorageProvider)
	}

	// Add some hardcoded storage for now.  Append to the map.
	if c.RESTStorageProviders == nil {
		c.RESTStorageProviders = map[string]genericapiserver.RESTStorageProvider{}
	}
	c.RESTStorageProviders[batchapiv1alpha1.GroupName] = batchrest.RESTStorageProvider{}
	m.InstallAPIs(c.Config, restOptionsFactory.NewFor)

	//m.InstallGeneralEndpoints(c.Config)

	return m, nil
}

func (m *Master) InstallLegacyAPI(c *Config, restOptionsGetter genericapiserver.RESTOptionsGetter, legacyRESTStorageProvider corerest.LegacyRESTStorageProvider) {
	legacyRESTStorage, apiGroupInfo, err := legacyRESTStorageProvider.NewLegacyRESTStorage(restOptionsGetter)
	if err != nil {
		glog.Fatalf("Error building core storage: %v", err)
	}

	if c.EnableCoreControllers {
		bootstrapController := c.NewBootstrapController(legacyRESTStorage, c.GenericConfig.ReadWritePort)
		if err := m.GenericAPIServer.AddPostStartHook("bootstrap-controller", bootstrapController.PostStartHook); err != nil {
			glog.Fatalf("Error registering PostStartHook %q: %v", "bootstrap-controller", err)
		}
	}

	if err := m.GenericAPIServer.InstallLegacyAPIGroup(genericapiserver.DefaultLegacyAPIPrefix, &apiGroupInfo); err != nil {
		glog.Fatalf("Error in registering group versions: %v", err)
	}
}

// TODO this needs to be refactored so we have a way to add general health checks to genericapiserver
// TODO profiling should be generic
// func (m *Master) InstallGeneralEndpoints(c *Config) {
// 	// Run the tunneler.
// 	healthzChecks := []healthz.HealthzChecker{}
// 	if c.Tunneler != nil {
// 		c.Tunneler.Run(m.getNodeAddresses)
// 		healthzChecks = append(healthzChecks, healthz.NamedCheck("SSH Tunnel Check", genericapiserver.TunnelSyncHealthChecker(c.Tunneler)))
// 		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
// 			Name: "apiserver_proxy_tunnel_sync_latency_secs",
// 			Help: "The time since the last successful synchronization of the SSH tunnels for proxy requests.",
// 		}, func() float64 { return float64(c.Tunneler.SecondsSinceSync()) })
// 	}
// 	healthz.InstallHandler(&m.GenericAPIServer.HandlerContainer.NonSwaggerRoutes, healthzChecks...)
//
// 	if c.GenericConfig.EnableProfiling {
// 		routes.MetricsWithReset{}.Install(m.GenericAPIServer.HandlerContainer)
// 	} else {
// 		routes.DefaultMetrics{}.Install(m.GenericAPIServer.HandlerContainer)
// 	}
//
// }

func (m *Master) InstallAPIs(c *Config, restOptionsGetter genericapiserver.RESTOptionsGetter) {
	apiGroupsInfo := []genericapiserver.APIGroupInfo{}

	// stabilize order.
	// TODO find a better way to configure priority of groups
	for _, group := range sets.StringKeySet(c.RESTStorageProviders).List() {
		if !c.GenericConfig.APIResourceConfigSource.AnyResourcesForGroupEnabled(group) {
			glog.V(1).Infof("Skipping disabled API group %q.", group)
			continue
		}
		restStorageBuilder := c.RESTStorageProviders[group]
		apiGroupInfo, enabled := restStorageBuilder.NewRESTStorage(c.GenericConfig.APIResourceConfigSource, restOptionsGetter)
		if !enabled {
			glog.Warningf("Problem initializing API group %q, skipping.", group)
			continue
		}
		glog.V(1).Infof("Enabling API group %q.", group)

		if postHookProvider, ok := restStorageBuilder.(genericapiserver.PostStartHookProvider); ok {
			name, hook, err := postHookProvider.PostStartHook()
			if err != nil {
				glog.Fatalf("Error building PostStartHook: %v", err)
			}
			if err := m.GenericAPIServer.AddPostStartHook(name, hook); err != nil {
				glog.Fatalf("Error registering PostStartHook %q: %v", name, err)
			}
		}

		apiGroupsInfo = append(apiGroupsInfo, apiGroupInfo)
	}

	for i := range apiGroupsInfo {
		if err := m.GenericAPIServer.InstallAPIGroup(&apiGroupsInfo[i]); err != nil {
			glog.Fatalf("Error in registering group versions: %v", err)
		}
	}
}

// func getServersToValidate(storageFactory genericapiserver.StorageFactory) map[string]apiserver.Server {
// 	serversToValidate := map[string]apiserver.Server{
// 		"controller-manager": {Addr: "127.0.0.1", Port: ports.ControllerManagerPort, Path: "/healthz"},
// 		"scheduler":          {Addr: "127.0.0.1", Port: ports.SchedulerPort, Path: "/healthz"},
// 	}
//
// 	for ix, machine := range storageFactory.Backends() {
// 		etcdUrl, err := url.Parse(machine)
// 		if err != nil {
// 			glog.Errorf("Failed to parse etcd url for validation: %v", err)
// 			continue
// 		}
// 		var port int
// 		var addr string
// 		if strings.Contains(etcdUrl.Host, ":") {
// 			var portString string
// 			addr, portString, err = net.SplitHostPort(etcdUrl.Host)
// 			if err != nil {
// 				glog.Errorf("Failed to split host/port: %s (%v)", etcdUrl.Host, err)
// 				continue
// 			}
// 			port, _ = strconv.Atoi(portString)
// 		} else {
// 			addr = etcdUrl.Host
// 			port = 2379
// 		}
// 		// TODO: etcd health checking should be abstracted in the storage tier
// 		serversToValidate[fmt.Sprintf("etcd-%d", ix)] = apiserver.Server{
// 			Addr:        addr,
// 			EnableHTTPS: etcdUrl.Scheme == "https",
// 			Port:        port,
// 			Path:        "/health",
// 			Validate:    etcdutil.EtcdHealthCheck,
// 		}
// 	}
// 	return serversToValidate
// }

type restOptionsFactory struct {
	deleteCollectionWorkers int
	enableGarbageCollection bool
	storageFactory          genericapiserver.StorageFactory
	storageDecorator        generic.StorageDecorator
}

func (f restOptionsFactory) NewFor(resource unversioned.GroupResource) generic.RESTOptions {
	storageConfig, err := f.storageFactory.NewConfig(resource)
	if err != nil {
		glog.Fatalf("Unable to find storage destination for %v, due to %v", resource, err.Error())
	}

	return generic.RESTOptions{
		StorageConfig:           storageConfig,
		Decorator:               f.storageDecorator,
		DeleteCollectionWorkers: f.deleteCollectionWorkers,
		EnableGarbageCollection: f.enableGarbageCollection,
		ResourcePrefix:          f.storageFactory.ResourcePrefix(resource),
	}
}

func DefaultAPIResourceConfigSource() *genericapiserver.ResourceConfig {
	ret := genericapiserver.NewResourceConfig()
	ret.EnableVersions(
		apiv1.SchemeGroupVersion,
		batchapiv1alpha1.SchemeGroupVersion,
	)

	// all extensions resources except these are disabled by default
	// ret.EnableResources(
	// 	extensionsapiv1beta1.SchemeGroupVersion.WithResource(""),
	// 	extensionsapiv1beta1.SchemeGroupVersion.WithResource("batchaccsevers"),
	// )

	return ret
}
