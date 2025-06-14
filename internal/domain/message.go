package domain

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type Message struct {
	Header   Header
	Question DnsQuestion
}

func NewMessage(packetId uint16) Message {
	return Message{Header: Header{
		Id: packetId,
	}}
}

func MessageFromBytes(bytes []byte) (Message, error) {
	if len(bytes) < headerSize {
		return Message{}, fmt.Errorf("unexpected message size: %d", len(bytes))
	}
	return Message{
		Header:   decodeHeader(bytes[:headerSize]),
		Question: decodeQuestion(bytes[headerSize:]),
	}, nil
}

func (m Message) Encode() []byte {
	return append(m.encodeHeader(), m.encodeQuestion()...)
}

func (m Message) encodeHeader() []byte {
	buf := make([]byte, 0, headerSize)
	fields := uint16(m.Header.Qr)<<15 |
		uint16(m.Header.OpCode)<<11 |
		uint16(m.Header.Aa)<<10 |
		uint16(m.Header.Tc)<<9 |
		uint16(m.Header.Rd)<<8 |
		uint16(m.Header.Ra)<<7 |
		uint16(m.Header.Z)<<4 |
		uint16(m.Header.RespCode)
	return appendBigEndianUint16(buf, m.Header.Id, fields, m.Header.QdCount, m.Header.AnCount, m.Header.NsCount, m.Header.ArCount)
}

func (m Message) encodeQuestion() []byte {
	buf := make([]byte, 0)
	labels := strings.Split(m.Question.Name, ".")
	for _, v := range labels {
		buf = append(buf, byte(len(v)))
		buf = append(buf, []byte(v)...)
	}
	buf = append(buf, '\x00')

	return appendBigEndianUint16(buf, m.Question.Type, m.Question.Class)

}

func appendBigEndianUint16(buf []byte, values ...uint16) []byte {
	for _, v := range values {
		buf = binary.BigEndian.AppendUint16(buf, v)
	}
	return buf
}

func (m Message) AsReply() Message {
	m.Header.Qr = 1
	return m
}

func (m Message) WithQuestion(question DnsQuestion) Message {
	m.Header.QdCount++
	m.Question = question
	return m
}

type Header struct {
	Id       uint16
	Qr       uint8
	OpCode   uint8
	Aa       uint8
	Tc       uint8
	Rd       uint8
	Ra       uint8
	Z        uint8
	RespCode uint8
	QdCount  uint16
	AnCount  uint16
	NsCount  uint16
	ArCount  uint16
}

func decodeHeader(header []byte) Header {
	data := binary.BigEndian.Uint16(header[queryResponseIndicatorOffset:questionCountOffset])
	return Header{
		Id:       binary.BigEndian.Uint16(header[packetIdentifierOffset:]),
		Qr:       uint8((data & queryResponseIndicatorMask) >> 15),
		OpCode:   uint8((data & operationCodeMask) >> 11),
		Aa:       uint8((data & authoritativeAnswerMask) >> 10),
		Tc:       uint8((data & truncationMask) >> 9),
		Rd:       uint8((data & recursionDesiredMask) >> 8),
		Ra:       uint8((data & recursionAvailableMask) >> 7),
		Z:        uint8((data & reservedMask) >> 4),
		RespCode: uint8(data & responseCodeMask),
		QdCount:  binary.BigEndian.Uint16(header[questionCountOffset:]),
		AnCount:  binary.BigEndian.Uint16(header[answerRecordCountOffset:]),
		NsCount:  binary.BigEndian.Uint16(header[authorityRecordCountOffset:]),
		ArCount:  binary.BigEndian.Uint16(header[additionalRecordCountOffset:]),
	}
}

type DnsQuestion struct {
	Name  string
	Type  uint16
	Class uint16
}

func decodeQuestion(question []byte) DnsQuestion {
	labelLength := question[0]
	question = question[1:]

	builder := new(strings.Builder)
	for labelLength != '\x00' {
		if builder.Len() > 0 {
			builder.WriteByte('.')
		}
		builder.Write(question[:labelLength])

		question = question[labelLength:]
		labelLength = question[0]
		question = question[1:]
	}

	return DnsQuestion{
		Name:  builder.String(),
		Type:  binary.BigEndian.Uint16(question),
		Class: binary.BigEndian.Uint16(question[recordClassOffset:]),
	}
}
