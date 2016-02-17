package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func (c *Client) genGetUsersUrl() string {
	return fmt.Sprintf("%s/users?key=%s", c.baseUrl, c.key)
}

func (c *Client) genUserTrafficUrl(id int) string {
	return fmt.Sprintf("%s/users/%d/traffic?key=%s", c.baseUrl, id, c.key)
}

func (c *Client) httpGet(urlStr string) (string, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *Client) httpPostUserTraffic(userId int, u, d string) (string, error) {
	nodeId := strconv.Itoa(c.nodeId)
	urlStr := c.genUserTrafficUrl(userId)
	resp, err := http.PostForm(urlStr,
		url.Values{"u": {u}, "d": {d}, "node_id": {nodeId}})

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
