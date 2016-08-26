package users

import (
	"errors"
	"net"
	"shadowsocks-go/pkg/storage"

	"github.com/golang/glog"
)

//User is a mysql users map
type Node struct {
	ID          int64  `column:"id"`
	Name        string `column:"name"`
	Type        int64  `column:"type"`
	Host        string `column:"server"`
	Method      string `column:"method"`
	Status      string `column:"status"`      //traffic for per user
	StartUserID int64  `column:"startUserID"` //upload traffic for per user
	EndUserID   int64  `column:"endUserID"`   //download traffic for per user
}

func getnodes(handle storage.Interface) (*Node, error) {
	host, err := externalIP()
	glog.Infof("Got host ip is %v\r\n", host)

	fileds := []string{"id", "name", "type", "server", "method", "status", "startUserID", "endUserID"}
	query := string("server = ?")
	queryArgs := []interface{}{host}
	selection := NewSelection(fileds, query, queryArgs)

	ctx := createContextWithValue(nodeTableName)

	var nodes []Node
	err = handle.GetToList(ctx, selection, &nodes)
	if nodes == nil {
		return nil, err
	}

	return &nodes[0], err
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}
