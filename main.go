package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func vital(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	serverMode := flag.Bool("server", false, "run in server mode")
	join := flag.String("join", "", "join a game")
	flag.Parse()
	if *serverMode {
		fmt.Println(StartServer(1337))
	}
	if *join != "" {
		addr, err := net.ResolveUDPAddr("udp", *join)
		if err != nil {
			fmt.Println(err)
			return
		}
		update := make(chan *State)
		go func() {
			for x := range update {
				fmt.Println(x)
			}
		}()
		fmt.Println(StartClient(addr, update))
	}
}
