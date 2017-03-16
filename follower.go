package main

import (
    "fmt"
    "math/rand"
    "time"
)

// RandomTimeout runs every timeout period with a lower and upper bound. These
// bounds can be set in the const section
func RandomTimeout(s *Server) bool {
	// Milliseconds, Min = 5000, max = 10000
	sleep := rand.Int()%10000 + 5000

	for i := sleep; i > 0; i-- {
		// Every 100 milliseconds, check for an update
		if i % 100 == 0 {
			// Check the channel for some input, but don't block on the channel
			select {
			case value := <-s.Hb:
				fmt.Println(s.ID, value)
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

// Heartbeat to a hearbeat request
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

	s.Hb <- true
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
