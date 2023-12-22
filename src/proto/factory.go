package proto

import (
	"fmt"
	"net"
	"time"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func NewClientPolicy(client net.IP, allowed []net.IP, allow_until time.Time) (*ClientPolicy, error) {
	client_proto, err := ToProtoIpAddress(client)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize client IP: %w", err)
	}

	allowed_proto, err := ToProtoIpAddresses(allowed)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize allowed IPs: %w", err)
	}

	return &ClientPolicy{
		Client:     client_proto,
		AllowedIps: allowed_proto,
		AllowUntil: timestamppb.New(allow_until),
	}, nil
}
