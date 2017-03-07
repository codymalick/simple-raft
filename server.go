package main

import (
	"fmt"
	"math/rand"
	"time"
	"log"
	"net"
	"net/rpc"
	"strconv"
)

// Server has three possible states:
//   1. Follower
//   2. Candidate
//   3. Leader
type Server struct {
	ID    int
	Epoch int   // Life span of the current leader
	Log   []Log // List of all logs
	Port  string   // tcp Port
	State int   // Current mode of the server
}


// Respond to a request
func (s *Server) Respond(line *string, ack *bool) error {
	fmt.Printf("Server %v: %v\n", s.ID, *line)
	*ack = true
	return nil
}

// Heartbeat function to let the leader know they are alive
// func (s *Server) Heartbeat() {
//
// }
//
// // GetHeartbeat sends a heartbeat to all servers, and requests a heartbeat from
// // all servers
// func (s *Server) GetHeartbeat() {
//
// }
//
// // Elect respond to an election vote, and wait for confirmation or timeout
// // A server cannot vote for a leader who has a log index less than their own
// func (s *Server) Elect() {
//
// }
//
// // StartElection starts an election, sends a request for a vote with the new
// // epoch, and index of the last log
// func (s *Server) StartElection() {
//
// }
//
// // GenerateRequest adds an event to the log. The leader should be the only one
// // able to run this
// func (s *Server) GenerateRequest() {
//
// }
//
// // ReplicateRequest sends a new request to be recorded in all server logs
// func (s *Server) ReplicateRequest() {
//
// }
//
// // CommitRequest sets a log to be commited, only increment the index of logs
// // once a log is commmited
// func (s *Server) CommitRequest() {
//
// }
//
// // IncrementEpoch is called for two reason:
// //   1. The server starts a new election
// //   2. The server receives confirmation an election is successful
// func (s *Server) IncrementEpoch() {
// 	s.Epoch++
// }
//
// // ServerKill deactivates a server
// func (s *Server) ServerKill() {
//
// }
//
// // ServerStart restarts a server
// func (s *Server) ServerStart() {
//
// }

// RandomTimeout runs every timeout period with a lower and upper bound. These
// bounds can be set in the const section
func RandomTimeout(s *Server) {
	// Milliseconds, Min = 5000, max = 10000
	sleep := rand.Int() % 10000 + 5000

	for i := sleep; i > 0; i-- {
		// Every 100 milliseconds, check for an update
		if i % 100 == 0 {
			fmt.Printf("Server %v, Checking for request:%v\n", s.ID, i)
		} else {
			time.Sleep(time.Millisecond)
		}
	}

	// Timed out, start an election
	// StartElection(s)
}

// Run is the main loop of the server, which starts by activating the server,
// and looping it through timeouts.
func (s *Server) Run() {
	// Seed random for timeout
	rand.Seed(time.Now().UTC().UnixNano())

	go func() {
		address, err := net.ResolveTCPAddr("tcp", s.Port)
		if err != nil {
			log.Fatal(err)
		}

		inbound, err := net.ListenTCP("tcp", address)
		if err != nil {
			log.Fatal(err)
		}

		rpc.Register(s)
		rpc.Accept(inbound)
	}()

	// main loop
	for {
		// Wait for heartbeat request
		RandomTimeout(s)
		client, err := rpc.Dial("tcp", ":50000")
		if err != nil {
			log.Fatal(err)
		}

		line := "test " + strconv.Itoa(s.ID)
		if err != nil {
			log.Fatal(err)
		}
		var reply bool
		err = client.Call("Server.Respond", line, &reply)
		if err != nil {
			log.Fatal(err)
		}
	}
}
