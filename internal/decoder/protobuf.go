package decoder

import (
	"fmt"
	"collector/pb"
	"google.golang.org/protobuf/proto"
)

// DecodePBDNSMessage takes a raw byte payload and decodes it into a PBDNSMessage.
func DecodePBDNSMessage(payload []byte) (*pb.PBDNSMessage, error) {
	msg := &pb.PBDNSMessage{}
	err := proto.Unmarshal(payload, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}
	return msg, nil
}
