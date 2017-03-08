// Simple implementation of the raft concensus algorithm using golang rpc
// Cody Malick and Jacob Broderick
package main

import (
	//"fmt"
	// "net"
	// "net/rpc"
	// "strconv"
	// "log"
	// "net/http"
	"time"
)

// Reference: https://golang.org/pkg/net/rpc/
const (
	numServers = 2
)

// Leader Election: Choose a leader
// Log replication: Make sure all systems have the same view of the log
// Safety: Make sure a leader who is behind cannot be elected

func main() {
	// Server 0 is our test leader
	server0 := CreateServer(0, ":50000", 2)
	server1 := CreateServer(1, ":50001", 0)
	server2 := CreateServer(2, ":50002", 0)
	server3 := CreateServer(3, ":50003", 0)
	server4 := CreateServer(4, ":50004", 0)

	servers := []*Server{server0,server1,server2,server3,server4}

	for _,val := range servers {
		val.Servers = servers
		go val.Run()
	}

	// server0.Run()
	// Give the threads some time to run
	time.Sleep(time.Second * 60)

}
