package main

import (
	"collector/internal/receiver"
	"log"
)

func main() {
	// Cấu hình dnsdist
	cfg := receiver.Config{
		Address: "0.0.0.0:6060",
		Source:  "dnsdist",
	}

	err := receiver.StartTCPListener(cfg)
	if err != nil {
		log.Fatalf("Fatal error starting server: %v", err)
	}
}
