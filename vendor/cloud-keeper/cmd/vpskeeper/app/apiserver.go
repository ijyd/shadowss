package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"apistack/pkg/apiserver/authenticator"
	"apistack/pkg/genericapiserver"
	"apistack/pkg/genericapiserver/authorizer"
	genericoptions "apistack/pkg/genericapiserver/options"
	"apistack/pkg/master"
	"apistack/pkg/version"
	"gofreezer/pkg/api"
	utilerrors "gofreezer/pkg/util/errors"
	"gofreezer/pkg/util/exec"
	"gofreezer/pkg/util/wait"

	"cloud-keeper/cmd/vpskeeper/app/options"
	"cloud-keeper/pkg/masterhook"

	"github.com/golang/glog"
)

const (
	licenseProgram = "vpslicense"
)

func execLicenseVerify() bool {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	program := dir + "/" + licenseProgram

	execCom := exec.New()
	cmd := execCom.Command(program, "check")
	out, err := cmd.CombinedOutput()
	if err == exec.ErrExecutableNotFound {
		glog.Errorf("Expected error ErrExecutableNotFound but got %v", err)
		return false
	}

	result := strings.Contains(string(out), "license result true")

	return result
}

func checkLicense() bool {
	return execLicenseVerify()
}

// Run runs the specified APIServer.  This should never exit.
func Run(s *options.ServerOption) error {
	if checkLicense() == false {
		return fmt.Errorf("not allow on this server, please contact administrator")
	}

	if errs := s.Storage.Validate(); len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}
	if err := s.GenericServerRunOptions.DefaultExternalAddress(s.SecureServing, s.InsecureServing); err != nil {
		return err
	}

	genericapiserver.DefaultAndValidateRunOptions(s.GenericServerRunOptions)
	genericConfig := genericapiserver.NewConfig(). // create the new config
							ApplyOptions(s.GenericServerRunOptions). // apply the options selected
							ApplySecureServingOptions(s.SecureServing).
							ApplyInsecureServingOptions(s.InsecureServing).
							ApplyAuthenticationOptions(s.Authentication)

	if err := genericConfig.MaybeGenerateServingCerts(); err != nil {
		glog.Fatalf("Failed to generate service certificate: %v", err)
	}

	keeperVersion := version.Get()

	storageGroupsToEncodingVersion, err := s.GenericServerRunOptions.StorageGroupsToEncodingVersion()
	if err != nil {
		glog.Fatalf("error generating storage version map: %s", err)
	}
	storageFactory, err := genericapiserver.BuildDefaultStorageFactory(
		s.Storage.StorageConfig, s.GenericServerRunOptions.DefaultStorageMediaType, api.Codecs,
		genericapiserver.NewDefaultResourceEncodingConfig(), storageGroupsToEncodingVersion,
		// FIXME: this GroupVersionResource override should be configurable
		nil,
		masterhook.DefaultAPIResourceConfigSource(), s.GenericServerRunOptions.RuntimeConfig)

	authenticatorConfig := s.Authentication.ToAuthenticationConfig(s.SecureServing.ClientCA)
	//if allow innerhook append our innerhook
	if s.Authentication.InnerHook.Allow {
		authenticatorConfig.InnerHookFunc = masterhook.InnerHookHandler.AuthenticateTokenInnerHook
	}
	apiAuthenticator, securityDefinitions, err := authenticator.New(authenticatorConfig)
	if err != nil {
		glog.Fatalf("Invalid Authentication config: %v", err)
	}

	authorizationConfig := s.Authorization.ToAuthorizationConfig(nil)
	apiAuthorizer, err := authorizer.NewAuthorizerFromAuthorizationConfig(authorizationConfig)
	if err != nil {
		glog.Fatalf("Invalid Authorization Config: %v", err)
	}

	// TODO(dims): We probably need to add an option "EnableLoopbackToken"
	privilegedLoopbackToken := "49acafe7e63682e1e6b6983580c4ee56" //uuid.NewRandom().String()
	selfClientConfig, err := genericoptions.NewSelfClientConfig(s.SecureServing, s.InsecureServing, privilegedLoopbackToken)
	if err != nil {
		glog.Fatalf("Failed to create clientset: %v", err)
	}

	genericConfig.Version = &keeperVersion
	genericConfig.Authenticator = apiAuthenticator
	genericConfig.Authorizer = apiAuthorizer
	genericConfig.LoopbackClientConfig = selfClientConfig
	genericConfig.APIResourceConfigSource = storageFactory.APIResourceConfigSource
	//genericConfig.EnableMetrics = false
	genericConfig.EnableOpenAPISupport = false
	genericConfig.OpenAPIConfig.SecurityDefinitions = securityDefinitions

	config := &master.Config{
		GenericConfig: genericConfig,

		StorageFactory:          storageFactory,
		EnableWatchCache:        s.GenericServerRunOptions.EnableWatchCache,
		DeleteCollectionWorkers: s.GenericServerRunOptions.DeleteCollectionWorkers,
		EnableUISupport:         false,
		EnableLogsSupport:       true,
	}

	m, err := config.ExtraComplete(masterhook.InstallLegacyAPI, masterhook.InstallAPIs).New()
	if err != nil {
		return err
	}

	m.GenericAPIServer.PrepareRun().Run(wait.NeverStop)
	return nil
}
