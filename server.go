package main

import (
	"fmt"
	"math/rand"
	"time"
	"log"
	"net"
	"net/rpc"

)

// Server has three possible states:
//   0. Follower
//   1. Candidate
//   2. Leader
type Server struct {
	ID    int
	Epoch int   // Life span of the current leader
	Log   []Log // List of all logs
	Port  string   // tcp Port
	State int   // Current mode of the server
	Servers []*Server
	Hb 		chan bool
}

// CreateServer makes it easy to quickly create a server
func CreateServer(id int, port string, startState int) *Server {
	server := new(Server)
	server.ID = id
	server.Epoch = 0
	server.Port = port
	server.State = startState
	return server
}
// Heartbeat to a hearbeat request
func (s *Server) Heartbeat(message *Message, response *Message) error {

	response.Source = s.Port
	response.Destination = message.Source
	response.Epoch = s.Epoch
	response.Index = len(s.Log)


	// Flip bool to let client thread know we sent a heartbeat
	s.Hb <- true
	return nil
}

// Heartbeat function to let the leader know they are alive
// func Heartbeat(s *Server) {
// 	var mutex = &sync.Mutex{}
// 	// send response
// 	client, err := rpc.Dial("tcp", s.LeaderPort)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	line := "test " + strconv.Itoa(s.ID)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	var reply bool
// 	err = client.Call("Server.Respond", line, &reply)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	mutex.Lock()
// 	s.HbRequest = false
// 	mutex.Unlock()
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
func RandomTimeout(s *Server, hb chan bool) bool {
	// Milliseconds, Min = 5000, max = 10000
	sleep := rand.Int() % 10000 + 5000

	timeout := make(chan bool, false)

	select {
	case <-hb:
		return true
	case <-timeout:
		return false
	}
	// for i := sleep; i > 0; i-- {
	// 	// Every 100 milliseconds, check for an update
	// 	if i % 100 == 0 {
	// 		//fmt.Printf("Server %v, Checking for request:%v\n", s.ID, i)
	// 		if s.HbRequest {
	// 			return true
	// 		}
	// 	} else {
	// 		time.Sleep(time.Millisecond)
	// 	}


	// Timed out, start an election
	// StartElection(s)
}

// Run is the main loop of the server, which starts by activating the server,
// and looping it through timeouts.
func (s *Server) Run() {
	// Seed random for timeout
	rand.Seed(time.Now().UTC().UnixNano())

	s.Hb = make(chan bool)

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
		switch(s.State) {
			// Server is a follower
			case 0:
				// Wait for heartbeat request
				if !RandomTimeout(s) {
					// start election
					fmt.Printf("%v Started an election\n", s.ID)
					// StartElection(s)
				}
			// Server is a candidate for leader
			case 1:

			// Server is a leader
			case 2:
				fmt.Printf("%v is leader\n", s.ID)
				// Get heartbeat from all servers
				time.Sleep(time.Second)

				GetHeartbeats(s)
				// GetHeartbeat(s)
		}
	}
}
