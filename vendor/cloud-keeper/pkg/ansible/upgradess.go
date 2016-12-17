package ansible

import (
	"fmt"

	"golib/pkg/util/exec"

	"github.com/golang/glog"
)

func UpgradeShadowss(hosts []string) error {

	var hostpath, sshkeypath, playbook, keypem string
	hostpath = absdir + ansibleUpgradeSSHostFile
	playbook = absdir + ansibleUpgradeSSPlayBook
	sshkeypath = absdir + ansibleUpgradeSSHKeyFile
	keypem = absdir + ansibleUpgradePrivateKey

	WriteDeplossConfigFile(hosts, privateKey, hostpath, sshkeypath, keypem)

	execCom := exec.New()
	cmd := execCom.Command("ansible-playbook", "-i", hostpath, playbook)
	out, err := cmd.CombinedOutput()
	if err == exec.ErrExecutableNotFound {
		glog.Errorf("Expected error ErrExecutableNotFound but got %v", err)
		return fmt.Errorf("ansible-playbook error:%v output:%v", err, string(out))
	}

	if err != nil {
		return fmt.Errorf("ansible-playbook error:%v output:%v", err, string(out))
	}
	return nil
}
