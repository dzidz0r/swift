package listener

import (
	"fmt"
	"net"
	"time"
)

func Listener(ch chan string) error {
	defer close(ch)
	// Resolve the broadcast address and port
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%v", 5050))
	if err != nil {
		fmt.Println("unable to resolve address, ", err)
		return err
	}

	// Create a UDP socket to listen on
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Unable to listen for broadcast: ", err)
		return err
	}
	defer conn.Close()

	// Set a timeout for the socket
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	fmt.Println("Now listening on port", 5050)

	//wait for message
	buffer := make([]byte, 256)
	_, remoteAddr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Unable to read broadcast: ", err)
		return err
	}
	// fmt.Printf("%v=%v", buffer, remoteAddr)
	ch <- fmt.Sprintf("%v=%v", buffer, remoteAddr)

	return nil

}
