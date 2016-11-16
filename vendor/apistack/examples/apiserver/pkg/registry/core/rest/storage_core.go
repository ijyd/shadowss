package rest

import (
	"apistack/pkg/apimachinery/registered"
	"apistack/pkg/genericapiserver"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/api/unversioned"

	"apistack/examples/apiserver/pkg/api"
	"apistack/examples/apiserver/pkg/registry/core/login"
	"apistack/examples/apiserver/pkg/registry/core/user"
	usermysql "apistack/examples/apiserver/pkg/registry/core/user/mysql"
	"apistack/examples/apiserver/pkg/registry/core/usertoken"
	usertokenmysql "apistack/examples/apiserver/pkg/registry/core/usertoken/mysql"
)

// LegacyRESTStorageProvider provides information needed to build RESTStorage for core, but
// does NOT implement the "normal" RESTStorageProvider (yet!)
type LegacyRESTStorageProvider struct {
	StorageFactory genericapiserver.StorageFactory
}

// LegacyRESTStorage returns stateful information about particular instances of REST storage to
// master.go for wiring controllers.
// TODO remove this by running the controller as a poststarthook
type LegacyRESTStorage struct {
	UserRegistry  user.Registry
	TokenRegistry usertoken.Registry
}

func (c LegacyRESTStorageProvider) NewLegacyRESTStorage(restOptionsGetter genericapiserver.RESTOptionsGetter) (LegacyRESTStorage, genericapiserver.APIGroupInfo, error) {
	apiGroupInfo := genericapiserver.APIGroupInfo{
		GroupMeta:                    *registered.GroupOrDie(api.GroupName),
		VersionedResourcesStorageMap: map[string]map[string]rest.Storage{},
		Scheme:                      api.Scheme,
		ParameterCodec:              api.ParameterCodec,
		NegotiatedSerializer:        api.Codecs,
		SubresourceGroupVersionKind: map[string]unversioned.GroupVersionKind{},
	}

	restStorage := LegacyRESTStorage{}

	userStorage := usermysql.NewREST(restOptionsGetter(api.Resource("users")))
	restStorage.UserRegistry = user.NewRegistry(userStorage)

	tokenStorage := usertokenmysql.NewREST(restOptionsGetter(api.Resource("usertokens")))
	restStorage.TokenRegistry = usertoken.NewRegistry(tokenStorage, tokenStorage, tokenStorage)

	loginStorage := login.NewREST(restStorage.UserRegistry, restStorage.TokenRegistry)

	restStorageMap := map[string]rest.Storage{
		"users": userStorage,

		"logins": loginStorage,

		"tokens": tokenStorage,
	}

	apiGroupInfo.VersionedResourcesStorageMap["v1"] = restStorageMap

	return restStorage, apiGroupInfo, nil
}

func (p LegacyRESTStorageProvider) GroupName() string {
	return api.GroupName
}
