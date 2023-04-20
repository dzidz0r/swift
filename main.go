package main

import (
	"sync"
	"time"

	"github.com/321swift/swift/client"
	"github.com/321swift/swift/server"
)

var wg sync.WaitGroup

func main() {

	file := "./go.mod"
	cli := client.NewClient()

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 2)
		cli.Connect(":5050")

		time.AfterFunc(time.Second*2, func() {
			cli.Send(file)
		})

	}()
	//
	srvr := server.NewServer()
	srvr.Start()
	srvr.Receive()

	wg.Wait()
}
