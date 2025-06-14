package domain

import (
	"encoding/binary"
)

type dnsHeader struct {
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

func decodeHeader(header []byte) dnsHeader {
	data := binary.BigEndian.Uint16(header[queryResponseIndicatorOffset:questionCountOffset])
	return dnsHeader{
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
