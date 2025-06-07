package ecs

import (
	"fmt"
	"iter"

	"github.com/PieterD/boevig/rang"
	"github.com/google/btree"
)

type cStorePageG[T Component, TP PTRContract[T]] struct {
	tree *btree.BTreeG[tuple[T]]
}

type tuple[VAL any] struct {
	ID  EntityID
	Val VAL
}

func newCStorePageG[T Component, TP PTRContract[T]]() *cStorePageG[T, TP] {
	less := func(a, b tuple[T]) bool {
		if a.ID < b.ID {
			return true
		}
		return false
	}
	return &cStorePageG[T, TP]{
		tree: btree.NewG(5, less),
	}
}

func (cp *cStorePageG[T, TP]) Add(id EntityID, iv Component) {
	v, ok := iv.(T)
	if !ok {
		panic(fmt.Errorf("fetching %v component %T: invalid component type, expected %T", id, iv, v))
	}
	cp.AddG(id, v)
}

func (cp *cStorePageG[T, TP]) AddG(id EntityID, v T) {
	cp.tree.ReplaceOrInsert(tuple[T]{ID: id, Val: v})
}

func (cp *cStorePageG[T, TP]) Remove(id EntityID) {
	cp.tree.Delete(tuple[T]{ID: id})
}

func (cp *cStorePageG[T, TP]) Get(id EntityID, iv Component) bool {
	vp, ok := iv.(TP)
	if !ok {
		panic(fmt.Errorf("fetching %v component %T: invalid component type, expected %T", id, iv, vp))
	}
	v, ok := cp.GetG(id)
	if !ok {
		return false
	}
	*vp = v
	return true
}

func (cp *cStorePageG[T, TP]) GetG(id EntityID) (T, bool) {
	var zero T
	got, ok := cp.tree.Get(tuple[T]{ID: id})
	if !ok {
		return zero, false
	}
	return got.Val, true
}

func (cp *cStorePageG[T, TP]) SeekSeq() iter.Seq[rang.Seekable[EntityID]] {
	idLess := func(a, b EntityID) bool {
		return a < b
	}
	return rang.NewOrdered[EntityID](idLess).SeekIterator(func(start *EntityID) iter.Seq[EntityID] {
		return func(yield func(EntityID) bool) {
			if start == nil {
				first, ok := cp.tree.Min()
				if !ok {
					return
				}
				start = &first.ID
			}
			cp.tree.AscendGreaterOrEqual(tuple[T]{ID: *start}, func(item tuple[T]) bool {
				if !yield(item.ID) {
					return false
				}
				return true
			})
		}
	})
}

func (cp *cStorePageG[T, TP]) AllG() iter.Seq2[EntityID, T] {
	return func(yield func(EntityID, T) bool) {
		cp.tree.Ascend(func(item tuple[T]) bool {
			if !yield(item.ID, item.Val) {
				return false
			}
			return true
		})
	}
}
