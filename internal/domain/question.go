package domain

import (
	"encoding/binary"
	"strings"
)

type dnsQuestion struct {
	name        string
	recordType  uint16
	recordClass uint16
}

func decodeQuestion(question []byte) dnsQuestion {
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

	return dnsQuestion{
		name:        builder.String(),
		recordType:  binary.BigEndian.Uint16(question),
		recordClass: binary.BigEndian.Uint16(question[recordClassOffset:]),
	}
}
