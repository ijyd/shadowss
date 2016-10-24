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

// func pageForURL(urlText string) (int, error) {
// 	u, err := url.ParseRequestURI(urlText)
// 	if err != nil {
// 		return 0, err
// 	}
//
// 	pageStr := u.Query().Get("page")
// 	page, err := strconv.Atoi(pageStr)
// 	if err != nil {
// 		return 0, err
// 	}
//
// 	return page, nil
// }
//
// func getItemsCount(p *godo.Pages, requirePerPage, actualPageNum int) (int, error) {
//
// 	if p != nil {
// 		firstPage, err := pageForURL(p.First)
// 		if err != nil {
// 			return 0, err
// 		}
//
// 		LastPage, err := pageForURL(p.Last)
// 		if err != nil {
// 			return 0, err
// 		}
//
// 		if firstPage == LastPage {
//
// 		} else {
// 			pages := LastPage - firstPage
// 			return pages * requirePerPage, nil
// 		}
//
// 		return firstPage + 1, nil
// 	}
//
// 	return 1 * actualPageNum, nil
// }

//GetServers get account server information
func (c *Client) GetServers(page pagination.Pager) ([]byte, error) {

	//first get all
	list := []godo.Droplet{}
	opt := &godo.ListOptions{}
	droplets, _, err := c.client.Droplets.List(opt)
	if err != nil {
		return nil, err
	}
	dropletsSize := len(droplets)

	var notPage bool
	if page == nil {
		notPage = true
	} else {
		notPage = page.Empty()
	}

	pagenum, perPage := page.RequirePage()
	if !notPage {
		opt.Page = int(pagenum)
		opt.PerPage = int(perPage)
	}

	var response *godo.Response

	for {
		droplets, resp, err := c.client.Droplets.List(opt)
		response = resp
		if err != nil {
			return nil, err
		}

		// append the current page's droplets to our list
		for _, d := range droplets {
			list = append(list, d)
		}

		// if we are at the last page, break out the for loop
		if response.Links == nil || response.Links.IsLastPage() {
			break
		}

		page, err := response.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	//need to update our list total
	if !notPage {
		page.SetItemTotal(uint64(dropletsSize))
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
		return ansible.DeployShadowss(api.OperatorDigitalOcean, execParam.Deploy.HostList, sshkey, execParam.Deploy.Attribute)
	case "restartSS":
		return ansible.RestartShadowss(api.OperatorDigitalOcean, execParam.Deploy.HostList, sshkey)
	case "reboot":
		_, _, err := c.client.DropletActions.Reboot(int(execParam.ID))
		return err
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
