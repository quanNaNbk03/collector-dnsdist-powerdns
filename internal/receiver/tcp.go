package receiver

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"

	"collector/internal/decoder"
)

// Config holds the configuration for the TCP listener
type Config struct {
	Address string
	Source  string // e.g., "dnsdist" or "powerdns"
}

// StartTCPListener starts a TCP server with the given configuration.
func StartTCPListener(cfg Config) error {
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return fmt.Errorf("failed to start listener on %s: %w", cfg.Address, err)
	}
	defer listener.Close()

	log.Printf("Listening for %s protobuf messages on TCP %s", cfg.Source, cfg.Address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn, cfg.Source)
	}
}

func handleConnection(conn net.Conn, source string) {
	defer conn.Close()
	log.Printf("[%s] Accepted connection from %s", source, conn.RemoteAddr().String())

	for {
		// Đọc 2 bytes đầu tiên để lấy độ dài (length) của frame
		frameLength := make([]byte, 2)
		if _, err := io.ReadFull(conn, frameLength); err != nil {
			if err != io.EOF {
				log.Printf("[%s] Error reading frame length from %s: %v", source, conn.RemoteAddr().String(), err)
			}
			break
		}

		// dnsdist (PowerDNS) gửi 16-bit length bằng Network Byte Order (Big Endian)
		length := binary.BigEndian.Uint16(frameLength)

		// Đọc payload với độ dài chính xác bằng length đã đọc được
		payload := make([]byte, length)
		if _, err := io.ReadFull(conn, payload); err != nil {
			log.Printf("[%s] Error reading payload of length %d from %s: %v", source, length, conn.RemoteAddr().String(), err)
			break
		}

		// Decode the protobuf
		msg, err := decoder.DecodePBDNSMessage(payload)
		if err != nil {
			log.Printf("[%s] Failed to decode message: %v", source, err)
			continue
		}

		// PoC: Sử dụng Formatter để chuyển Protobuf struct thành JSON đã được format các trường human-readable
		jsonBytes, err := decoder.FormatPBDNSMessage(msg, source)
		if err != nil {
			log.Printf("[%s] Error formatting to JSON: %v", source, err)
		} else {
			log.Printf("====================\n[%s] Received Message (%d bytes):\n%s\n", source, length, string(jsonBytes))
		}
	}
	log.Printf("[%s] Connection from %s closed", source, conn.RemoteAddr().String())
}
