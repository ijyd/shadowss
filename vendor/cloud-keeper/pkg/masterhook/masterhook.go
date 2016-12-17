package masterhook

import (
	"apistack/pkg/genericapiserver"
	"apistack/pkg/master"

	apiv1 "cloud-keeper/pkg/api/v1"
	abacbeta1 "cloud-keeper/pkg/apis/abacpolicys/v1beta1"
	batchapiv1alpha1 "cloud-keeper/pkg/apis/batch/v1alpha1"

	abacrest "cloud-keeper/pkg/registry/abacpolicys/rest"
	batchrest "cloud-keeper/pkg/registry/batch/rest"
	corerest "cloud-keeper/pkg/registry/core/rest"

	"github.com/golang/glog"
)

func InstallLegacyAPI(c *master.Config,
	restOptionsGetter genericapiserver.RESTOptionsGetter,
	storageFactory genericapiserver.StorageFactory,
	genericAPIServer *genericapiserver.GenericAPIServer) {

	legacyRESTStorageProvider := corerest.LegacyRESTStorageProvider{
		StorageFactory: c.StorageFactory,
	}

	legacyRESTStorage, apiGroupInfo, err := legacyRESTStorageProvider.NewLegacyRESTStorage(restOptionsGetter)
	if err != nil {
		glog.Fatalf("Error building core storage: %v", err)
	}

	bootstrapController := NewBootstrapController(c, legacyRESTStorage)
	if err := genericAPIServer.AddPostStartHook("bootstrap-controller", bootstrapController.PostStartHook); err != nil {
		glog.Fatalf("Error registering PostStartHook %q: %v", "bootstrap-controller", err)
	}

	InnerHookHandler.SetRegistry(legacyRESTStorage)

	if err := master.InstallLegacyAPI(genericAPIServer, restOptionsGetter, &apiGroupInfo); err != nil {
		glog.Fatalf("Error in registering group versions: %v", err)
	}

}

func InstallAPIs(apiResourceConfigSource genericapiserver.APIResourceConfigSource,
	restOptionsGetter genericapiserver.RESTOptionsGetter,
	genericAPIServer *genericapiserver.GenericAPIServer) {

	restStorageProviders := []genericapiserver.RESTStorageProvider{
		batchrest.RESTStorageProvider{},
		abacrest.RESTStorageProvider{},
	}

	postHook := func(restStorageBuilder genericapiserver.RESTStorageProvider) {
		master.InstallAPIsSimpleHook(genericAPIServer, restStorageBuilder)
	}

	master.InstallAPIs(genericAPIServer, apiResourceConfigSource, restOptionsGetter, postHook, restStorageProviders...)

}

func DefaultAPIResourceConfigSource() *genericapiserver.ResourceConfig {
	ret := genericapiserver.NewResourceConfig()
	ret.EnableVersions(
		apiv1.SchemeGroupVersion,
		batchapiv1alpha1.SchemeGroupVersion,
		abacbeta1.SchemeGroupVersion,
	)

	return ret
}
