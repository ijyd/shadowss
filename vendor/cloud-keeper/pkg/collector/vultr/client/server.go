package client

import (
	"fmt"
	"gofreezer/pkg/pagination"
	"sort"
	"strconv"

	"cloud-keeper/pkg/ansible"
	api "cloud-keeper/pkg/api"
	"cloud-keeper/pkg/collector/common"
	"cloud-keeper/pkg/collector/vultr/client/lib"

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

func (c *Client) newAccserverSpec(srv *lib.Server) *api.AccServer {
	accsrv := &api.AccServer{}

	accsrv.Spec.Vultr = api.VultrServerInfo{
		Location:    srv.Location,
		Status:      srv.Status,
		Name:        srv.Name,
		CreatedTime: srv.Created,

		IPV4Addr:    srv.MainIP,
		IPV4NetMask: srv.NetmaskV4,
		IPV4Gateway: srv.GatewayV4,

		PendingCharges:   srv.PendingCharges,
		CostPerMonth:     srv.Cost,
		AllowedBandWidth: srv.AllowedBandwidth,
		CurrentBandwidth: srv.CurrentBandwidth,
	}
	accsrv.Spec.ID = srv.ID
	accsrv.Spec.Size = srv.RAM

	return accsrv
}

func (c *Client) GetServer(id int) (*api.AccServer, error) {
	idStr := strconv.FormatInt(int64(id), 10)
	srv, err := c.vultrClient.GetServer(idStr)
	if err != nil {
		return nil, err
	}

	accsrv := c.newAccserverSpec(&srv)

	return accsrv, nil
}

//GetServers get account server information
func (c *Client) GetServers(page pagination.Pager) ([]api.AccServer, error) {
	info, err := c.vultrClient.GetServers()
	if err != nil {
		return nil, err
	}

	var srvList []api.AccServer

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
		hasPage, perPage, skip = false, 0, 0 //api.PagerToCondition(page, uint64(count))
		glog.V(5).Infof("Got page has %v  perpage %v skip %v\r\n", hasPage, perPage, skip)
		if hasPage {
			goto Pieces
		} else {
			goto AllPage
		}
	}

AllPage:
	for _, v := range info {
		accsrv := c.newAccserverSpec(&v)
		srvList = append(srvList, *accsrv)
	}
	goto Out

Pieces:
	for index := uint64(0); index < perPage; index++ {
		accsrv := c.newAccserverSpec(&info[index+skip])
		srvList = append(srvList, *accsrv)
	}

	goto Out

Out:
	return srvList, nil
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

func (c *Client) Exec(param *api.AccExec) error {

	execParam := param

	var sshkey string
	switch execParam.Spec.Command {
	case "deployss":
		glog.V(5).Infof("Got request %+v\r\n", execParam.Spec.Deploy.HostList)
		return ansible.DeployShadowss(api.OperatorVultr, execParam.Spec.Deploy.HostList, sshkey, execParam.Spec.Deploy.Attribute)
	case "restartSS":
		return ansible.RestartShadowss(api.OperatorVultr, execParam.Spec.Deploy.HostList, sshkey)
	case "reboot":
		idStr := strconv.FormatInt(execParam.Spec.ID, 10)
		return c.vultrClient.RebootServer(idStr)
	default:
		return fmt.Errorf("not support command %s", execParam.Spec.Command)
	}

}

func (c *Client) GetSSHKey() (*api.AccSSHKey, error) {

	accKey := &api.AccSSHKey{}
	key, err := c.vultrClient.GetSSHKeys()
	if err != nil {
		return nil, err
	}

	for _, v := range key {

		sshKey := api.SSHKey{
			KeyID: v.ID,
			Key:   v.Key,
			Name:  v.Name,
		}

		accKey.Spec.Keys = append(accKey.Spec.Keys, sshKey)
	}

	return accKey, nil
}

func (c *Client) ServerExec(serverid int, cmd string) error {
	switch cmd {
	case common.ServerExecRestart:
		idStr := strconv.FormatInt(int64(serverid), 10)
		return c.vultrClient.RebootServer(idStr)
	}

	return fmt.Errorf("not support command %v", cmd)
}
