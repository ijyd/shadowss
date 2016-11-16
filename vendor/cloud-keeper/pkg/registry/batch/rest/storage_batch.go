package rest

import (
	"cloud-keeper/pkg/apis/batch"
	batchv1alpha1 "cloud-keeper/pkg/apis/batch/v1alpha1"
	"cloud-keeper/pkg/registry/batch/batchserver"

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

	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(batch.GroupName)

	if apiResourceConfigSource.AnyResourcesForVersionEnabled(batchv1alpha1.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[batchv1alpha1.SchemeGroupVersion.Version] = p.v1beta1Storage(apiResourceConfigSource, restOptionsGetter)
		apiGroupInfo.GroupMeta.GroupVersion = batchv1alpha1.SchemeGroupVersion
	}

	return apiGroupInfo, true
}

func (p RESTStorageProvider) v1beta1Storage(apiResourceConfigSource genericapiserver.APIResourceConfigSource, restOptionsGetter genericapiserver.RESTOptionsGetter) map[string]rest.Storage {
	version := batchv1alpha1.SchemeGroupVersion

	storage := map[string]rest.Storage{}
	if apiResourceConfigSource.AnyResourcesForVersionEnabled(batchv1alpha1.SchemeGroupVersion) {
		if apiResourceConfigSource.ResourceEnabled(version.WithResource("batchaccsevers")) {
			batchserverStorage := batchserver.NewREST()
			storage["batchaccsevers"] = batchserverStorage
		}
	}

	return storage
}

func (p RESTStorageProvider) GroupName() string {
	return batch.GroupName
}
