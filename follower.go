package main

import (
    "fmt"
    "math/rand"
    "time"
)

// RandomTimeout runs every timeout period with a lower and upper bound. These
// bounds can be set in the const section
func RandomTimeout(s *Server) bool {
	// Milliseconds, Min = timeout, max = timeout*2
	sleep := rand.Int() % (timeout*1000) + (timeout*1000)

	for i := sleep; i > 0; i-- {
		// Every 100 milliseconds, check for an update
		if i % 100 == 0 {
			// Check the channel for some input, but don't block on the channel
			select {
			case value := <-s.Hb:
        fmt.Printf("%v: heartbeat to %v\n", s.ID, value)
				return true
			case <-s.VoteRequested:
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

// Heartbeat responds to a hearbeat request
func (s *Server) Heartbeat(message *Message, response *Message) error {

	response.Source = s.Port
	response.Destination = message.Source
	response.Epoch = s.Epoch
	response.Index = len(s.Log)
	if message.NumServers != s.NumAliveServers {
		s.NumAliveServers = message.NumServers
		fmt.Printf("Server Status updated\n")
		s.AliveServers = message.ServerStatus
	}
	// Flip bool to let client thread know we sent a heartbeat

	s.Hb <- message.SourceID
	return nil
}

// Elect respond to an election vote, and wait for confirmation or timeout
// A server cannot vote for a leader who has a log index less than their own
func (s *Server) Elect(message *Message, response *Message) error {

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
	return nil
}
