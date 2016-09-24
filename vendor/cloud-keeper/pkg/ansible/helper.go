package ansible

import (
	"fmt"
	"io/ioutil"
	"os"
)

func WriteCreateVPSVarFile(data []byte, file string) error {

	err := ioutil.WriteFile(file, data, os.FileMode(0644))
	if err != nil {
		return err
	}

	return nil

}

func WriteDeleteVPSVarFile(data []byte, file string) error {

	err := ioutil.WriteFile(file, data, os.FileMode(0644))
	if err != nil {
		return err
	}

	return nil

}

func WriteDeplossConfigFile(hosts []string, sshkey, hostpath, sshkeypath string) error {

	hostList := fmt.Sprint("[cloud]\r\n")

	for _, v := range hosts {
		hostList += fmt.Sprintf("%s ", v)

		hostList += fmt.Sprintf("ansible_connection=ssh ")
		hostList += fmt.Sprintf("ansible_ssh_user=root ")
		hostList += fmt.Sprintf("ansible_ssh_private_key_file=%s \r\n", keyFile)
	}

	ioutil.WriteFile(hostpath, []byte(hostList), os.FileMode(0644))
	ioutil.WriteFile(sshkeypath, []byte(sshkey), os.FileMode(0600))

	return nil
}
