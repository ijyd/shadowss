package web

import (
	"encoding/json"
	"github.com/orvice/shadowsocks-go/mu/log"
	"github.com/orvice/shadowsocks-go/mu/user"
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
	var tempUser []user.User
	res, err := c.httpGet(c.genGetUsersUrl())
	if err != nil {
		return tempUser, err
	}
	var resData UserDataRet
	err = json.Unmarshal([]byte(res), resData)
	if err != nil {
		return tempUser, err
	}
	userData := resData.Data
	log.Log.Info(len(userData))
	users := make([]user.User, len(userData))
	for k, v := range userData {
		users[k] = v
	}
	return users, nil
}
