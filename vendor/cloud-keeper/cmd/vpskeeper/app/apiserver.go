package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"

	"cloud-keeper/cmd/vpskeeper/app/options"
	"cloud-keeper/pkg/apiserver"
	"cloud-keeper/pkg/backend"

	"golib/pkg/util/exec"
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
	glog.V(5).Infof("Got license program %v\r\n", program)

	execCom := exec.New()
	cmd := execCom.Command(program, "check")
	out, err := cmd.CombinedOutput()
	if err == exec.ErrExecutableNotFound {
		glog.Errorf("Expected error ErrExecutableNotFound but got %v", err)
		return false
	}

	result := strings.Contains(string(out), "license result true")

	glog.V(5).Infof("check licese resutl %s:%v", string(out), result)
	return result
}

func checkLicense() bool {
	return execLicenseVerify()
}

//Run start a api server
func Run(options *options.ServerOption) error {

	be := backend.NewBackend(options.Storage.Type, options.Storage.ServerList)

	config := apiserver.Config{
		InsecurePort:       options.InsecurePort,
		SecurePort:         options.SecurePort,
		StorageClient:      be,
		SwaggerPath:        options.SwaggerPath,
		TLSCertFile:        options.TLSCertFile,
		TLSPrivateKeyFile:  options.TLSPrivateKeyFile,
		EtcdStorageOptions: options.EtcdStorageConfig,
	}

	if checkLicense() == false {
		return fmt.Errorf("not allow on this server, please contact administrator")
	}

	serverHandler := apiserver.NewApiServer(config)
	if serverHandler == nil {
		return fmt.Errorf("api server init failure")
	}

	if err := serverHandler.Run(); err != nil {
		return err
	}

	return nil
}
