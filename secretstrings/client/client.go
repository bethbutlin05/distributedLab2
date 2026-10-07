package main

import (
	//	"net/rpc"
	"flag"
	"log"
	"net/rpc"

	//	"bufio"
	//	"os"
	//	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
	"fmt"

	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
)

func main() {
	server := flag.String("server", "127.0.0.1:8030", "IP:port string to connect to as server")
	flag.Parse()
	fmt.Println("Server: ", *server)
	client, err := rpc.Dial("tcp", *server)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer client.Close()

	request := stubs.Request{Message: "Hello"}
	response := new(stubs.Response)
	client.Call(stubs.PremiumReverseHandler, request, response)
	fmt.Println("Response: ", response.Message)
	//TODO: connect to the RPC server and send the request(s)
}
