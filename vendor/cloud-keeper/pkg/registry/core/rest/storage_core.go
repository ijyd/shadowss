package rest

import (
	"apistack/pkg/apimachinery/registered"
	"apistack/pkg/genericapiserver"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/api/unversioned"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/account"
	accmysql "cloud-keeper/pkg/registry/core/account/mysql"
	accserveretcd "cloud-keeper/pkg/registry/core/accserver/etcd"
	"cloud-keeper/pkg/registry/core/apiserver"
	apiserveretcd "cloud-keeper/pkg/registry/core/apiserver/etcd"
	"cloud-keeper/pkg/registry/core/login"
	"cloud-keeper/pkg/registry/core/node"
	nodeetcd "cloud-keeper/pkg/registry/core/node/etcd"
	nodemysql "cloud-keeper/pkg/registry/core/node/mysql"
	noderest "cloud-keeper/pkg/registry/core/node/rest"
	"cloud-keeper/pkg/registry/core/nodeuser"
	nodeuseretcd "cloud-keeper/pkg/registry/core/nodeuser/etcd"
	"cloud-keeper/pkg/registry/core/user"
	usermysql "cloud-keeper/pkg/registry/core/user/mysql"
	userrest "cloud-keeper/pkg/registry/core/user/rest"
	"cloud-keeper/pkg/registry/core/userfile"
	"cloud-keeper/pkg/registry/core/userservice"
	userserviceetcd "cloud-keeper/pkg/registry/core/userservice/etcd"
	userservicerest "cloud-keeper/pkg/registry/core/userservice/rest"
	"cloud-keeper/pkg/registry/core/usertoken"
	usertokenmysql "cloud-keeper/pkg/registry/core/usertoken/mysql"
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
	UserRegistry      user.Registry
	TokenRegistry     usertoken.Registry
	APIServerRegistry apiserver.Registry
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

	nodeuserStorage := nodeuseretcd.NewREST(restOptionsGetter(api.Resource("nodeusers")))
	nodeuserRegistry := nodeuser.NewRegistry(nodeuserStorage, nodeuserStorage)

	nodeEtcdStorage := nodeetcd.NewREST(restOptionsGetter(api.Resource("nodes")))
	nodeMysqlStorage := nodemysql.NewREST(restOptionsGetter(api.Resource("nodes")))
	nodeStorage := noderest.NewREST(nodeEtcdStorage, nodeMysqlStorage)
	nodeRegistry := node.NewRegistry(nodeStorage, nodeStorage, nodeStorage)

	userMysqlStorage := usermysql.NewREST(restOptionsGetter(api.Resource("users")))
	userRegistry := user.NewRegistry(userMysqlStorage, userMysqlStorage, userMysqlStorage, userMysqlStorage, userMysqlStorage)

	userserviceEtcdStorage := userserviceetcd.NewREST(restOptionsGetter(api.Resource("userservices")))
	userserviceStorage := userservicerest.NewREST(userserviceEtcdStorage, userRegistry, nodeRegistry, nodeuserRegistry)
	userserviceRegistry := userservice.NewRegistry(userserviceStorage, userserviceStorage, userserviceStorage, userserviceStorage)
	userserviceBindingNodeStorage, userservicePropertiesStorage := userservicerest.NewExtendREST(userserviceRegistry)

	//todo: it is not place here,  disordered resource design
	nodeuserStorage.SetRequireRegistry(userserviceRegistry, nodeRegistry)

	userStorage := userrest.NewREST(userRegistry, userserviceRegistry)
	restStorage.UserRegistry = userRegistry

	nodeBindingUserStorage, nodeAPINodeStorage := noderest.NewExtendREST(nodeRegistry, userserviceRegistry)

	tokenStorage := usertokenmysql.NewREST(restOptionsGetter(api.Resource("usertokens")))
	restStorage.TokenRegistry = usertoken.NewRegistry(tokenStorage, tokenStorage, tokenStorage)

	loginStorage := login.NewREST(restStorage.UserRegistry, restStorage.TokenRegistry)

	accStorage := accmysql.NewREST(restOptionsGetter(api.Resource("accounts")))
	accRegistry := account.NewRegistry(accStorage, accStorage, accStorage, accStorage, accStorage)

	accExtendStorage := account.NewREST(accRegistry)

	cloudServerStorage := accserveretcd.NewREST(restOptionsGetter(api.Resource("accservers")), accRegistry)

	userFileStorage := userfile.NewREST()

	apiserverStorage := apiserveretcd.NewREST(restOptionsGetter(api.Resource("apiservers")))
	restStorage.APIServerRegistry = apiserver.NewRegistry(apiserverStorage, apiserverStorage, apiserverStorage.Store, apiserverStorage.Store)

	restStorageMap := map[string]rest.Storage{
		"users":              userStorage,
		"users/bindingnodes": userserviceBindingNodeStorage,
		"users/properties":   userservicePropertiesStorage,

		"logins": loginStorage,

		"accounts":         accStorage,
		"accounts/info":    accExtendStorage.AccInfo,
		"accounts/sshkeys": accExtendStorage.SSHKeys,
		"accounts/exec":    accExtendStorage.Exec,

		"accservers": cloudServerStorage,

		"nodes":               nodeStorage,
		"nodes/bindingusers":  nodeBindingUserStorage,
		"nodes/activeapinode": nodeAPINodeStorage,

		"nodes/nodeusers": nodeuserStorage,

		"userfile":        userFileStorage.File,
		"userfile/stream": userFileStorage.FileStream,
		"userfile/desc":   userFileStorage.FileDesc,

		"apiservers": apiserverStorage,
	}

	apiGroupInfo.VersionedResourcesStorageMap["v1"] = restStorageMap

	return restStorage, apiGroupInfo, nil
}

func (p LegacyRESTStorageProvider) GroupName() string {
	return api.GroupName
}
