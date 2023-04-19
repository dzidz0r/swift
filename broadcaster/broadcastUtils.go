package broadcaster

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

// calcBroadcastAddress calculates the broadcast address for a given IP address and subnet
func calcBroadcastAddress(ipAddress string) (string, error) {
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

// The broadcast function sends the given message to all the nodes on its network,
// if the timeout has been elapsed, the killsignal is used to tell the calling function to close the channel, the broadcast is then stopped and returns
// again, if a killSignal is received, then the broadcast is also stopped
func broadcast(msg broadcastMessage, timeout time.Time) {
	// get all possible broadcast addresses
	ifaces, err := getUpnRunninginterfaces()
	if err != nil {
		log.Println(err)
		return
	}
	addrs := getAllBroadcasts(ifaces)

	/*
	 send to all possible broadcast addresses
	 if timeout stop broadcast
	 if kill
	*/
	for time.Now().Before(timeout) {
		log.Println("bursting")
		broadcastBurst(addrs, msg)
		time.Sleep(time.Second)
	}
}

func broadcastBurst(addrs []string, msg broadcastMessage) {
	for _, addr := range addrs {
		// convert message to a byte array
		messageInBytes := []byte(msg.String())

		// Resolve the IP address and port
		udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", addr, BROADCASTPORT))
		if err != nil {
			log.Println(err)
			return
		}

		// Create the UDP socket
		conn, err := net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			log.Println(err)
			return
		}

		// Send the message
		_, err = conn.Write(messageInBytes)
		if err != nil {
			log.Println(err)
			return
		}

		conn.Close()
	}
}

func getAllBroadcasts(ifaces []net.Interface) []string {
	broacastAddrs := make([]string, 0)

	for _, iface := range ifaces {
		address, err := iface.Addrs()
		if err != nil {
			log.Println(err)
			continue
		}

		addr, err := calcBroadcastAddress(address[len(address)-1].String())
		if err != nil {
			log.Println(err)
			continue
		}
		broacastAddrs = append(
			broacastAddrs, addr,
		)
	}

	return broacastAddrs
}

func getUpnRunninginterfaces() ([]net.Interface, error) {
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
