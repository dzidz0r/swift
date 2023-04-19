package broadcaster

import (
	"fmt"
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		broadcast(
			broadcastMessage{"me", SERVERPORT},
			time.Now().Add(time.Second*BROADCASTTIMEOUT),
		)
	}()

	wg.Wait()
}
