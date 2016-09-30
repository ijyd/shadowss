package ansible

import (
	"cloud-keeper/pkg/api"
	"fmt"
	"golib/pkg/util/exec"

	"github.com/golang/glog"
)

func DeployShadowss(typ api.OperatorType, hosts []string, sshkey string, attr map[string]string) error {

	var hostpath, sshkeypath, playbook, attrpath string
	switch typ {
	case api.OperatorVultr:
		hostpath = ansibleVulHostFile
		sshkeypath = ansibleVulSSHKeyFile
		playbook = ansibleVulDeploySSPlayBook
		attrpath = ansibleVulAttrFile
	case api.OperatorDigitalOcean:
		hostpath = ansibleDGOCHostFile
		sshkeypath = ansibleDGOCSSHKeyFile
		attrpath = ansibleDGOCAttrFile
		playbook = ansibleDGOCDeploySSPlayBook
	}

	err := WriteSSAttrFile(attrpath, attr)
	if err != nil {
		glog.Errorf("Unexpected error write attr to file %v", err)
		return err
	}

	WriteDeplossConfigFile(hosts, privateKey, hostpath, sshkeypath)

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

	var hostpath, sshkeypath, playbook string
	switch typ {
	case api.OperatorVultr:
		hostpath = ansibleVulHostFile
		sshkeypath = ansibleVulSSHKeyFile
		playbook = ansibleVulRestartSSPlayBook
	case api.OperatorDigitalOcean:
		hostpath = ansibleDGOCHostFile
		sshkeypath = ansibleDGOCSSHKeyFile
		playbook = ansibleDGOCRestartSSPlayBook
	}

	WriteDeplossConfigFile(hosts, privateKey, hostpath, sshkeypath)

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
		playbook = ansibleVulDeployVPSPlayBook
		varfile = ansibleVulVarFile

		vpsVar = fmt.Sprintf("api_key: %s\r\n", key)
		vpsVar += fmt.Sprintf("DCID: %s\r\n", srv.Region)
		vpsVar += fmt.Sprintf("OSID: %s\r\n", srv.Image)
		vpsVar += fmt.Sprintf("SSHKEYID: %q\r\n", srv.SSHKeyID)
		vpsVar += fmt.Sprintf("VPSPLANID: %s\r\n", srv.Size)
		vpsVar += fmt.Sprintf("name: %s\r\n", srv.Name)

	case api.OperatorDigitalOcean:
		playbook = ansibleDGOCDeployVPSPlayBook
		varfile = ansibleDGOCVarFile

		vpsVar = fmt.Sprintf("do_token: %s\r\n", key)
		vpsVar += fmt.Sprintf("image_id: %s\r\n", srv.Image)
		vpsVar += fmt.Sprintf("ssh_key_ids: %s\r\n", srv.SSHKeyID)
		vpsVar += fmt.Sprintf("region: %s\r\n", srv.Region)
		vpsVar += fmt.Sprintf("size_id: %s\r\n", srv.Size)
		vpsVar += fmt.Sprintf("name: %s\r\n", srv.Name)
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
		playbook = ansibleVulDeleteVPSPlayBook
		varfile = ansibleVulVarFile

		//vultr not delete into here
	case api.OperatorDigitalOcean:
		playbook = ansibleDGOCDeleteVPSPlayBook
		varfile = ansibleDGOCVarFile

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
