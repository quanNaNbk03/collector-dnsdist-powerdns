package decoder

import (
	"encoding/json"
	"fmt"
	"net"

	"collector/pb"

	"google.golang.org/protobuf/encoding/protojson"
)

// QType mapping based on standard DNS record types
var qtypeMap = map[uint32]string{
	1:   "A",
	2:   "NS",
	5:   "CNAME",
	6:   "SOA",
	12:  "PTR",
	15:  "MX",
	16:  "TXT",
	28:  "AAAA",
	33:  "SRV",
	255: "ANY",
}

// RCode mapping
var rcodeMap = map[uint32]string{
	0: "NOERROR",
	1: "FORMERR",
	2: "SERVFAIL",
	3: "NXDOMAIN",
	4: "NOTIMP",
	5: "REFUSED",
}

// FormatPBDNSMessage formats the message for better readability, extracting IP addresses and latency.
func FormatPBDNSMessage(msg *pb.PBDNSMessage, source string) ([]byte, error) {
	// Bước 1: Parse toàn bộ cấu trúc Protobuf thành dạng JSON tiêu chuẩn
	jsonBytes, err := protojson.MarshalOptions{
		EmitUnpopulated: false,
	}.Marshal(msg)
	if err != nil {
		return nil, err
	}

	// Bước 2: Chuyển lại thành map để ta có thể sửa đổi và thêm các trường tiện ích (helper fields)
	var logMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &logMap); err != nil {
		return nil, err
	}

	// Hiển thị IP
	if msg.From != nil {
		logMap["clientIP"] = net.IP(msg.GetFrom()).String()
		delete(logMap, "from") // Xóa cái base64 vô nghĩa đi
	}

	if msg.To != nil {
		logMap["serverIP"] = net.IP(msg.GetTo()).String()
		delete(logMap, "to")
	}

	// Giải mã serverIdentity thành chuỗi rõ ràng (thay vì base64)
	if msg.ServerIdentity != nil {
		logMap["serverID"] = string(msg.GetServerIdentity())
		delete(logMap, "serverIdentity")
	}

	// Hiển thị Query Type
	if msg.Question != nil {
		qtype := msg.GetQuestion().GetQType()
		if desc, ok := qtypeMap[qtype]; ok {
			logMap["queryTypeDesc"] = desc
		} else {
			logMap["queryTypeDesc"] = fmt.Sprintf("UNKNOWN(%d)", qtype)
		}
	}

	// Hiển thị RCode và tính Latency (chỉ áp dụng cho Response)
	if msg.Response != nil {
		rcode := msg.GetResponse().GetRcode()
		if desc, ok := rcodeMap[rcode]; ok {
			logMap["rCodeDesc"] = desc
		} else {
			logMap["rCodeDesc"] = fmt.Sprintf("UNKNOWN(%d)", rcode)
		}

		// Tính toán độ trễ: (timeSec - queryTimeSec)*1,000,000 + (timeUsec - queryTimeUsec)
		queryTimeSec := msg.GetResponse().GetQueryTimeSec()
		queryTimeUsec := msg.GetResponse().GetQueryTimeUsec()

		if queryTimeSec > 0 {
			latencyUs := int64(msg.GetTimeSec()-queryTimeSec)*1000000 + int64(msg.GetTimeUsec()) - int64(queryTimeUsec)
			logMap["latencyUs"] = latencyUs
		}
	}

	// Đính kèm source (vd: dnsdist, powerdns) để phân biệt
	logMap["source"] = source

	// Bước 3: Đóng gói lại thành chuỗi JSON đẹp mắt
	return json.MarshalIndent(logMap, "", "  ")
}
