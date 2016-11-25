package rest

import (
	"apistack/pkg/apimachinery/registered"
	"apistack/pkg/genericapiserver"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/api/unversioned"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/account"
	accserverrest "cloud-keeper/pkg/registry/core/account/accserver/rest"
	accmysql "cloud-keeper/pkg/registry/core/account/mysql"
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
	userdynamo "cloud-keeper/pkg/registry/core/user/dynamodb"
	useretcd "cloud-keeper/pkg/registry/core/user/etcd"
	usermysql "cloud-keeper/pkg/registry/core/user/mysql"
	userrest "cloud-keeper/pkg/registry/core/user/rest"
	"cloud-keeper/pkg/registry/core/userfile"
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
	userEtcdStorage := useretcd.NewREST(restOptionsGetter(api.Resource("users")))
	userDynamoStorage := userdynamo.NewREST(restOptionsGetter(api.Resource("users")))
	userStorage := userrest.NewREST(userEtcdStorage, userMysqlStorage, userDynamoStorage, nodeRegistry, nodeuserRegistry)
	userRegistry := user.NewRegistry(userStorage, userStorage, userStorage, userStorage, userStorage, userStorage, userStorage)

	userBindingNodeStorage, userPropertiesStorage := userrest.NewExtendREST(userRegistry)

	//todo: it is not place here,  disordered resource design
	nodeuserStorage.SetRequireRegistry(nodeRegistry, userRegistry)
	nodeStorage.SetRequireRegistry(userRegistry, nodeuserRegistry)

	nodeAPINodeStorage, nodeRefreshStorage := noderest.NewExtendREST(nodeRegistry)

	tokenStorage := usertokenmysql.NewREST(restOptionsGetter(api.Resource("usertokens")))
	tokenRegistry := usertoken.NewRegistry(tokenStorage, tokenStorage, tokenStorage)

	loginStorage := login.NewREST(userRegistry, tokenRegistry)

	accStorage := accmysql.NewREST(restOptionsGetter(api.Resource("accounts")))
	accRegistry := account.NewRegistry(accStorage, accStorage, accStorage, accStorage, accStorage)

	accExtendStorage := account.NewREST(accRegistry)

	cloudServerStorage := accserverrest.NewREST(accRegistry)

	userFileStorage := userfile.NewREST()

	apiserverStorage := apiserveretcd.NewREST(restOptionsGetter(api.Resource("apiservers")))
	apiServerRegistry := apiserver.NewRegistry(apiserverStorage, apiserverStorage, apiserverStorage.Store, apiserverStorage.Store)

	restStorage.UserRegistry = userRegistry
	restStorage.TokenRegistry = tokenRegistry
	restStorage.APIServerRegistry = apiServerRegistry

	restStorageMap := map[string]rest.Storage{
		"users":              userStorage,
		"users/bindingnodes": userBindingNodeStorage,
		"users/properties":   userPropertiesStorage,

		"logins": loginStorage,

		"accounts":         accStorage,
		"accounts/info":    accExtendStorage.AccInfo,
		"accounts/sshkeys": accExtendStorage.SSHKeys,
		"accounts/exec":    accExtendStorage.Exec,
		"accounts/servers": cloudServerStorage,

		// "acc": cloudServerStorage,

		"nodes":         nodeStorage,
		"nodes/refresh": nodeRefreshStorage,

		"activeapinode": nodeAPINodeStorage,

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
