package proxyserver

import (
	"shadowsocks-go/pkg/config"
	"shadowsocks-go/pkg/connection/tcp"
	"shadowsocks-go/pkg/connection/udp"

	"github.com/golang/glog"
)

//Servers hold on a list of proxyserver
type Servers struct {
	tcpSrv    []ProxyServer
	udpSrv    []ProxyServer
	tcpSrvMap map[int64]int //id to srv index map
	udpSrvMap map[int64]int //id to srv index map
	enableUDP bool
}

//NewServers create a servers
func NewServers(udp bool) *Servers {
	return &Servers{
		tcpSrv:    make([]ProxyServer, 0),
		udpSrv:    make([]ProxyServer, 0),
		tcpSrvMap: make(map[int64]int),
		udpSrvMap: make(map[int64]int),
		enableUDP: udp,
	}
}

func appendSrv(slice []ProxyServer, elements ProxyServer) ([]ProxyServer, int) {
	n := len(slice)

	total := len(slice) + 1
	if total > cap(slice) {
		// Reallocate. Grow to 1.5 times the new size, so we can still grow.
		newSize := total*3/2 + 1
		newSlice := make([]ProxyServer, total, newSize)
		copy(newSlice, slice)
		slice = newSlice
	}

	slice = slice[:total]
	slice = append(slice, elements)
	return slice, n + 1
}

func (srv *Servers) storeSrv(tcp ProxyServer, udp ProxyServer, cfg *config.ConnectionInfo) {
	idx := 0

	srv.tcpSrv, idx = appendSrv(srv.tcpSrv, tcp)
	srv.tcpSrvMap[cfg.ID] = idx

	if srv.enableUDP {
		srv.udpSrv, idx = appendSrv(srv.udpSrv, udp)
		srv.udpSrvMap[cfg.ID] = idx
	}
}

//CheckServer create new server for users
func (srv *Servers) CheckServer(client *config.ConnectionInfo) (bool, bool) {

	var equal bool
	v, exist := srv.tcpSrvMap[client.ID]
	if exist {
		tcpSrv := srv.tcpSrv[v]
		equal = tcpSrv.Compare(client)
	}

	return exist, equal
}

//StopServer create new server for users
func (srv *Servers) StopServer(client *config.ConnectionInfo) {
	v, ok := srv.tcpSrvMap[client.ID]
	if !ok {
		glog.Warningf("not found tcp server %s\r\n", client.Port)
	} else {
		tcpSrv := srv.tcpSrv[v]
		tcpSrv.Stop()
		//remove it from slice
		srv.tcpSrv = append(srv.tcpSrv[:v], srv.tcpSrv[v+1:]...)
		//remove from map
		delete(srv.tcpSrvMap, client.ID)
	}

	if srv.enableUDP {
		v, ok := srv.udpSrvMap[client.ID]
		if !ok {
			glog.Warningf("not found tcp server %s\r\n", client.Port)
		} else {
			udpSrv := srv.udpSrv[v]
			udpSrv.Stop()
			//remove it from slice
			srv.udpSrv = append(srv.udpSrv[:v], srv.udpSrv[v+1:]...)
			//remove from map
			delete(srv.udpSrvMap, client.ID)
		}
	}
}

//StartWithConfig create new server for users
func (srv *Servers) StartWithConfig(v *config.ConnectionInfo) {

	glog.V(5).Infof("Start with %v at %p\r\n", v, v)

	tcpSrv := tcp.NewTCPServer(v)
	go tcpSrv.Run()

	var udpSrv ProxyServer
	if srv.enableUDP {
		udpSrv = udp.NewUDPServer(v)
		go udpSrv.Run()
	}

	srv.storeSrv(tcpSrv, udpSrv, v)
}

//Start create new server for user
func (srv *Servers) Start() {
	for idx := range config.ServerCfg.Clients {
		config := &config.ServerCfg.Clients[idx]
		srv.StartWithConfig(config)
	}
}
