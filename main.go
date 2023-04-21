package main

import (
	"fmt"
	"sync"

	"github.com/321swift/swift/server"
)

var wg sync.WaitGroup

func main() {
	fmt.Println(server.NewServer().GetActiveIps())
}
