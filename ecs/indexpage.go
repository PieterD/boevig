package ecs

import (
	"fmt"
	"iter"
	"sort"

	"github.com/PieterD/boevig/rang"
	"github.com/samber/lo"
)

type indexPageG[T comparable] struct {
	idToValue  map[EntityID]T
	valueToIDs map[T]map[EntityID]struct{}
}

func newIndexPageG[T comparable]() *indexPageG[T] {
	return &indexPageG[T]{
		idToValue:  make(map[EntityID]T),
		valueToIDs: make(map[T]map[EntityID]struct{}),
	}
}

func (page *indexPageG[T]) Set(id EntityID, vi any) {
	v, ok := vi.(T)
	if !ok {
		panic(fmt.Errorf("adding to field index %v %+v: invalid type, got %T, want %T", id, vi, vi, v))
	}
	page.SetG(id, v)
}

func (page *indexPageG[T]) SetG(id EntityID, value T) {
	existingValue, ok := page.idToValue[id]
	if ok {
		if existingValue == value {
			return
		}
		delete(page.idToValue, id)
		delete(page.valueToIDs[existingValue], id)
	}
	page.idToValue[id] = value
	if _, ok := page.valueToIDs[value]; !ok {
		page.valueToIDs[value] = make(map[EntityID]struct{})
	}
	page.valueToIDs[value][id] = struct{}{}
}

func (page *indexPageG[T]) Remove(id EntityID) {
	value, ok := page.idToValue[id]
	if !ok {
		return
	}
	delete(page.valueToIDs[value], id)
	if len(page.valueToIDs[value]) == 0 {
		delete(page.valueToIDs, value)
	}
}

func (page *indexPageG[T]) SeekSeq(vi any) iter.Seq[rang.Seekable[EntityID]] {
	v, ok := vi.(T)
	if !ok {
		panic(fmt.Errorf("seekseq on field index %+v: invalid type, got %T, want %T", vi, vi, v))
	}
	return page.SeekSeqG(v)
}

func (page *indexPageG[T]) SeekSeqG(value T) iter.Seq[rang.Seekable[EntityID]] {
	return rang.NewOrdered(func(a, b EntityID) bool { return a < b }).SeekIterator(func(start *EntityID) iter.Seq[EntityID] {
		return func(yield func(EntityID) bool) {
			m := page.valueToIDs[value]
			if len(m) == 0 {
				return
			}
			ids := lo.Keys(m)
			sort.Slice(ids, func(i, j int) bool {
				return ids[i] < ids[j]
			})
			id := ids[0]
			if start != nil {
				id = *start
			}
			n := sort.Search(len(ids), func(i int) bool { return ids[i] >= id })
			if n == len(ids) {
				return
			}
			for i := n; i < len(ids); i++ {
				if !yield(ids[i]) {
					return
				}
			}
		}
	})
}
