package main

import (
	"collector/internal/receiver"
	"log"
)

func main() {
	// Cổng mặc định để lắng nghe (dnsdist RemoteLogger config)
	address := "0.0.0.0:6060"

	err := receiver.StartTCPListener(address)
	if err != nil {
		log.Fatalf("Fatal error starting server: %v", err)
	}
}
