// Simple implementation of the raft concensus algorithm using golang rpc
// Cody Malick and Jacob Broderick
package main

import (
	"flag"
)

const (
	numServers = 5
	timeout = 5 //seconds
)

// Leader Election: Choose a leader
// Log replication: Make sure all systems have the same view of the log
// Safety: Make sure a leader who is behind cannot be elected

func main() {
	cmdID := flag.Int("id", 0, "Usage: -id <id>")
	cmdPort := flag.String("port", ":50000", "Usage: -port <portnumber>")
	cmdState := flag.Int("state", 0, "Usage: -state <start-state>")
	flag.Parse()

	server := CreateServer(*cmdID, *cmdPort, *cmdState)

	// Exit channel keeps the program alive indefinitely
	exit := make(chan bool)

	// Tracking state of all declared servers
	server0 := CreateServer(0, ":50000", 0)
	server1 := CreateServer(1, ":50001", 0)
	server2 := CreateServer(2, ":50002", 0)
	server3 := CreateServer(3, ":50003", 0)
	server4 := CreateServer(4, ":50004", 0)

	servers := []*Server{server0,server1,server2,server3,server4}

	server.Servers = servers

	go Run(server)

	// Stay alive
	<-exit
}
