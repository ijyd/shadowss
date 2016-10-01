package apiserverproxy

import (
	"cloud-keeper/pkg/api"
	"fmt"
	"net"
	"net/http"

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
	//apiServerList = srv

	for k, spec := range srv {
		port := spec.Server.Port
		for _, host := range spec.HostList {
			requestURL := "http://" + host + ":" + spec.Server.Port
			resp, err := http.Get("")
			if err == nil {
				api := APIServerPair{
					Host: host,
					Port: port,
				}
				glog.V(5).Infof("got a active api server %v \r\n", api)
				apiServerList = append(apiServerList, api)
			}
			defer resp.Body.Close()
		}
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
