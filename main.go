package main

import (
	"sync"

	"github.com/321swift/swift/server"
)

var wg sync.WaitGroup

func main() {
	serv := server.NewServer()
	serv.Broadcast()

}
