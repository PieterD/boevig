package rang

import (
	"iter"
)

type Ordered[T any] struct {
	less func(a, b T) bool
}

func NewOrdered[T any](less func(a, b T) bool) *Ordered[T] {
	return &Ordered[T]{
		less: less,
	}
}

func (o Ordered[T]) newSeekable() Seekable[T] {
	return Seekable[T]{
		seek:      new(bool),
		seekValue: new(T),
		less:      o.less,
	}
}

func (o Ordered[T]) SeekIterator(cons SeekConstructor[T]) iter.Seq[Seekable[T]] {
	return func(yield func(Seekable[T]) bool) {
		seekable := o.newSeekable()
		seq := cons(nil)
		for {
			for v := range seq {
				seekable.value = v
				if !yield(seekable) {
					return
				}
				if *seekable.seek {
					break
				}
			}
			if !*seekable.seek {
				return
			}
			*seekable.seek = false
			seq = cons(seekable.seekValue)
		}
	}
}

func (o Ordered[T]) Intersect(seqs ...iter.Seq[Seekable[T]]) iter.Seq[Seekable[T]] {
	return func(yield func(Seekable[T]) bool) {
		holders := newSeqHolders(seqs)
		defer holders.Stop()
		seekable := o.newSeekable()
		for holders.AllAlive() {
			maxValue, ok := holders.MaxValue()
			if !ok {
				return
			}
			if !holders.AllEqual(maxValue) {
				holders.SeekAll(maxValue)
				continue
			}
			seekable.value = maxValue
			if !yield(seekable) {
				return
			}
			if *seekable.seek {
				holders.SeekAll(*seekable.seekValue)
				*seekable.seek = false
				continue
			}
			holders.Next()
		}
	}
}

func (o Ordered[T]) Union(seqs ...iter.Seq[Seekable[T]]) iter.Seq[Seekable[T]] {
	return func(yield func(Seekable[T]) bool) {
		holders := newSeqHolders(seqs)
		defer holders.Stop()
		seekable := o.newSeekable()
		for !holders.AllStopped() {
			minValue, ok := holders.MinValue()
			if !ok {
				return
			}
			seekable.value = minValue
			if !yield(seekable) {
				return
			}
			if *seekable.seek {
				holders.SeekAll(*seekable.seekValue)
				*seekable.seek = false
				continue
			}
			holders.NextEqual(minValue)
		}
	}
}

type seqHolder[T any] struct {
	Seq      iter.Seq[Seekable[T]]
	PullNext func() (Seekable[T], bool)
	PullStop func()
	Value    Seekable[T]
}

func (sh *seqHolder[T]) Alive() bool {
	return sh.PullNext != nil
}

func (sh *seqHolder[T]) Next() {
	if !sh.Alive() {
		return
	}
	v, ok := sh.PullNext()
	if !ok {
		sh.Stop()
		return
	}
	sh.Value = v
}

func (sh *seqHolder[T]) Seek(to T) {
	if !sh.Alive() {
		return
	}
	sh.Value.Seek(to)
	sh.Next()
}

func (sh *seqHolder[T]) Stop() {
	if !sh.Alive() {
		return
	}
	sh.PullStop()
	sh.PullStop = nil
	sh.PullNext = nil
	sh.Value = Seekable[T]{}
}

type seqHolders[T any] []*seqHolder[T]

func newSeqHolders[T any](seqs []iter.Seq[Seekable[T]]) seqHolders[T] {
	var holders []*seqHolder[T]
	for _, seq := range seqs {
		holder := &seqHolder[T]{
			Seq:   seq,
			Value: Seekable[T]{},
		}
		holder.PullNext, holder.PullStop = iter.Pull(seq)
		holder.Next()
		holders = append(holders, holder)
	}
	return holders
}

func (holders seqHolders[T]) Stop() {
	for _, holder := range holders {
		holder.Stop()
	}
}

func (holders seqHolders[T]) AllAlive() bool {
	for _, holder := range holders {
		if !holder.Alive() {
			return false
		}
	}
	return true
}

func (holders seqHolders[T]) AllStopped() bool {
	for _, holder := range holders {
		if holder.Alive() {
			return false
		}
	}
	return true
}

func (holders seqHolders[T]) MaxValue() (T, bool) {
	var highest T
	first := true
	for _, holder := range holders {
		if !holder.Alive() {
			continue
		}
		v := holder.Value.Value()
		if first {
			highest = v
			first = false
			continue
		}
		if holder.Value.Less()(highest, v) {
			highest = v
		}
	}
	if first {
		return highest, false
	}
	return highest, true
}
func (holders seqHolders[T]) MinValue() (T, bool) {
	var lowest T
	first := true
	for _, holder := range holders {
		if !holder.Alive() {
			continue
		}
		v := holder.Value.Value()
		if first {
			lowest = v
			first = false
			continue
		}
		if holder.Value.Less()(v, lowest) {
			lowest = v
		}
	}
	if first {
		return lowest, false
	}
	return lowest, true
}

func (holders seqHolders[T]) AllEqual(comparValue T) bool {
	for _, holder := range holders {
		if !holder.Alive() {
			return false
		}
		v := holder.Value.Value()
		less := holder.Value.Less()
		if less(comparValue, v) || less(v, comparValue) {
			return false
		}
	}
	return true
}

func (holders seqHolders[T]) NextEqual(comparValue T) {
	for _, holder := range holders {
		if !holder.Alive() {
			continue
		}
		v := holder.Value.Value()
		less := holder.Value.Less()
		if less(comparValue, v) || less(v, comparValue) {
			continue
		}
		holder.Next()
	}
}

func (holders seqHolders[T]) SeekAll(v T) {
	for _, holder := range holders {
		holder.Seek(v)
	}
}

func (holders seqHolders[T]) Next() {
	for _, holder := range holders {
		holder.Next()
	}
}

type SeekConstructor[T any] func(start *T) iter.Seq[T]

func UnSeek[T any](seq iter.Seq[Seekable[T]]) iter.Seq[T] {
	return Map(seq, func(from Seekable[T]) T {
		return from.Value()
	})
}

type Seekable[T any] struct {
	value     T
	seek      *bool
	seekValue *T
	less      func(a, b T) bool
}

func (s Seekable[T]) Value() T {
	return s.value
}

func (s Seekable[T]) Seek(value T) {
	if s.less(value, s.value) {
		return
	}
	*s.seek = true
	*s.seekValue = value
}

func (s Seekable[T]) Less() func(a, b T) bool {
	return s.less
}
