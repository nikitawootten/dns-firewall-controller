package dns

import (
	"fmt"
	"net"

	dnstap "github.com/dnstap/golang-dnstap"
	"github.com/miekg/dns"
)

type DNSRecordType string

const (
	ARecord    DNSRecordType = "A"
	AAAARecord DNSRecordType = "AAAA"
)

type DNSRecord struct {
	RecordType DNSRecordType
	RecordIP   net.IP
	TTL        uint32
}

type DNSResponse struct {
	SourceAddress net.IP
	Records       []DNSRecord
}

func ParseDNSTapMessage(dnstapMessage *dnstap.Message) (*DNSResponse, error) {
	sourceAddress, err := extractSourceAddress(dnstapMessage)
	if err != nil {
		return nil, err
	}

	records, err := extractRecords(dnstapMessage)
	if err != nil {
		return nil, err
	}

	record := DNSResponse{
		SourceAddress: sourceAddress,
		Records:       records,
	}

	return &record, nil
}

func ParseDNSTap(dt *dnstap.Dnstap) (*DNSResponse, error) {
	if dt == nil {
		return nil, fmt.Errorf("nil dnstap message")
	}

	dnstapMessage := dt.GetMessage()
	if dnstapMessage == nil {
		return nil, fmt.Errorf("nil dnstap message")
	}

	return ParseDNSTapMessage(dnstapMessage)
}

func dnsTapCallbackAdapter(callback func(*DNSResponse)) func(*dnstap.Dnstap) {
	return func(dt *dnstap.Dnstap) {
		response, err := ParseDNSTap(dt)
		if err != nil {
			return
		}

		callback(response)
	}
}

func extractSourceAddress(dnstapMessage *dnstap.Message) (net.IP, error) {
	if dnstapMessage == nil {
		return nil, fmt.Errorf("nil dnstap message")
	}

	rawAddress := dnstapMessage.GetQueryAddress()
	if rawAddress == nil {
		return nil, fmt.Errorf("nil query address")
	}

	return rawAddress, nil
}

func extractRecords(dnstapMessage *dnstap.Message) ([]DNSRecord, error) {
	if dnstapMessage == nil {
		return nil, fmt.Errorf("nil dnstap message")
	}

	records := []DNSRecord{}

	if dnstapMessage.ResponseMessage == nil {
		return nil, fmt.Errorf("nil response message")
	}

	dnsMsg := new(dns.Msg)
	if err := dnsMsg.Unpack(dnstapMessage.ResponseMessage); err != nil {
		return nil, fmt.Errorf("failed to unpack response message: %v", err)
	}

	for _, answer := range dnsMsg.Answer {
		record, err := extractRecord(answer)
		if err != nil {
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

func extractRecord(answer dns.RR) (DNSRecord, error) {
	switch answer := answer.(type) {
	case *dns.A:
		return DNSRecord{
			RecordType: ARecord,
			RecordIP:   answer.A,
			TTL:        answer.Hdr.Ttl,
		}, nil
	case *dns.AAAA:
		return DNSRecord{
			RecordType: AAAARecord,
			RecordIP:   answer.AAAA,
			TTL:        answer.Hdr.Ttl,
		}, nil
	default:
		return DNSRecord{}, fmt.Errorf("unsupported record type: %v", answer)
	}
}
