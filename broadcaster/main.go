package broadcaster

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	BROADCASTPORT    = 5056
	BROADCASTTIMEOUT = 5
	SERVERPORT       = 5055
	SERVERTIMEOUT    = BROADCASTTIMEOUT * 2
)

type broadcastMessage struct {
	hostname   string
	serverPort int64
}

func (msg broadcastMessage) String() string {
	return fmt.Sprintf("%s__%d", msg.hostname, msg.serverPort)
}

var wg sync.WaitGroup

func StartServer() {
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", SERVERPORT))
	if err != nil {
		log.Println(err)
	}
	listener, err := net.Listen("tcp4", addr.String())
	wg.Add(2)

	if err != nil {
		log.Println(err)
	}

	conn, err := listener.Accept()
	go func() {
		defer wg.Done()
		if err != nil {
			log.Println(err)
		}
		conn.SetDeadline(time.Now().Add(time.Second * SERVERTIMEOUT))
		for {
			_, err := conn.Write([]byte("this is main"))
			if err != nil {
				return
			}
		}

	}()

	go func() {
		defer wg.Done()
		broadcast(
			broadcastMessage{"me", SERVERPORT},
			time.Now().Add(time.Second*BROADCASTTIMEOUT),
		)
	}()

	wg.Wait()
}
