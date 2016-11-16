package ansible

import (
	api "cloud-keeper/pkg/api"
	"fmt"
	"golib/pkg/util/exec"
	"log"
	"os"
	"path/filepath"

	"github.com/golang/glog"
)

var absdir string

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	absdir = dir + "/"
}

func DeployShadowss(typ api.OperatorType, hosts []string, sshkey string, attr map[string]string) error {

	var hostpath, sshkeypath, playbook, attrpath, keypem string
	switch typ {
	case api.OperatorVultr:
		hostpath = absdir + ansibleVulHostFile
		sshkeypath = absdir + ansibleVulSSHKeyFile
		playbook = absdir + ansibleVulDeploySSPlayBook
		attrpath = absdir + ansibleVulAttrFile
		keypem = absdir + ansibleVulPrivateKey
	case api.OperatorDigitalOcean:
		hostpath = absdir + ansibleDGOCHostFile
		sshkeypath = absdir + ansibleDGOCSSHKeyFile
		attrpath = absdir + ansibleDGOCAttrFile
		playbook = absdir + ansibleDGOCDeploySSPlayBook
		keypem = absdir + ansibleDGOCPrivateKey
	}

	err := WriteSSAttrFile(attrpath, attr)
	if err != nil {
		glog.Errorf("Unexpected error write attr to file %v", err)
		return err
	}

	WriteDeplossConfigFile(hosts, privateKey, hostpath, sshkeypath, keypem)

	execCom := exec.New()
	cmd := execCom.Command("ansible-playbook", "-i", hostpath, playbook)
	out, err := cmd.CombinedOutput()
	if err == exec.ErrExecutableNotFound {
		glog.Errorf("Expected error ErrExecutableNotFound but got %v", err)
		return err
	}
	glog.V(5).Infof("playbook result %s", string(out))
	return err
}

func RestartShadowss(typ api.OperatorType, hosts []string, sshkey string) error {

	var hostpath, sshkeypath, playbook, keypem string
	switch typ {
	case api.OperatorVultr:
		hostpath = absdir + ansibleVulHostFile
		sshkeypath = absdir + ansibleVulSSHKeyFile
		playbook = absdir + ansibleVulRestartSSPlayBook
		keypem = absdir + ansibleVulPrivateKey
	case api.OperatorDigitalOcean:
		hostpath = absdir + ansibleDGOCHostFile
		sshkeypath = absdir + ansibleDGOCSSHKeyFile
		playbook = absdir + ansibleDGOCRestartSSPlayBook
		keypem = absdir + ansibleDGOCPrivateKey
	}

	WriteDeplossConfigFile(hosts, privateKey, hostpath, sshkeypath, keypem)

	execCom := exec.New()
	cmd := execCom.Command("ansible-playbook", "-i", hostpath, playbook, "--tags", "restartSS")
	out, err := cmd.CombinedOutput()
	if err == exec.ErrExecutableNotFound {
		glog.Errorf("Expected error ErrExecutableNotFound but got %v", err)
		return err
	}
	glog.V(5).Infof("playbook result %s", string(out))
	return err
}

func DeployVPS(typ api.OperatorType, srv *api.AccServer, key string) error {

	var playbook, vpsVar, varfile string

	switch typ {
	case api.OperatorVultr:
		playbook = absdir + ansibleVulDeployVPSPlayBook
		varfile = absdir + ansibleVulVarFile

		vpsVar = fmt.Sprintf("api_key: %s\r\n", key)
		vpsVar += fmt.Sprintf("DCID: %s\r\n", srv.Spec.Region)
		vpsVar += fmt.Sprintf("OSID: %s\r\n", srv.Spec.Image)
		vpsVar += fmt.Sprintf("SSHKEYID: %q\r\n", srv.Spec.SSHKeyID)
		vpsVar += fmt.Sprintf("VPSPLANID: %s\r\n", srv.Spec.Size)
		vpsVar += fmt.Sprintf("name: %s\r\n", srv.Name)

	case api.OperatorDigitalOcean:
		playbook = absdir + ansibleDGOCDeployVPSPlayBook
		varfile = absdir + ansibleDGOCVarFile

		vpsVar = fmt.Sprintf("do_token: %s\r\n", key)
		vpsVar += fmt.Sprintf("image_id: %s\r\n", srv.Spec.Image)
		vpsVar += fmt.Sprintf("ssh_key_ids: %s\r\n", srv.Spec.SSHKeyID)
		vpsVar += fmt.Sprintf("region: %s\r\n", srv.Spec.Region)
		vpsVar += fmt.Sprintf("size_id: %s\r\n", srv.Spec.Size)
		vpsVar += fmt.Sprintf("name: %s\r\n", srv.Spec.Name)
		vpsVar += fmt.Sprintf("timeout: 500\r\n")

	}

	err := WriteCreateVPSVarFile([]byte(vpsVar), varfile)
	if err != nil {
		return err
	}

	execCom := exec.New()
	cmd := execCom.Command("ansible-playbook", playbook)
	out, err := cmd.CombinedOutput()
	if err == exec.ErrExecutableNotFound {
		glog.Errorf("Expected error ErrExecutableNotFound but got %v", err)
		return err
	}
	glog.V(5).Infof("playbook result %s", string(out))

	return err
}

func DeleteVPS(typ api.OperatorType, id int64, key string) error {

	var playbook, varData, varfile string
	switch typ {
	case api.OperatorVultr:
		playbook = absdir + ansibleVulDeleteVPSPlayBook
		varfile = absdir + ansibleVulVarFile

		//vultr not delete into here
	case api.OperatorDigitalOcean:
		playbook = absdir + ansibleDGOCDeleteVPSPlayBook
		varfile = absdir + ansibleDGOCVarFile

		varData = fmt.Sprintf("do_token: %q\r\n", key)
		varData += fmt.Sprintf("id: %d\r\n", id)

	}

	err := WriteDeleteVPSVarFile([]byte(varData), varfile)
	if err != nil {
		return err
	}

	execCom := exec.New()
	cmd := execCom.Command("ansible-playbook", playbook)
	out, err := cmd.CombinedOutput()
	if err == exec.ErrExecutableNotFound {
		glog.Errorf("Expected error ErrExecutableNotFound but got %v", err)
		return err
	}
	glog.V(5).Infof("playbook result %s", string(out))

	return err
}
