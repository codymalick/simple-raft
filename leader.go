package main

import(
    "fmt"
	// "math/rand"
	// "time"
	"log"
	// "net"
	"net/rpc"
	// "strconv"
	// "sync"
)

// SendHeartbeatRequest creates and sends a message
func SendHeartbeatRequest(source *Server, destination *Server) {
    var mes = new(Message)
    mes.SourceID = source.ID
    mes.Source = source.Port
    mes.Destination = destination.Port
    mes.Index = len(source.Log)
    mes.Epoch = source.Epoch

    // send response
    client, err := rpc.Dial("tcp", mes.Destination)
    if err != nil {
        // Fail silently
        //log.Print(err)
        fmt.Printf("No response from %v\n", destination.ID)
        return
    }
    //
    // if err != nil {
    //     log.Print(err)
    // }
    var reply = new(Message)
    err = client.Call("Server.Heartbeat", mes, reply)
    if err != nil {
        log.Print(err)
    }

	fmt.Printf("Heartbeat from %v, Epoch %v, Index %v\n", destination.ID, reply.Epoch, reply.Index)
}

// GetHeartbeat sends a heartbeat to all servers, and requests a heartbeat from
// all servers
func GetHeartbeats(s *Server) {
    for _, val := range s.Servers {
        if val.ID != s.ID {
            //fmt.Printf("%v, calling %v at %v\n",s.ID, val.ID, val.Port)
            go SendHeartbeatRequest(s, val)
        }
    }
}
