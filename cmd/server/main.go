package main

import (
	"log"
	"net"

	"com.ityurika/go-redis-clone/internal/core"
)

func main() {
	addr := ":6378"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", addr, err)
	}
	log.Printf("server listening on %s", addr)

	srv := core.NewServer(listener)
	go func() {
		if err := srv.Start(); err != nil {
			ErrPrint(
				err,
			)
		}
	}()

	conn, err := listener.Accept();
	if err != nil {
		ErrPrint(err)
	}
	defer conn.Close()

	conn.Write([]byte("Hello from server\n"))
}

func ErrPrint(err error) {
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}