package app

// Run runs the specified APIServer.  This should never exit.
import (
	"apistack/pkg/apiserver/authenticator"
	"apistack/pkg/genericapiserver"
	"apistack/pkg/genericapiserver/authorizer"
	genericvalidation "apistack/pkg/genericapiserver/validation"
	"apistack/pkg/version"
	authenticatorunion "apistack/plugin/pkg/auth/authenticator/request/union"
	"database/sql/driver"

	"gofreezer/pkg/api"
	"gofreezer/pkg/api/unversioned"
	"gofreezer/pkg/auth/user"
	"gofreezer/pkg/util/wait"

	"apistack/examples/apiserver/cmd/server/app/options"
	"apistack/examples/apiserver/pkg/master"

	"github.com/golang/glog"
	"github.com/pborman/uuid"
)

func Testting() {
	var t unversioned.Time
	var tInface interface{}
	tInface = t
	_, ok := tInface.(driver.Valuer)

	glog.V(5).Infof("assert ok:%v", ok)
}

func Run(s *options.ServerOption) error {
	Testting()

	genericvalidation.VerifyEtcdServersList(s.GenericServerRunOptions)
	genericvalidation.VerifyMysqlServersList(s.GenericServerRunOptions)
	genericapiserver.DefaultAndValidateRunOptions(s.GenericServerRunOptions)
	genericConfig := genericapiserver.NewConfig(). // create the new config
							ApplyOptions(s.GenericServerRunOptions). // apply the options selected
							Complete()                               // set default values based on the known values

	if err := genericConfig.MaybeGenerateServingCerts(); err != nil {
		glog.Fatalf("Failed to generate service certificate: %v", err)
	}

	keeperVersion := version.Get()

	storageGroupsToEncodingVersion, err := s.GenericServerRunOptions.StorageGroupsToEncodingVersion()
	if err != nil {
		glog.Fatalf("error generating storage version map: %s", err)
	}
	storageFactory, err := genericapiserver.BuildDefaultStorageFactory(
		s.GenericServerRunOptions.StorageConfig, s.GenericServerRunOptions.DefaultStorageMediaType, api.Codecs,
		genericapiserver.NewDefaultResourceEncodingConfig(), storageGroupsToEncodingVersion,
		// FIXME: this GroupVersionResource override should be configurable
		nil,
		master.DefaultAPIResourceConfigSource(), s.GenericServerRunOptions.RuntimeConfig)

	apiAuthenticator, securityDefinitions, err := authenticator.New(authenticator.AuthenticatorConfig{
		Anonymous:           s.GenericServerRunOptions.AnonymousAuth,
		AnyToken:            s.GenericServerRunOptions.EnableAnyToken,
		BasicAuthFile:       s.GenericServerRunOptions.BasicAuthFile,
		ClientCAFile:        s.GenericServerRunOptions.ClientCAFile,
		TokenAuthFile:       s.GenericServerRunOptions.TokenAuthFile,
		OIDCIssuerURL:       s.GenericServerRunOptions.OIDCIssuerURL,
		OIDCClientID:        s.GenericServerRunOptions.OIDCClientID,
		OIDCCAFile:          s.GenericServerRunOptions.OIDCCAFile,
		OIDCUsernameClaim:   s.GenericServerRunOptions.OIDCUsernameClaim,
		OIDCGroupsClaim:     s.GenericServerRunOptions.OIDCGroupsClaim,
		KeystoneURL:         s.GenericServerRunOptions.KeystoneURL,
		RequestHeaderConfig: s.GenericServerRunOptions.AuthenticationRequestHeaderConfig(),
		InnerHookFunc:       master.InnerHookHandler.AuthenticateTokenInnerHook,
	})

	apiAuthorizer := authorizer.NewAlwaysAllowAuthorizer()

	// TODO(dims): We probably need to add an option "EnableLoopbackToken"
	privilegedLoopbackToken := uuid.NewRandom().String()
	if apiAuthenticator != nil {
		var uid = uuid.NewRandom().String()
		tokens := make(map[string]*user.DefaultInfo)
		tokens[privilegedLoopbackToken] = &user.DefaultInfo{
			Name:   user.APIServerUser,
			UID:    uid,
			Groups: []string{user.SystemPrivilegedGroup},
		}

		//append system default user in token authenticator
		tokenAuthenticator := authenticator.NewAuthenticatorFromTokens(tokens)
		apiAuthenticator = authenticatorunion.New(tokenAuthenticator, apiAuthenticator)

		// tokenAuthorizer := authorizer.NewPrivilegedGroups(user.SystemPrivilegedGroup)
		// apiAuthorizer = authorizerunion.New(tokenAuthorizer, apiAuthorizer)
	}

	genericConfig.Version = &keeperVersion
	genericConfig.Authenticator = apiAuthenticator
	genericConfig.Authorizer = apiAuthorizer
	genericConfig.APIResourceConfigSource = storageFactory.APIResourceConfigSource
	//genericConfig.OpenAPIConfig.Info.Title = "Keeper"
	//genericConfig.OpenAPIConfig.Definitions = generatedopenapi.OpenAPIDefinitions
	genericConfig.EnableOpenAPISupport = false
	//genericConfig.EnableMetrics = false
	genericConfig.OpenAPIConfig.SecurityDefinitions = securityDefinitions

	config := &master.Config{
		GenericConfig: genericConfig.Config,

		StorageFactory:          storageFactory,
		EnableWatchCache:        s.GenericServerRunOptions.EnableWatchCache,
		EnableCoreControllers:   true,
		DeleteCollectionWorkers: s.GenericServerRunOptions.DeleteCollectionWorkers,
		EnableUISupport:         false,
		EnableLogsSupport:       true,
	}

	m, err := config.Complete().New()
	if err != nil {
		return err
	}

	m.GenericAPIServer.PrepareRun().Run(wait.NeverStop)
	return nil
}
