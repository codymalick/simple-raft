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
func (s *Server) Respond(message *Message, reply *Message) error {
  //message.Value++
  *reply = *message
  return nil
}

// Heartbeat function to let the leader know they are alive
func (s *Server) Heartbeat() {

}

// GetHeartbeat sends a heartbeat to all servers, and requests a heartbeat from
// all servers
func (s *Server) GetHeartbeat() {

}

// Elect respond to an election vote, and wait for confirmation or timeout
// A server cannot vote for a leader who has a log index less than their own
func (s *Server) Elect() {

}

// StartElection starts an election, sends a request for a vote with the new
// epoch, and index of the last log
func (s *Server) StartElection() {

}

// GenerateRequest adds an event to the log. The leader should be the only one
// able to run this
func (s *Server) GenerateRequest() {

}

// ReplicateRequest sends a new request to be recorded in all server logs
func (s *Server) ReplicateRequest() {

}

// CommitRequest sets a log to be commited, only increment the index of logs
// once a log is commmited
func (s *Server) CommitRequest() {

}

// IncrementEpoch is called for two reason:
//   1. The server starts a new election
//   2. The server receives confirmation an election is successful
func (s *Server) IncrementEpoch() {
  s.Epoch++
}

// ServerKill deactivates a server
func (s *Server) ServerKill() {

}

// ServerStart restarts a server
func (s *Server) ServerStart() {

}

// RandomTimeout runs every timeout period with a lower and upper bound. These
// bounds can be set in the const section
func (s *Server) RandomTimeout() {
  // Min = 5, max = 20
  sleep := rand.Int() % 20 + 5
  fmt.Printf("sleeping for %v seconds\n", sleep)
  time.Sleep(time.Second * time.Duration(sleep))
}

// SpawnServer runs the initial setup of the server
func (s *Server) SpawnServer() {
  fmt.Printf("Server %v is running\n", s.ID)
  for i := 1; i > 0; i++ {
    s.RandomTimeout()
  }
}

// ServerRun is the main loop of the server, which starts by activating the server,
// and looping it through timeouts.
func (s *Server) ServerRun() {
  s.SpawnServer()
}
