package listener

import (
	"fmt"
	"log"
	"net"
	"time"
)

type ListenerChannel struct {
	string
	*net.UDPAddr
	error
}

func StartListener(broadcastPort int64, channel chan ListenerChannel) {
	defer close(channel)

	// Resolve the broadcast address and port
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%v", broadcastPort))
	if err != nil {
		// return "", nil, err
		log.Fatal("Error while resolving udp address: \n", err)
	}

	// Create a UDP socket to listen on
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		// return "", nil, err
		log.Fatal("Error while listening for broadcast: \n", err)
	}
	defer conn.Close()

	// Set a timeout for the socket
	conn.SetReadDeadline(time.Now().Add(time.Second * 3))
	println("now listening on port ", broadcastPort)

	var buffers string
	var remoteAddr *net.UDPAddr
	var n int

	for {
		buffer := make([]byte, 1024)
		println("looping")
		go func() {
			n, remoteAddr, err = conn.ReadFromUDP(buffer)
			if err != nil {
				log.Fatal("Error while reading broadcast: \n", err)
			}
		}()

		buffers = fmt.Sprintf("%s+%s", buffers, buffer[:n])
		if conn == nil {
			break
		}

		channel <- ListenerChannel{buffers, remoteAddr, nil}
	}

}
