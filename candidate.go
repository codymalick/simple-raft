package main

import (
    "fmt"
    "net/rpc"
    "log"
)

// StartElection is called when a server times out.
func StartElection(s *Server) {
	s.Voted = s.ID
	s.TotalVotes[s.ID] = true
	for _, val := range s.Servers {
		// Let's assume he votes for himself
		if val.ID != s.ID && s.AliveServers[val.ID] {
			fmt.Printf("%v alive servers\n", s.NumAliveServers)
			go RequestVote(s, val)
		}
	}

	for i := 0; i < s.NumAliveServers; i++ {
		<-s.VoteReceived
		if x := CheckVotes(s); x > s.NumAliveServers/2 {
			s.State = 2
			return
		}
	}

}

// CheckVotes for election win
func CheckVotes(s *Server) int {
	votes := 0
	fmt.Println(s.TotalVotes)
	for _, val := range s.TotalVotes {
		if val {
			votes++
		}
	}
	return votes
}

// RequestVote sends a request for a vote from destination
func RequestVote(source *Server, destination *Server) {
	var mes = new(Message)
	mes.Source = source.Port
	mes.Destination = destination.Port
	mes.Index = len(source.Log)
	mes.Epoch = source.Epoch

	// send response
	client, err := rpc.Dial("tcp", mes.Destination)
	if err != nil {
		//fmt.Printf("Cannot connect to %v for vote\n",destination.ID)
		return
		//log.Fatal(err)
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
