package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bdshkaaa/tcp-server-go/internal/server"
)

func main() {
	addr := ":9000"
	workerCount := 4

	srv := server.New(addr, workerCount)

	go func() {
		log.Printf("Starting TCP server on %s with %d workers", addr, workerCount)
		if err := srv.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	srv.Stop()
	log.Println("Server stopped")
}
