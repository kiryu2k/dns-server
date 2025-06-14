package domain

import (
	"encoding/binary"
	"strings"

	"github.com/kiryu2k/dns-server/pkg/utils"
)

type dnsQuestion struct {
	name        string
	recordType  uint16
	recordClass uint16
}

func (q dnsQuestion) encode() []byte {
	buf := make([]byte, 0)
	labels := strings.Split(q.name, ".")
	for _, v := range labels {
		buf = append(buf, byte(len(v)))
		buf = append(buf, []byte(v)...)
	}
	buf = append(buf, '\x00')

	return utils.AppendBigEndianUint16(buf, q.recordType, q.recordClass)
}

func decodeQuestions(packet []byte, count uint16) []dnsQuestion {
	result := make([]dnsQuestion, count)
	for i := range count {
		question, bytesRead := decodeQuestion(packet)
		result[i] = question
		packet = packet[bytesRead:]
	}
	return result
}

func decodeQuestion(question []byte) (dnsQuestion, int) {
	labelLength := question[0]
	question = question[1:]

	var (
		builder = new(strings.Builder)
		counter = int(1 + labelLength)
	)
	for labelLength != '\x00' {
		if builder.Len() > 0 {
			builder.WriteByte('.')
		}
		builder.Write(question[:labelLength])

		question = question[labelLength:]
		labelLength = question[0]
		question = question[1:]

		counter += int(1 + labelLength)
	}

	return dnsQuestion{
		name:        builder.String(),
		recordType:  binary.BigEndian.Uint16(question),
		recordClass: binary.BigEndian.Uint16(question[recordClassOffset:]),
	}, counter + 4
}
