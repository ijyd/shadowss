package multiuser

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"shadowss/pkg/api"
	"shadowss/pkg/multiuser/apiserverproxy"
	muconfig "shadowss/pkg/multiuser/config"
)

type DoReq func(client *http.Client, url string) error

func Request(req DoReq) error {
	secure := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
		DisableKeepAlives:  true,
	}
	timeout := 5 * time.Second
	secureClient := &http.Client{
		Transport: secure,
		Timeout:   timeout,
	}

	url := string("")
	if len(apiserverproxy.ApiServerList) != 0 {
		url = fmt.Sprintf("https://%s:%d", apiserverproxy.ApiServerList[0].Host, apiserverproxy.ApiServerList[0].Port)
	}

	return req(secureClient, url)
}

func GetAPIServers(urladdr string) (*api.APIServerList, error) {

	apiservers := &api.APIServerList{}

	err := Request(func(client *http.Client, url string) error {
		url = fmt.Sprintf("%s/api/v1/apiservers", urladdr)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		muconfig.AddAuthHTTPHeader(req)

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if resp.StatusCode == 200 {
			return json.Unmarshal(body, apiservers)
		} else {
			return fmt.Errorf("get api servers error:%v", string(body))
		}

	})

	return apiservers, err
}

func UpdateNode(node *api.Node, ttl uint64) error {

	body, err := json.Marshal(node)
	if err != nil {
		return err
	}

	err = Request(func(client *http.Client, url string) error {
		url = fmt.Sprintf("%s/api/v1/nodes/%s/refresh", url, node.Name)
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
		if err != nil {
			return err
		}

		muconfig.AddAuthHTTPHeader(req)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 300 {
			return fmt.Errorf("update node %v  error:%v", node, string(body))
		}

		return nil
	})

	return err
}

func UpdateNodeUser(user *api.NodeUser) error {
	name := user.Spec.User.Name
	//need to update our node user
	user.Name = name
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = Request(func(client *http.Client, url string) error {
		url = fmt.Sprintf("%s/api/v1/nodeusers/%s", url, name)
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
		if err != nil {
			return err
		}

		muconfig.AddAuthHTTPHeader(req)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 300 {
			return fmt.Errorf("update nodeuser %v  error:%v", name, string(body))
		}

		return nil
	})

	return err
}
