// Simple implementation of the raft concensus algorithm using golang rpc
// Cody Malick and Jacob Broderick
package main

import (
  "fmt"
  "net"
  "net/rpc"
  "strconv"
  "log"
  "net/http"
)

type Log struct {
	Item  int
	Epoch int
}

type Message struct {
	Source      int
	Destination int
  Value       int
}

type Server struct {
	Id    int
	Epoch int
	Log   []Log
  Port  int
}

func (t *Server) Respond(message *Message, reply *Message) error {
  message.Value += 1
  *reply = *message
  return nil
}

const (
	numServers = 5
)
// Reference: https://golang.org/pkg/net/rpc/
func spawnServer(server Server) {
  fmt.Println(server.Id)
  fmt.Println(server.Port)
  rpc.Register(server.Id)
  rpc.HandleHTTP()
  listen, err := net.Listen("tcp", ":" + strconv.Itoa(server.Port))
  if err != nil {
	   log.Fatal("listen error:", err)
   }
   fmt.Println("Test")
   go http.Serve(listen, nil)
}

func main() {
	server0 := Server{0, 0, make([]Log, numServers), 50000}
	server1 := Server{1, 0, make([]Log, numServers), 50001}
	server2 := Server{2, 0, make([]Log, numServers), 50002}
	server3 := Server{3, 0, make([]Log, numServers), 50003}
	server4 := Server{4, 0, make([]Log, numServers), 50004}
	servers := []Server{server0, server1, server2, server3, server4}

	for i := 0; i < numServers; i++ {
		go spawnServer(servers[i])
	}
  //
  // client, err := rpc.DialHTTP("tcp", "localhost:50004")
  // if err != nil {
  //   log.Fatal("dialing:", err)
  // }
  //
  // testMessage := Message{1,1,2}
  // var reply Message
  // //server6 := Server{5, 0, make([]Log, 6), 50005}
  // mesCall := client.Go("Server.Respond", testMessage, reply, nil)
  // replyCall := <-mesCall.Done
  //
  // fmt.Println("%v\n", replyCall.Reply)

}
