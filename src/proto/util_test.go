package proto_test

import (
	"net"
	"testing"

	"github.com/nikitawootten/dns-firewall-controller/src/proto"
	"github.com/stretchr/testify/assert"
)

func TestToProtoIpAddress(t *testing.T) {
	proto_ip, err := proto.ToProtoIpAddress(net.ParseIP("192.168.1.1"))
	assert.NoError(t, err, "no error on ipv4")
	assert.NotNil(t, proto_ip, "not nil on ipv4")

	proto_ip, err = proto.ToProtoIpAddress(net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334"))
	assert.NoError(t, err, "no error on ipv6")
	assert.NotNil(t, proto_ip, "not nil on ipv6")

	proto_ip, err = proto.ToProtoIpAddress(net.ParseIP("invalid"))
	assert.Error(t, err, "error on invalid")
	assert.Nil(t, proto_ip, "nil on invalid")
}

func TestFromProtoIpAddress(t *testing.T) {
	ip := net.ParseIP("192.168.1.1").To4()
	proto_ip, _ := proto.ToProtoIpAddress(ip)
	roundtrip_ip, err := proto.FromProtoIpAddress(proto_ip)
	assert.NoError(t, err, "no error on ipv4")
	assert.Equal(t, ip, roundtrip_ip, "round trip ipv4")

	ip = net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	proto_ip, _ = proto.ToProtoIpAddress(ip)
	roundtrip_ip, err = proto.FromProtoIpAddress(proto_ip)
	assert.NoError(t, err, "no error on ipv6")
	assert.Equal(t, ip, roundtrip_ip, "round trip ipv6")
}
