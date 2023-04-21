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
	cl1 := client.NewClient()
	cl2 := client.NewClient()
	cl3 := client.NewClient()

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second)
		cl1.Connect(":5050")
		cl2.Connect(":5050")
		cl3.Connect(":5050")

		// time.AfterFunc(time.Second*2, func() {
		// cl1.Send(file)
		// })

	}()
	//
	srvr := server.NewServer()
	wg.Add(1)
	go func() {
		defer wg.Done()
		srvr.Start()
	}()
	time.AfterFunc(time.Second*3, func() {
		srvr.Send(file)
	})

	wg.Wait()
}
