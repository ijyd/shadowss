package network

import (
	"errors"
	"net"
)

func checkExternalIPV4(iface net.Interface) (external bool, ipv4 string) {
	if iface.Flags&net.FlagUp == 0 && iface.Flags&net.FlagLoopback != 0 {
		return
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return
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
		external = true
		ipv4 = ip.String()
		break
	}
	return
}

func ExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		external, ipv4 := checkExternalIPV4(iface)
		if external {
			return ipv4, nil
		} else {
			continue
		}
	}
	return "", errors.New("are you connected to the network?")
}

func ExternalMAC() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		external, _ := checkExternalIPV4(iface)
		if external {
			return iface.HardwareAddr.String(), nil
		} else {
			continue
		}
	}
	return "", errors.New("are you connected to the network?")
}
