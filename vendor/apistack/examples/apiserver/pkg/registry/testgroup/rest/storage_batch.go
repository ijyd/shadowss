package rest

import (
	"apistack/examples/apiserver/pkg/apis/testgroup"
	testgroupv1 "apistack/examples/apiserver/pkg/apis/testgroup/v1"
	testtyperest "apistack/examples/apiserver/pkg/registry/testgroup/testtype/etcd"

	"gofreezer/pkg/api/rest"

	"apistack/pkg/genericapiserver"
)

type RESTStorageProvider struct {
}

var _ genericapiserver.RESTStorageProvider = &RESTStorageProvider{}

func (p RESTStorageProvider) NewRESTStorage(apiResourceConfigSource genericapiserver.APIResourceConfigSource, restOptionsGetter genericapiserver.RESTOptionsGetter) (genericapiserver.APIGroupInfo, bool) {
	// TODO figure out how to make the swagger generation stable, while allowing this endpoint to be disabled.
	// if p.Authenticator == nil {
	// 	return genericapiserver.APIGroupInfo{}, false
	// }

	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(testgroup.GroupName)

	if apiResourceConfigSource.AnyResourcesForVersionEnabled(testgroupv1.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[testgroupv1.SchemeGroupVersion.Version] = p.v1beta1Storage(apiResourceConfigSource, restOptionsGetter)
		apiGroupInfo.GroupMeta.GroupVersion = testgroupv1.SchemeGroupVersion
	}

	return apiGroupInfo, true
}

func (p RESTStorageProvider) v1beta1Storage(apiResourceConfigSource genericapiserver.APIResourceConfigSource, restOptionsGetter genericapiserver.RESTOptionsGetter) map[string]rest.Storage {
	version := testgroupv1.SchemeGroupVersion

	storage := map[string]rest.Storage{}
	if apiResourceConfigSource.AnyResourcesForVersionEnabled(testgroupv1.SchemeGroupVersion) {
		if apiResourceConfigSource.ResourceEnabled(version.WithResource("testtypes")) {
			testtypestorage := testtyperest.NewREST(restOptionsGetter(testgroup.Resource("testtypes")))
			storage["testtypes"] = testtypestorage
		}
	}

	return storage
}

func (p RESTStorageProvider) GroupName() string {
	return testgroup.GroupName
}
