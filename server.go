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
	Hb    chan bool
	VoteRequested chan bool
	VoteReceived chan bool
	Voted int
	TotalVotes []bool
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
	server.TotalVotes = []bool{false,false,false,false,false}
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

// Elect respond to an election vote, and wait for confirmation or timeout
// A server cannot vote for a leader who has a log index less than their own
func (s *Server) Elect(message *Message, response *Message) {

	response.SourceID = s.ID
	response.Source = s.Port
	response.Destination = message.Source
	response.Epoch = s.Epoch
	response.Index = len(s.Log)

	// If you haven't already voted
	if s.Voted != -1 {
		response.Vote = false
	} else {
		// We've voted
		response.Vote = true
		s.VoteRequested <- true
	}
}

// StartElection starts an election, sends a request for a vote with the new
// epoch, and index of the last log
func (s *Server) RequestVote(source *Server, destination *Server) {
	var mes = new(Message)
	mes.Source = source.Port
	mes.Destination = destination.Port
	mes.Index = len(source.Log)
	mes.Epoch = source.Epoch

	// send response
	client, err := rpc.Dial("tcp", mes.Destination)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	var reply = new(Message)
	err = client.Call("Server.Elect", mes, reply)
	if err != nil {
		log.Fatal(err)
	}

	if reply.Vote {
		source.TotalVotes[destination.ID] = true
		source.VoteReceived <- true
	}

	fmt.Printf("%v: Vote %v, from %v\n", source.ID, destination.ID, reply.Vote)
}

func StartElection(s *Server) {

}
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
func RandomTimeout(s *Server) bool {
	// Milliseconds, Min = 5000, max = 10000
	sleep := rand.Int() % 10000 + 5000

	for i := sleep; i > 0; i-- {
		// Every 100 milliseconds, check for an update
		if i % 100 == 0 {
			// Check the channel for some input, but don't block on the channel
			select {
			case value := <-s.Hb:
				fmt.Println(s.ID, value)
				return true
			default:
				// Do nothing for now
			}
		} else {
			time.Sleep(time.Millisecond)
		}
	}
	return false
}

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
		switch(s.State) {
			// Server is a follower
			case 0:
				select {
				case value := <-s.VoteRequested:
					fmt.Printf("Vote requested %v\n", value)
				default:
					// Wait for heartbeat request
					if !RandomTimeout(s) {
						// start election
						fmt.Printf("%v Started an election\n", s.ID)
						// Switch State
						s.State=1
					}
				}

			// Server is a candidate for leader
			case 1:
				// Start an election
				StartElection(s)

			// Server is a leader
			case 2:
				select {
				case value :=<-s.VoteRequested:
					// If a vote is request while leader, surrender leadership
					s.State = 0
					fmt.Printf("Giving up leadership, vote requested %v\n", value)
				default:
					fmt.Printf("%v is leader\n", s.ID)
					// Get heartbeat from all servers
					time.Sleep(time.Second)
					GetHeartbeats(s)
				}

		}
	}
}
