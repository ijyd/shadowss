package rest

import (
	"cloud-keeper/pkg/apis/abacpolicys"
	abacv1beta1 "cloud-keeper/pkg/apis/abacpolicys/v1beta1"

	"apistack/pkg/genericapiserver"
	"cloud-keeper/pkg/registry/abacpolicys/role"

	"gofreezer/pkg/api/rest"
)

type RESTStorageProvider struct {
	policyFilePath string
}

var _ genericapiserver.RESTStorageProvider = &RESTStorageProvider{}

func (p RESTStorageProvider) NewRESTStorage(apiResourceConfigSource genericapiserver.APIResourceConfigSource, restOptionsGetter genericapiserver.RESTOptionsGetter) (genericapiserver.APIGroupInfo, bool) {
	// TODO figure out how to make the swagger generation stable, while allowing this endpoint to be disabled.
	// if p.Authenticator == nil {
	// 	return genericapiserver.APIGroupInfo{}, false
	// }

	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(abacpolicys.GroupName)

	if apiResourceConfigSource.AnyResourcesForVersionEnabled(abacv1beta1.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[abacv1beta1.SchemeGroupVersion.Version] = p.v1beta1Storage(apiResourceConfigSource, restOptionsGetter)
		apiGroupInfo.GroupMeta.GroupVersion = abacv1beta1.SchemeGroupVersion
	}

	return apiGroupInfo, true
}

func (p RESTStorageProvider) v1beta1Storage(apiResourceConfigSource genericapiserver.APIResourceConfigSource, restOptionsGetter genericapiserver.RESTOptionsGetter) map[string]rest.Storage {
	version := abacv1beta1.SchemeGroupVersion

	storage := map[string]rest.Storage{}
	if apiResourceConfigSource.AnyResourcesForVersionEnabled(abacv1beta1.SchemeGroupVersion) {
		if apiResourceConfigSource.ResourceEnabled(version.WithResource("policys")) {
			abacStorage := role.NewREST("")
			storage["policys"] = abacStorage
		}
	}

	return storage
}

func (p RESTStorageProvider) GroupName() string {
	return abacpolicys.GroupName
}
