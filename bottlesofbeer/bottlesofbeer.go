package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	//	"net/rpc"
	//	"fmt"
	//	"time"
	//	"net"
)

var nextAddr string

type Buddy struct {
	turns chan int
}

// RPC handler receives a turn, queues it, acknowledges receipt
func (b *Buddy) Sing(bottles int, accepted *bool) error {
	if bottles < 0 {
		return fmt.Errorf("bottles must be >= 0")
	}
	b.turns <- bottles
	*accepted = true
	return nil
}

func bottleText(n int) string {
	switch n {
	case 0:
		return "no more bottles"
	case 1:
		return "one more bottle"
	default:
		return fmt.Sprintf("%d more bottles", n)
	}
}

func singVerse(n int) {
	if n == 0 {
		fmt.Println("No more bottles of beer on the wall, " +
			"no more bottles of beer. The song is finished!")
		return
	}
	fmt.Printf("%s of beer on the wall, %s of beer. "+
		"Take one down, pass it around, "+
		"%s of beer on the wall.\n",
		bottleText(n), bottleText(n), bottleText(n-1))
}

func main() {
	thisPort := flag.String("this", "8030", "Port for this process to listen on")
	next := flag.String("next", "127.0.0.1:8040", "IP:Port string for next member of the round.")
	start := flag.Int("n", 0, "Bottles of Beer (launches song if not 0)")
	flag.Parse()
	//TODO: Up to you from here! Remember, you'll need to both listen for
	//RPC calls and make your own.
	if *start < 0 {
		log.Fatal("Bottles must be >= 0")
	}

	buddy := &Buddy{
		turns: make(chan int, 1), //each buddy has a channel
	}

	err := rpc.Register(buddy)
	if err != nil {
		log.Fatal("Register failed: ", err)
	}

	listener, err := net.Listen("tcp", ":"+*thisPort)
	if err != nil {
		log.Fatal("Listener failed: ", err)
	}
	defer listener.Close()

	//Accept incoming RPC requests while main processes turns
	go rpc.Accept(listener)

	log.Printf("Listening on port %s; next buddy is %s", *thisPort, *next)

	//Only the final buddy launched should receive a positive -n
	if *start > 0 {
		buddy.turns <- *start
	}

	//main loop takes the turn
	for bottles := range buddy.turns {
		singVerse(bottles)

		//Zero is the final turn, don't forward it.
		if bottles == 0 {
			continue
		}

		client, err := rpc.Dial("tcp", *next)
		if err != nil {
			log.Fatal("Dial failed: ", err)
		}

		var accepted bool
		//passes next count onwards
		err = client.Call("Buddy.Sing", bottles-1, &accepted)
		client.Close()

		if err != nil {
			log.Fatal("Call failed: ", err)
		}

		if !accepted {
			log.Fatal("Not accepted")
		}
	}
}
