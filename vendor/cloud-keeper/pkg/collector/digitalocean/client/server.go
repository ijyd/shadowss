package client

import (
	"encoding/json"
	"fmt"
	"gofreezer/pkg/api/unversioned"

	"cloud-keeper/pkg/ansible"
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/pagination"

	"github.com/digitalocean/godo"
)

//GetServers get account server information
func (c *Client) GetServers(page pagination.Pager) ([]byte, error) {
	list := []godo.Droplet{}
	opt := &godo.ListOptions{}

	var notPage bool
	if page == nil {
		notPage = true
	} else {
		notPage = page.Empty()
	}

	if !notPage {
		page, perPage := page.RequirePage()
		opt.Page = int(page)
		opt.PerPage = int(perPage)
	}

	for {
		droplets, resp, err := c.client.Droplets.List(opt)
		if err != nil {
			return nil, err
		}

		// append the current page's droplets to our list
		for _, d := range droplets {
			list = append(list, d)
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	srvList := api.AccServerList{
		TypeMeta: api.AccountInfoType,
		ListMeta: unversioned.ListMeta{
			SelfLink: "/api/v1/apiservers",
		},
	}

	for _, v := range list {
		information := make(map[string]interface{}, 1)
		information[api.OperatorDigitalOcean] = v

		srv := api.AccServer{
			TypeMeta:    api.AccServerType,
			Information: information,
		}
		srvList.Items = append(srvList.Items, srv)
	}

	return json.Marshal(&srvList)
}

//CreateServer create server
func (c *Client) CreateServer(server interface{}) error {
	srv, ok := server.(*api.AccServer)
	if !ok {
		return fmt.Errorf("invalid obj type")
	}

	return ansible.DeployVPS(api.OperatorDigitalOcean, srv, c.key)
}

//DeleteServer delete server by id
func (c *Client) DeleteServer(id int64) error {

	return ansible.DeleteVPS(api.OperatorDigitalOcean, id, c.key)
}

func (c *Client) Exec(param interface{}) error {
	execParam, ok := param.(*api.AccServerCommand)
	if !ok {
		return fmt.Errorf("invalid param")
	}

	var sshkey string
	switch execParam.Command {
	case "deployss":
		return ansible.DeployShadowss(api.OperatorDigitalOcean, execParam.Deploy.HostList, sshkey)
	case "restartSS":
		return ansible.RestartShadowss(api.OperatorDigitalOcean, execParam.Deploy.HostList, sshkey)
	default:
		return fmt.Errorf("not support command %s", execParam.Command)
	}

}

func (c *Client) GetSSHKey() ([]byte, error) {
	list := []godo.Key{}
	opt := &godo.ListOptions{}

	for {
		keys, resp, err := c.client.Keys.List(opt)
		if err != nil {
			return nil, err
		}

		// append the current page's droplets to our list
		for _, d := range keys {
			list = append(list, d)
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	sshKey := api.AccServerSSHKey{
		TypeMeta: api.AccServerSSHKeyType,
		Key:      list,
	}

	return json.Marshal(&sshKey)
}
