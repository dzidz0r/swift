package broadcaster

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

const (
	broadcastPort    = 5056
	broadcastTimeout = 25
	serverPort       = 5055
)

type broadcastMessage struct {
	hostname   string
	serverPort int64
}

func (msg broadcastMessage) String() string {
	return fmt.Sprintf("%s__%d", msg.hostname, msg.serverPort)
}

func StartServer() error {
	hostname, _ := os.Hostname()
	killSig := make(chan struct{})

	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		log.Println(err)
		return err
	}
	listener, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Println(err)
		return err
	}

	// begin broadcast here
	go broadcast(
		killSig,
		broadcastMessage{hostname: hostname, serverPort: int64(serverPort)},
		time.Now().Add(time.Second+broadcastTimeout),
	)

	conn, err := listener.AcceptTCP()
	if err != nil {
		log.Println("Unable to accept connection")
		return err
	}
	conn.SetDeadline(time.Now().Add(time.Second + time.Duration(broadcastTimeout+60)))

	return nil
}

// The broadcast function sends the given message to all the nodes on its network,
// if the timeout has been elapsed, the broadcast is stopped and returns
// again, if a killSignal is received, then the broadcast is also stopped
func broadcast(killSig chan struct{}, msg broadcastMessage, timeout time.Time) {
	// get all possible broadcast addresses
	ifaces, err := getUpnRunninginterfaces()
	if err != nil {
		log.Println(err)
		return
	}
	addrs := getAllBroadcasts(ifaces)

	// send to all possible broadcast addresses
	// if timeout stop broadcast
	// if kill
	select {
	case <-killSig:
		killSig <- struct{}{}
		return
	default:
		if time.Now().Before(timeout) {
			log.Println("bursting")
			broadcastBurst(addrs, msg)
		} else {
			return
		}
	}
}

func broadcastBurst(addrs []string, msg broadcastMessage) {

	for _, addr := range addrs {
		// convert message to a byte array
		messageInBytes := []byte(msg.String())

		// Resolve the IP address and port
		udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", addr, broadcastPort))
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
		}
		broacastAddrs = append(
			broacastAddrs, address[len(address)-1].String(),
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
