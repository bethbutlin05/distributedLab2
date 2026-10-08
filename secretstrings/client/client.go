package main

import (
	"bufio"
	//	"net/rpc"
	"flag"
	"log"
	"net/rpc"
	"os"

	//	"bufio"
	//	"os"
	//	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
	"fmt"

	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
)

func main() {
	server := flag.String("server", "127.0.0.1:8030", "IP:port string to connect to as server")
	filePath := flag.String("file", "wordlist", "Words file")
	premium := flag.Bool("premium", false, "Use FastReverse")
	flag.Parse()
	//fmt.Println("Server: ", *server)
	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal("Open file failed: ", err)
	}
	defer file.Close()

	client, err := rpc.Dial("tcp", *server)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer client.Close()

	handler := stubs.ReverseHandler
	if *premium {
		handler = stubs.PremiumReverseHandler
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		word := scanner.Text()
		request := stubs.Request{Message: word}
		var response stubs.Response

		err := client.Call(handler, request, &response)
		if err != nil {
			log.Fatal("RPC call failed: ", err)
		}

		fmt.Printf("%s: %s\n", word, response.Message)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Scan failed: ", err)
	}

	//request := stubs.Request{Message: "Hello"}
	//response := new(stubs.Response)
	//client.Call(stubs.PremiumReverseHandler, request, response)
	//fmt.Println("Response: ", response.Message)
	//TODO: connect to the RPC server and send the request(s)
}
