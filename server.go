package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"time"
)

// Server has three possible states:
//   0. Follower
//   1. Candidate
//   2. Leader
type Server struct {
	ID              int
	Epoch           int    // Life span of the current leader
	Log             []Log  // List of all logs
	Port            string // tcp Port
	State           int    // Current mode of the server
	Servers         []*Server
	AliveServers    []bool
	Hb              chan bool
	VoteRequested   chan bool
	VoteReceived    chan bool
	Voted           int
	TotalVotes      []bool
	NumAliveServers int
}

// CreateServer makes it easy to quickly create a server
func CreateServer(id int, port string, startState int) *Server {
	server := new(Server)
	server.ID = id
	server.Epoch = 0
	server.Port = port
	server.State = startState
	server.Hb = make(chan bool)
	server.Voted = -1
	server.VoteRequested = make(chan bool)
	server.VoteReceived = make(chan bool)
	server.TotalVotes = []bool{false, false, false, false, false}
	server.AliveServers = []bool{true, true, true, true, true}
	server.NumAliveServers = numServers
	return server
}

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


// Run is the main loop of the server, which starts by activating the server,
// and looping it through timeouts.
func Run(s *Server) {
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
		switch s.State {
		// Server is a follower
		case 0:
			select {
			case value := <-s.VoteRequested:
				fmt.Printf("Vote requested %v\n", value)
			default:
				// Wait for heartbeat request
				if !RandomTimeout(s) {
					// Switch State
					s.State = 1
				}
			}

		// Server is a candidate for leader
		case 1:
			select {
			case <-s.VoteRequested:
				s.State = 0
			default:
				// Start an election
				fmt.Printf("%v Started an election\n", s.ID)
				StartElection(s)
			}

		// Server is a leader
		case 2:
			select {
			case <-s.VoteRequested:
				s.Voted = -1
				s.TotalVotes = []bool{false, false, false, false, false}
				// If a vote is request while leader, surrender leadership
				s.State = 0
			default:
				fmt.Printf("%v is leader\n", s.ID)
				// Get heartbeat from all servers
				time.Sleep(time.Second)
				GetHeartbeats(s)
			}

		}
	}
}
