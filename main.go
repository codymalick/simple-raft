// Simple implementation of the raft concensus algorithm using golang rpc
// Cody Malick and Jacob Broderick
package main

import (
)

// Reference: https://golang.org/pkg/net/rpc/
const (
	numServers = 2
)

// Leader Election: Choose a leader
// Log replication: Make sure all systems have the same view of the log
// Safety: Make sure a leader who is behind cannot be elected

func main() {
	exit := make(chan bool)

	c0 := make(chan bool)
	c1 := make(chan bool)
	c2 := make(chan bool)
	c3 := make(chan bool)
	c4 := make(chan bool)

	// Server 0 is our test leader
	server0 := CreateServer(0, ":50000", 2, c0)
	server1 := CreateServer(1, ":50001", 0, c1)
	server2 := CreateServer(2, ":50002", 0, c2)
	server3 := CreateServer(3, ":50003", 0, c3)
	server4 := CreateServer(4, ":50004", 0, c4)

	servers := []*Server{server0,server1,server2,server3,server4}

	for _,val := range servers {
		val.Servers = servers

		// Spawn goroutine
		go Run(val)
	}

	// Wait for anything on the exit channel to quit
	<- exit
}
