package stretto

import (
	"io"
)

type seed struct {
	data []uint64
}

func newSeed(size int) (s seed) {
	s.data = make([]uint64, size)
	return s
}

func (s seed) Size() uint64 {
	return uint64(len(s.data))
}

func (s seed) At(index uint64) uint64 {
	return s.data[index%s.Size()]
}

func (s seed) ReSeed(r io.Reader) error {
	buf := make([]byte, 8)
	for i := range s.data {
		_, err := io.ReadFull(r, buf)
		if err != nil {
			return err
		}
		s.data[i] = uint64(buf[0])
		s.data[i] <<= 8
		s.data[i] |= uint64(buf[1])
		s.data[i] <<= 8
		s.data[i] |= uint64(buf[2])
		s.data[i] <<= 8
		s.data[i] |= uint64(buf[3])
		s.data[i] <<= 8
		s.data[i] |= uint64(buf[4])
		s.data[i] <<= 8
		s.data[i] |= uint64(buf[5])
		s.data[i] <<= 8
		s.data[i] |= uint64(buf[6])
		s.data[i] <<= 8
		s.data[i] |= uint64(buf[7])
	}
	return nil
}
