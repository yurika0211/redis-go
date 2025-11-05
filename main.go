package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"com.ityurika/go-redis-clone/internal/core"
)

func main() {
	addr := ":6379"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", addr, err)
	}
	log.Printf("server listening on %s", addr)

	srv := core.NewServer(listener)
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	// wait for shutdown signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("shutting down server...")
	srv.Stop()
}
