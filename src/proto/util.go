package proto

import (
	"encoding/binary"
	"fmt"
	"net"
)

func ToProtoIpAddress(ip net.IP) (*IpAddress, error) {
	if ip4 := ip.To4(); ip4 != nil {
		data := binary.BigEndian.Uint32(ip4)
		return &IpAddress{IpAddress: &IpAddress_Ipv4{Ipv4: data}}, nil
	} else if ip6 := ip.To16(); ip6 != nil {
		return &IpAddress{IpAddress: &IpAddress_Ipv6{Ipv6: ip6}}, nil
	}

	return nil, fmt.Errorf("invalid net.IP: %v", ip)
}

func ToProtoIpAddresses(ips []net.IP) ([]*IpAddress, error) {
	var proto_ips []*IpAddress
	for _, ip := range ips {
		proto_ip, err := ToProtoIpAddress(ip)
		if err != nil {
			return nil, err
		}
		proto_ips = append(proto_ips, proto_ip)
	}

	return proto_ips, nil
}

func FromProtoIpAddress(ip *IpAddress) (net.IP, error) {
	switch ip := ip.IpAddress.(type) {
	case *IpAddress_Ipv4:
		data := make([]byte, 4)
		binary.BigEndian.PutUint32(data, ip.Ipv4)
		return net.IP(data), nil
	case *IpAddress_Ipv6:
		return ip.Ipv6, nil
	}

	return nil, fmt.Errorf("invalid protobuf IP address: %v", ip)
}

func FromProtoIpAddresses(ips []*IpAddress) ([]net.IP, error) {
	var proto_ips []net.IP
	for _, ip := range ips {
		proto_ip, err := FromProtoIpAddress(ip)
		if err != nil {
			return nil, err
		}
		proto_ips = append(proto_ips, proto_ip)
	}

	return proto_ips, nil
}
