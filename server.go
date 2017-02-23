package main

import(
  "math/rand"
  "fmt"
  "time"
)

// Server has three possible states:
//   1. Follower
//   2. Candidate
//   3. Leader
type Server struct {
	ID    int
	Epoch int // Life span of the current leader
	Log   []Log // List of all logs
  Port  int // tcp Port
  State int // Current mode of the server
}

// Respond to a request
func (t *Server) Respond(message *Message, reply *Message) error {
  message.Value++
  *reply = *message
  return nil
}

// Heartbeat function to let the leader know they are alive
func (t *Server) Heartbeat() {

}

// GetHeartbeat sends a heartbeat to all servers, and requests a heartbeat from
// all servers
func (t *Server) GetHeartbeat() {

}

// Elect respond to an election vote, and wait for confirmation or timeout
func (t *Server) Elect() {

}

// StartElection starts an election
func (t *Server) StartElection() {

}

// GenerateRequest adds an event to the log. The leader should be the only one
// able to run this
func (t *Server) GenerateRequest() {

}

// IncrementEpoch is called for two reason:
//   1. The server starts a new election
//   2. The server receives confirmation an election is successful
func (t *Server) IncrementEpoch() {
  t.Epoch++
}



const (
	numServers = 5
)

// Timeout function
func (t *Server) randomTimeout() {
  sleep := rand.Intn(5)
  fmt.Printf("sleeping for %v seconds\n", sleep)
  time.Sleep(time.Second * time.Duration(sleep))
}
