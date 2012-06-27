package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

func vital(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	var addr *net.UDPAddr
	var err error

	rand.Seed(time.Now().UnixNano())
	serverMode := flag.Bool("server", false, "run in server mode")
	join := flag.String("join", "", "join a game")
	flag.Parse()
	if *serverMode {
		vital(StartServer(1337))
	}
	if *join != "" {
		addr, err = net.ResolveUDPAddr("udp", *join)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	update := make(chan *State)
	move := make(chan MoveMsg)
	quit := make(chan bool)
	id := rand.Uint32()
	go func() {
		vital(StartGraphics(id, update, move, quit))
	}()
	if addr != nil {
		vital(StartClient(addr, id, update, move, quit))
	} else {
		vital(StartSingle(id, update, move, quit))
	}
}
