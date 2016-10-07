package apiserverproxy

import (
	"cloud-keeper/pkg/api"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/golang/glog"
)

var localIPList = make(map[string]bool)

var apiServerList []APIServerPair

type APIServerPair struct {
	Host string
	Port int64
}

func init() {
	glog.V(5).Infof("collector ip list\r\n")
	ifaces, err := net.Interfaces()
	if err != nil {
		glog.Fatalf("Unexcepted error %v\r\v", err)
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			glog.Fatalf("Unexcepted error %v\r\v", err)
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				localIPList[v.IP.String()] = true
			case *net.IPAddr:
				localIPList[v.IP.String()] = true
			}
		}
	}

	glog.V(5).Infof("Got local ip list %v\r\n", localIPList)
}

func InitAPIServer(srv []api.APIServerSpec) {

	glog.V(5).Infof("got a  api server %v \r\n", srv)

	secure := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
		DisableKeepAlives:  true,
	}

	insecure := &http.Transport{
		DisableKeepAlives:  true,
		DisableCompression: true,
	}

	timeout := 10 * time.Second
	secureClient := &http.Client{
		Transport: secure,
		Timeout:   timeout,
	}
	insecureClient := &http.Client{
		Timeout:   timeout,
		Transport: insecure,
	}

	for _, spec := range srv {
		port := spec.Server.Port
		for _, host := range spec.HostList {
			requestInsecureURL := fmt.Sprintf("http://%s:%d", host, spec.Server.Port)
			requestSecureURL := fmt.Sprintf("https://%s:%d", host, spec.Server.Port)

			resp, err := insecureClient.Get(requestInsecureURL)
			if err == nil {
				api := APIServerPair{
					Host: host,
					Port: port,
				}
				glog.V(5).Infof("got a active api server %v \r\n", api)
				apiServerList = append(apiServerList, api)
				defer resp.Body.Close()
			}

			resp, err = secureClient.Get(requestSecureURL)
			if err == nil {
				api := APIServerPair{
					Host: host,
					Port: port,
				}
				glog.V(5).Infof("got a active api server %v \r\n", api)
				apiServerList = append(apiServerList, api)
				defer resp.Body.Close()
			}

		}
	}

	if len(apiServerList) == 0 {
		glog.Fatalf("not found any api server, shutdown node\r\n")
	}

}

func FilterRequest(addr *net.TCPAddr) string {
	host := addr.String()
	_, ok := localIPList[addr.IP.String()]
	if ok {
		host = fmt.Sprintf("%s:%d", apiServerList[0].Host, apiServerList[0].Port)
	}

	return host
}
