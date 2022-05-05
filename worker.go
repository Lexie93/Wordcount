package main

import (
	"os"
	"log"
	"net"
	"net/rpc"
	"wordcount_service"
)

func main() {
	
	if len(os.Args)!=2 {
		log.Fatalf("Usage: %s <port>\n", os.Args[0])
	}
	
	// Register a new rpc server
	count := new(wordcount_service.Count)
	server := rpc.NewServer()
	err := server.RegisterName("Counter", count)
	if err != nil {
		log.Fatal("Format of service Counter is not correct: ", err)
	}
	
	// Listen for incoming tcp packets on specified port.
	l, e := net.Listen("tcp", ":" + os.Args[1])
	if e != nil {
		log.Fatal("Listen error:", e)
	}
	
	// Link rpc server to the socket
	server.Accept(l)
}
