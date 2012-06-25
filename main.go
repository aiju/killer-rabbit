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
	rand.Seed(time.Now().UnixNano())
	serverMode := flag.Bool("server", false, "run in server mode")
	join := flag.String("join", "", "join a game")
	flag.Parse()
	if *serverMode {
		vital(StartServer(1337))
	}
	if *join != "" {
		addr, err := net.ResolveUDPAddr("udp", *join)
		if err != nil {
			fmt.Println(err)
			return
		}
		update := make(chan *State)
		move := make(chan Ent)
		go func() {
			vital(StartClient(addr, update, move))
		}()
		vital(StartGraphics(update, move))
	}
}
