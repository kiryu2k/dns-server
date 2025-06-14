package domain

import (
	"fmt"
	"strings"

	"github.com/kiryu2k/dns-server/pkg/utils"
)

type DnsMessage struct {
	Header   dnsHeader
	Question dnsQuestion
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
	return DnsMessage{
		Header:   decodeHeader(bytes[:headerSize]),
		Question: decodeQuestion(bytes[headerSize:]),
	}, nil
}

func (m DnsMessage) Encode() []byte {
	return append(m.encodeHeader(), m.encodeQuestion()...)
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

func (m DnsMessage) encodeQuestion() []byte {
	buf := make([]byte, 0)
	labels := strings.Split(m.Question.name, ".")
	for _, v := range labels {
		buf = append(buf, byte(len(v)))
		buf = append(buf, []byte(v)...)
	}
	buf = append(buf, '\x00')

	return utils.AppendBigEndianUint16(buf, m.Question.recordType, m.Question.recordClass)

}

func (m DnsMessage) AsReply() DnsMessage {
	m.Header.Qr = 1
	return m
}

func (m DnsMessage) WithQuestion(question dnsQuestion) DnsMessage {
	m.Header.QdCount++
	m.Question = question
	return m
}
