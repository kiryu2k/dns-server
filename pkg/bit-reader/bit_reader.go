package bit_reader

import (
	"fmt"
)

type bitReader struct {
	data     []byte
	readBits uint
}

func New(data []byte) *bitReader {
	return &bitReader{data: data}
}

func (r *bitReader) ReadBits(bits uint) ([]uint8, error) {
	var (
		result = make([]uint8, bits)
		err    error
	)
	for i := range bits {
		result[i], err = r.Read()
		if err != nil {
			return nil, fmt.Errorf("read: %w", err)
		}
	}
	return result, nil
}

func (r *bitReader) Read() (uint8, error) {
	if r.readBits >= uint(8*len(r.data)) {
		return 0, ErrNoDataToRead
	}
	var (
		bytePos = int(r.readBits) / 8
		bitPos  = r.readBits % 8
	)

	bit := r.data[bytePos] >> (7 - bitPos) & 1
	r.readBits++

	return bit, nil
}
