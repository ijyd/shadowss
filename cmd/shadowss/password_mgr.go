package shadowss

import (
	"net"
	"sync"
)

// PasswdManager is hold user port and password
type PasswdManager struct {
	sync.Mutex
	portListener map[int]*PortListener
	udpListener  map[int]*UDPListener
}

var passwdManager = PasswdManager{portListener: map[int]*PortListener{}}

func (pm *PasswdManager) add(password string, port int, listener net.Listener) {
	pm.Lock()
	pm.portListener[port] = &PortListener{password, listener}
	pm.Unlock()
}

func (pm *PasswdManager) addUDP(password string, port int, listener *net.UDPConn) {
	pm.Lock()
	pm.udpListener[port] = &UDPListener{password, listener}
	pm.Unlock()
}

func (pm *PasswdManager) get(port int) (pl *PortListener, ok bool) {
	pm.Lock()
	pl, ok = pm.portListener[port]
	pm.Unlock()
	return
}

func (pm *PasswdManager) getUDP(port int) (pl *UDPListener, ok bool) {
	pm.Lock()
	pl, ok = pm.udpListener[port]
	pm.Unlock()
	return
}

func (pm *PasswdManager) del(port int) {
	pl, ok := pm.get(port)
	if !ok {
		return
	}
	pl.listener.Close()
	pm.Lock()
	delete(pm.portListener, port)
	pm.Unlock()
}

// Update port password would first close a port and restart listening on that
// port. A different approach would be directly change the password used by
// that port, but that requires **sharing** password between the port listener
// and password manager.
// func (pm *PasswdManager) updatePortPasswd(password, method string, port int, auth bool) {
// 	pl, ok := pm.get(port)
// 	if !ok {
// 		log.Printf("new port %s added\n", port)
// 	} else {
// 		if pl.password == password {
// 			return
// 		}
// 		log.Printf("closing port %s to update password\n", port)
// 		pl.listener.Close()
// 	}
// 	// run will add the new port listener to passwdManager.
// 	// So there maybe concurrent access to passwdManager and we need lock to protect it.
// 	go run(password, method, port, auth)
// 	if udp {
// 		pl, _ := pm.getUDP(port)
// 		pl.listener.Close()
// 		go runUDP(password, method, port)
// 	}
// }
