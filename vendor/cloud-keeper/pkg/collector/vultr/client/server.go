package client

import (
	"encoding/json"
	"fmt"
	"gofreezer/pkg/api/unversioned"
	"sort"
	"strconv"

	"cloud-keeper/pkg/ansible"
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/collector/vultr/client/lib"
	"cloud-keeper/pkg/pagination"

	"github.com/golang/glog"
)

type servers []lib.Server

func (slice servers) Len() int {
	return len(slice)
}

func (slice servers) Less(i, j int) bool {
	return slice[i].Created < slice[j].Created
}

func (slice servers) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

//GetServers get account server information
func (c *Client) GetServers(page pagination.Pager) ([]byte, error) {
	info, err := c.vultrClient.GetServers()
	if err != nil {
		return nil, err
	}

	srvList := api.AccServerList{
		TypeMeta: api.AccountInfoType,
		ListMeta: unversioned.ListMeta{
			SelfLink: "/api/v1/apiservers",
		},
	}

	var notPage bool
	if page == nil {
		notPage = true
	} else {
		notPage = page.Empty()
	}

	srvs := servers(info)
	sort.Sort(srvs)

	var hasPage bool
	var perPage, skip uint64
	var count int
	if notPage {
		goto AllPage
	} else {
		count = len(info)
		glog.V(5).Infof("Got Total count %v \r\n", count)
		hasPage, perPage, skip = api.PagerToCondition(page, uint64(count))
		glog.V(5).Infof("Got page has %v  perpage %v skip %v\r\n", hasPage, perPage, skip)
		if hasPage {
			goto Pieces
		} else {
			goto AllPage
		}
	}

AllPage:
	for _, v := range info {
		information := make(map[string]interface{}, 1)
		information[api.OperatorVultr] = v

		srv := api.AccServer{
			TypeMeta:    api.AccServerType,
			Information: information,
		}
		srvList.Items = append(srvList.Items, srv)
	}
	goto Out

Pieces:
	for index := uint64(0); index < perPage; index++ {
		information := make(map[string]interface{}, 1)
		information[api.OperatorVultr] = info[index+skip]

		srv := api.AccServer{
			TypeMeta:    api.AccServerType,
			Information: information,
		}
		srvList.Items = append(srvList.Items, srv)
	}

	goto Out

Out:
	return json.Marshal(&srvList)
}

//CreateServer create server
func (c *Client) CreateServer(server interface{}) error {

	srv, ok := server.(*api.AccServer)
	if !ok {
		return fmt.Errorf("invalid obj type")
	}

	return ansible.DeployVPS(api.OperatorVultr, srv, c.key)
}

//DeleteServer delete server by id
func (c *Client) DeleteServer(id int64) error {

	idStr := strconv.FormatInt(id, 10)

	err := c.vultrClient.DeleteServer(idStr)

	return err
}

func (c *Client) Exec(param interface{}) error {

	execParam, ok := param.(*api.AccServerCommand)
	if !ok {
		return fmt.Errorf("invalid param")
	}
	var sshkey string
	switch execParam.Command {
	case "deployss":
		glog.V(5).Infof("Got request %+v\r\n", execParam.Deploy.HostList)
		return ansible.DeployShadowss(api.OperatorVultr, execParam.Deploy.HostList, sshkey, execParam.Deploy.Attribute)
	case "restartSS":
		return ansible.RestartShadowss(api.OperatorVultr, execParam.Deploy.HostList, sshkey)
	case "reboot":
		idStr := strconv.FormatInt(execParam.ID, 10)
		return c.vultrClient.RebootServer(idStr)
	default:
		return fmt.Errorf("not support command %s", execParam.Command)
	}

}

func (c *Client) GetSSHKey() ([]byte, error) {

	key, err := c.vultrClient.GetSSHKeys()
	if err != nil {
		return nil, err
	}

	sshKey := api.AccServerSSHKey{
		TypeMeta: api.AccServerSSHKeyType,
		Key:      key,
	}

	return json.Marshal(&sshKey)
}
