package domain

import (
	"fmt"

	"github.com/kiryu2k/dns-server/pkg/utils"
)

type DnsMessage struct {
	Header    dnsHeader
	Questions []dnsQuestion
}

func NewMessage(packetId uint16) DnsMessage {
	return DnsMessage{Header: dnsHeader{
		Id: packetId,
	}}
}

func MessageFromBytes(bytes []byte) (DnsMessage, error) {
	if len(bytes) < headerSize {
		return DnsMessage{}, fmt.Errorf("unexpected message size: %d", len(bytes))
	}

	header := decodeHeader(bytes[:headerSize])

	return DnsMessage{
		Header:    header,
		Questions: decodeQuestions(bytes[headerSize:], header.QdCount),
	}, nil
}

func (m DnsMessage) Encode() []byte {
	return append(m.encodeHeader(), m.encodeQuestions()...)
}

func (m DnsMessage) encodeHeader() []byte {
	buf := make([]byte, 0, headerSize)
	fields := uint16(m.Header.Qr)<<15 |
		uint16(m.Header.OpCode)<<11 |
		uint16(m.Header.Aa)<<10 |
		uint16(m.Header.Tc)<<9 |
		uint16(m.Header.Rd)<<8 |
		uint16(m.Header.Ra)<<7 |
		uint16(m.Header.Z)<<4 |
		uint16(m.Header.RespCode)
	return utils.AppendBigEndianUint16(buf, m.Header.Id, fields, m.Header.QdCount, m.Header.AnCount, m.Header.NsCount, m.Header.ArCount)
}

func (m DnsMessage) encodeQuestions() []byte {
	result := make([]byte, 0)
	for _, q := range m.Questions {
		result = append(result, q.encode()...)
	}
	return result
}

func (m DnsMessage) AsReply() DnsMessage {
	m.Header.Qr = 1
	return m
}

func (m DnsMessage) WithQuestions(questions ...dnsQuestion) DnsMessage {
	m.Header.QdCount += uint16(len(questions))
	m.Questions = append(m.Questions, questions...)
	return m
}
