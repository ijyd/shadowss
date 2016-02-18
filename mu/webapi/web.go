package web

import (
	"github.com/orvice/shadowsocks-go/mu/user"
	"encoding/json"
	"github.com/orvice/shadowsocks-go/mu/log"
)

var (
	client = new(Client)
)

type Client struct {
	baseUrl string
	key     string
	nodeId  int
}

func (c *Client) setBaseUrl(baseUrl string) {
	c.baseUrl = baseUrl
}

func (c *Client) setKey(key string) {
	c.key = key
}

func (c *Client) setNodeId(id int) {
	c.nodeId = id
}

func (c *Client) GetUsers() ([]user.User, error) {
	var datas []*User
	res,err := c.httpGet(c.genGetUsersUrl())
	if err != nil{
		return datas,err
	}
	var resData UserDataRet
	err = json.Unmarshal([]byte(res),resData)
	if err != nil{
		return datas,err
	}
	userData := resData.Data
	log.Log.Info(len(userData))
	users := make([]user.User, len(userData))
	for k, v := range userData {
		users[k] = v
	}
	return datas, nil
}
