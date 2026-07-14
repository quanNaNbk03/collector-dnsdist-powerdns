package receiver

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"

	"collector/internal/decoder"
)

// StartTCPListener starts a TCP server on the given address.
func StartTCPListener(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to start listener on %s: %w", address, err)
	}
	defer listener.Close()

	log.Printf("Listening for DNSDist protobuf messages on TCP %s", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Accepted connection from %s", conn.RemoteAddr().String())

	for {
		// Đọc 2 bytes đầu tiên để lấy độ dài (length) của frame
		frameLength := make([]byte, 2)
		if _, err := io.ReadFull(conn, frameLength); err != nil {
			if err != io.EOF {
				log.Printf("Error reading frame length from %s: %v", conn.RemoteAddr().String(), err)
			}
			break
		}

		// dnsdist (PowerDNS) gửi 16-bit length bằng Network Byte Order (Big Endian)
		length := binary.BigEndian.Uint16(frameLength)

		// Đọc payload với độ dài chính xác bằng length đã đọc được
		payload := make([]byte, length)
		if _, err := io.ReadFull(conn, payload); err != nil {
			log.Printf("Error reading payload of length %d from %s: %v", length, conn.RemoteAddr().String(), err)
			break
		}

		// Decode the protobuf
		msg, err := decoder.DecodePBDNSMessage(payload)
		if err != nil {
			log.Printf("Failed to decode message: %v", err)
			continue
		}

		// PoC: Sử dụng Formatter để chuyển Protobuf struct thành JSON đã được format các trường human-readable
		jsonBytes, err := decoder.FormatPBDNSMessage(msg)
		if err != nil {
			log.Printf("Error formatting to JSON: %v", err)
		} else {
			log.Printf("====================\nReceived Message (%d bytes):\n%s\n", length, string(jsonBytes))
		}
	}
	log.Printf("Connection from %s closed", conn.RemoteAddr().String())
}
