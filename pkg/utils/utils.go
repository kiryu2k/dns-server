package utils

import (
	"encoding/binary"
)

func AppendBigEndianUint16(buf []byte, values ...uint16) []byte {
	for _, v := range values {
		buf = binary.BigEndian.AppendUint16(buf, v)
	}
	return buf
}
