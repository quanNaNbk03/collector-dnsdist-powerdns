# Go DNSDist Collector

A lightweight TCP server written in Go that acts as a collector for `dnsdist`'s RemoteLogger feature. It receives Protobuf-encoded DNS logs over a persistent TCP connection, decodes them, formats the data into readable JSON (resolving IPs, decoding QTypes/RCodes, and calculating latency), and logs them to the console.

## Project Structure

This project follows the Standard Go Layout to ensure scalability:
- `cmd/collector/main.go`: Application entry point, starts the TCP listener.
- `internal/receiver/tcp.go`: TCP connection handling and frame extraction.
- `internal/decoder/protobuf.go`: Raw Protobuf decoding.
- `internal/decoder/formatter.go`: Data transformation into human-readable JSON (extracting Client IP, Latency, etc.).
- `pb/`: Auto-generated Go structs from `dnsmessage.proto`.

## Prerequisites

- Go 1.20+
- `protoc` compiler (optional, if you want to re-generate the protobuf structs)
- `dnsdist` running with `RemoteLogger` enabled.

## Quick Start

1. **Build the Collector:**
   ```bash
   go build -o collector-bin cmd/collector/main.go
   ```

2. **Run the Collector:**
   ```bash
   ./collector-bin
   ```
   *The server listens on `0.0.0.0:6060` by default.*

3. **Configure `dnsdist`:**
   Add the following to your `dnsdist.conf`:
   ```lua
   rl = newRemoteLogger("127.0.0.1:6060")
   
   -- Log incoming queries
   addAction(AllRule(), RemoteLogAction(rl, nil, {serverID="dnsdist-lab-query"}))
   
   -- Log outgoing responses (note the 'true' for includeCNAME)
   addResponseAction(AllRule(), RemoteLogResponseAction(rl, nil, true, {serverID="dnsdist-lab-resp"}))
   ```

4. **Test:**
   ```bash
   systemctl restart dnsdist
   dig @127.0.0.1 example.com
   ```
   Observe the formatted JSON output in the collector's terminal!
