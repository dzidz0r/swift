package server

import (
	"net"
	"strings"
)

func (s *server) Broadcast() {

}

func (s *server) GetActiveIps() []*net.IPNet {
	interfaces, err := s.getUpnRunninginterfaces()
	var addrs []*net.IPNet
	if err != nil {
		panic("error while getting up and runnign interfaces")
	}
	for _, interf := range interfaces {
		addr := s.extractIPV4Address(interf)
		if addr != nil {
			addrs = append(addrs, addr)
		}
	}
	return addrs
}

// GetBroadcastAddress calculates the broadcast address for a given IP address and subnet
// in the format ip/subnet
func (s *server) calcBroadcastAddress(ipAddress string) (string, error) {
	// Parse the IP address and subnet
	ip, ipNet, err := net.ParseCIDR(ipAddress)
	if err != nil {
		return "", err
	}

	// Get the network size in bits
	ones, bits := ipNet.Mask.Size()

	// Calculate the broadcast address
	mask := net.CIDRMask(ones, bits)
	network := ip.Mask(mask)
	broadcast := make(net.IP, len(network))
	for i := range network {
		broadcast[i] = network[i] | ^mask[i]
	}

	return broadcast.String(), nil
}

func (s *server) extractIPV4Address(iface net.Interface) *net.IPNet {
	addrs, err := iface.Addrs()
	if err != nil {
		return nil
	}
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet
		}
	}
	return nil
}

func (s *server) getUpnRunninginterfaces() ([]net.Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var upnRunning []net.Interface

	// get all up and running interfaces
	for _, interf := range interfaces {
		if flags := interf.Flags; strings.Contains(flags.String(), "running") &&
			strings.Contains(flags.String(), "up") &&
			!strings.Contains(interf.Name, "VirtualBox") &&
			!strings.Contains(interf.Name, "Loopback") {
			upnRunning = append(upnRunning, interf)
		}
	}
	return upnRunning, nil
}
