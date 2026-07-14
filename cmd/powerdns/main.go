package main

import (
	"collector/internal/receiver"
	"log"
)

func main() {
	// Cấu hình powerdns
	cfg := receiver.Config{
		Address: "0.0.0.0:6061",
		Source:  "powerdns",
	}

	err := receiver.StartTCPListener(cfg)
	if err != nil {
		log.Fatalf("Fatal error starting server: %v", err)
	}
}
