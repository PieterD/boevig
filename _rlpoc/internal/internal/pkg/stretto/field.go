package stretto

import (
	"io"
	"math"
)

const Max = math.MaxUint64

type Field struct {
	seeds []seed
}

func NewField(seedGeometry []uint64) *Field {
	seeds := make([]seed, len(seedGeometry))
	for i, seedSize := range seedGeometry {
		seeds[i] = newSeed(int(seedSize))
	}
	return &Field{
		seeds: seeds,
	}
}

func (f *Field) ReSeed(r io.Reader) error {
	for _, seed := range f.seeds {
		if err := seed.ReSeed(r); err != nil {
			return err
		}
	}
	return nil
}

func (f *Field) Uint64(indices ...uint64) uint64 {
	var collector uint64
	for _, seed := range f.seeds {
		for _, index := range indices {
			collector ^= seed.At(index)
		}
	}
	return collector
}

func (f *Field) Uint64n(max uint64, indices ...uint64) uint64 {
	generated := f.Float64(indices...)
	return uint64(generated * float64(max))
}

func (f *Field) Float64(indices ...uint64) float64 {
	u := f.Uint64(indices...)
	return float64(u) / float64(Max)
}
