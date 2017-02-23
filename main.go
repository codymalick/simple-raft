// Simple implementation of the raft concensus algorithm using golang rpc
// Cody Malick and Jacob Broderick
package main

import (
  "fmt"
  // "net"
  // "net/rpc"
  // "strconv"
  // "log"
  // "net/http"
  "time"
)

// Reference: https://golang.org/pkg/net/rpc/
const (
	numServers = 5
)

// Leader Election: Choose a leader
// Log replication: Make sure all systems have the same view of the log
// Safety: Make sure a leader who is behind cannot be elected

func main() {
	server0 := Server{0, 0, make([]Log, numServers), 50000, 0}
	server1 := Server{1, 0, make([]Log, numServers), 50001, 0}
	server2 := Server{2, 0, make([]Log, numServers), 50002, 0}
	server3 := Server{3, 0, make([]Log, numServers), 50003, 0}
	server4 := Server{4, 0, make([]Log, numServers), 50004, 0}
	servers := []Server{server0, server1, server2, server3, server4}

	for i := 0; i < numServers; i++ {
		go servers[i].ServerRun()
    fmt.Printf("Thread %v spawned\n", servers[i].ID)
	}

  // Give the threads some time to run
  time.Sleep(time.Second * 20)

}
