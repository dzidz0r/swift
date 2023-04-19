package broadcaster

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

type broadcastMessage struct {
	hostname   string
	serverPort int64
}

func StartServer() error {
	serverPort := 5055
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
	broadcastTimeout := 15
	go startBroadcast(
		killSig,
		broadcastMessage{hostname: hostname, serverPort: int64(serverPort)},
		time.Duration(broadcastTimeout),
	)

	conn, err := listener.AcceptTCP()
	if err != nil {
		log.Println("Unable to accept connection")
		return err
	}
	conn.SetDeadline(time.Now().Add(time.Second + time.Duration(broadcastTimeout+60)))

	return nil
}

// func startServer(killSig chan<- struct{}, serverPort int64) *net.TCPConn {
// 	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", serverPort))

// 	listener, err := net.ListenTCP("tcp4", addr)
// 	if err != nil {
// 		log.Println("Unable to setup server")
// 	}

// 	conn, err := listener.AcceptTCP()
// 	if err != nil {
// 		log.Println("Unable to accept connection")
// 	}
// 	return conn
// }

func startBroadcast(killSig <-chan struct{}, msg broadcastMessage, timeout time.Duration) {
	// get all possible broadcast addresses
	// send to all possible broadcast addresses

	// if timeout stop broadcast

	// if kill

}
