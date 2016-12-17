package rest

import (
	"apistack/pkg/apimachinery/registered"
	"apistack/pkg/genericapiserver"
	"gofreezer/pkg/api/rest"
	"gofreezer/pkg/runtime/schema"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/registry/core/account"
	"cloud-keeper/pkg/registry/core/account/accserver"
	accserverrest "cloud-keeper/pkg/registry/core/account/accserver/rest"
	accmysql "cloud-keeper/pkg/registry/core/account/mysql"
	accrest "cloud-keeper/pkg/registry/core/account/rest"
	"cloud-keeper/pkg/registry/core/apiserver"
	apiserveretcd "cloud-keeper/pkg/registry/core/apiserver/etcd"
	batchshadowssrest "cloud-keeper/pkg/registry/core/batchshadowss"
	batchusers "cloud-keeper/pkg/registry/core/batchusers"
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
	usertokenetcd "cloud-keeper/pkg/registry/core/usertoken/etcd"
	//usertokenmysql "cloud-keeper/pkg/registry/core/usertoken/mysql"
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
	NodeRegistry      node.Registry
}

func (c LegacyRESTStorageProvider) NewLegacyRESTStorage(restOptionsGetter genericapiserver.RESTOptionsGetter) (LegacyRESTStorage, genericapiserver.APIGroupInfo, error) {
	apiGroupInfo := genericapiserver.APIGroupInfo{
		GroupMeta:                    *registered.GroupOrDie(api.GroupName),
		VersionedResourcesStorageMap: map[string]map[string]rest.Storage{},
		Scheme:                      api.Scheme,
		ParameterCodec:              api.ParameterCodec,
		NegotiatedSerializer:        api.Codecs,
		SubresourceGroupVersionKind: map[string]schema.GroupVersionKind{},
	}

	restStorage := LegacyRESTStorage{}

	nodeuserStorage := nodeuseretcd.NewREST(restOptionsGetter(api.Resource("nodeusers")))
	nodeuserRegistry := nodeuser.NewRegistry(nodeuserStorage, nodeuserStorage)

	nodeEtcdStorage := nodeetcd.NewREST(restOptionsGetter(api.Resource("nodes")))
	nodeMysqlStorage := nodemysql.NewREST(restOptionsGetter(api.Resource("nodes")))
	nodeStorage := noderest.NewREST(nodeEtcdStorage, nodeMysqlStorage)
	nodeRegistry := node.NewRegistry(nodeStorage, nodeStorage, nodeStorage, nodeStorage)

	userMysqlStorage := usermysql.NewREST(restOptionsGetter(api.Resource("users")))
	userEtcdStorage := useretcd.NewREST(restOptionsGetter(api.Resource("users")))
	userDynamoStorage := userdynamo.NewREST(restOptionsGetter(api.Resource("users")))
	userStorage := userrest.NewREST(userEtcdStorage, userMysqlStorage, userDynamoStorage, nodeRegistry, nodeuserRegistry)
	userRegistry := user.NewRegistry(userStorage, userStorage, userStorage, userStorage, userStorage, userStorage, userStorage)

	userBindingNodeStorage, userPropertiesStorage, userActivation := userrest.NewExtendREST(userRegistry)

	//todo: it is not place here,  disordered resource design
	nodeuserStorage.SetRequireRegistry(nodeRegistry, userRegistry)
	nodeStorage.SetRequireRegistry(userRegistry, nodeuserRegistry)

	nodeAPINodeStorage, nodeRefreshStorage, nodeUserStorage := noderest.NewExtendREST(nodeRegistry, userRegistry)

	//tokenStorage := usertokenmysql.NewREST(restOptionsGetter(api.Resource("usertokens")))
	tokenStorage := usertokenetcd.NewREST(restOptionsGetter(api.Resource("usertokens")))
	tokenRegistry := usertoken.NewRegistry(tokenStorage, tokenStorage, tokenStorage)

	loginStorage := login.NewREST(userRegistry, tokenRegistry)

	accMysqlStorage := accmysql.NewREST(restOptionsGetter(api.Resource("accounts")))
	accStorage := accrest.NewREST(accMysqlStorage)
	accRegistry := account.NewRegistry(accStorage, accStorage, accStorage, accStorage, accStorage)

	accExtendStorage := accrest.NewExtendREST(accRegistry)

	cloudServerStorage := accserverrest.NewREST(accRegistry)
	cloudServerRegistry := accserver.NewRegistry(cloudServerStorage)

	userFileStorage := userfile.NewREST()

	apiserverStorage := apiserveretcd.NewREST(restOptionsGetter(api.Resource("apiservers")))
	apiServerRegistry := apiserver.NewRegistry(apiserverStorage, apiserverStorage, apiserverStorage.Store, apiserverStorage.Store)

	batchUsersStorage := batchusers.NewREST(userRegistry)

	batchShadowssStorage := batchshadowssrest.NewREST(accRegistry, cloudServerRegistry)

	restStorage.UserRegistry = userRegistry
	restStorage.TokenRegistry = tokenRegistry
	restStorage.APIServerRegistry = apiServerRegistry
	restStorage.NodeRegistry = nodeRegistry

	restStorageMap := map[string]rest.Storage{
		"users":              userStorage,
		"users/bindingnodes": userBindingNodeStorage,
		"users/properties":   userPropertiesStorage,
		"users/activation":   userActivation,

		"logins": loginStorage,

		"accounts":         accStorage,
		"accounts/info":    accExtendStorage.AccInfo,
		"accounts/sshkeys": accExtendStorage.SSHKeys,
		"accounts/exec":    accExtendStorage.Exec,
		"accounts/servers": cloudServerStorage,

		// "acc": cloudServerStorage,

		"nodes":           nodeStorage,
		"nodes/refresh":   nodeRefreshStorage,
		"nodes/nodeusers": nodeUserStorage,

		"activeapinode": nodeAPINodeStorage,

		"nodeusers": nodeuserStorage,

		"userfile":        userFileStorage.File,
		"userfile/stream": userFileStorage.FileStream,
		"userfile/desc":   userFileStorage.FileDesc,

		"apiservers": apiserverStorage,

		"batchusers": batchUsersStorage,

		"batchshadowss": batchShadowssStorage,
	}

	apiGroupInfo.VersionedResourcesStorageMap["v1"] = restStorageMap

	return restStorage, apiGroupInfo, nil
}

func (p LegacyRESTStorageProvider) GroupName() string {
	return api.GroupName
}
