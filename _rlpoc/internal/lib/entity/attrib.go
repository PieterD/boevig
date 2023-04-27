package entity

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/google/btree"
)

type attribNode struct {
	entityId Id
	value    Value
}

func (a *attribNode) Less(than btree.Item) bool {
	b, ok := than.(*attribNode)
	if !ok {
		panic(fmt.Errorf("invalid type for than: want %T, got %T", b, than))
	}
	if a.entityId.Less(b.entityId) {
		return true
	}
	if b.entityId.Less(a.entityId) {
		return false
	}
	// if both values are nil, consider the two nodes equal.
	if a.value == nil && b.value == nil {
		return false
	}
	// nil goes left for easy searching.
	if a.value == nil {
		return true
	} else if b.value == nil {
		return false
	}
	return a.value.Less(b.value)
}

var _ btree.Item = &attribNode{}

type attribStore struct {
	btree *btree.BTree // *attribNode
}

func newAttribStore() *attribStore {
	return &attribStore{
		btree: btree.New(50),
	}
}

func (s *attribStore) Give(entityId Id, v Value) {
	node := &attribNode{
		entityId: entityId,
		value:    v,
	}
	s.btree.ReplaceOrInsert(node)
}

func (s *attribStore) Destroy(entityId Id) {
	node := &attribNode{
		entityId: entityId,
	}
	s.btree.AscendGreaterOrEqual(node, func(i btree.Item) bool {
		found, ok := i.(*attribNode)
		if !ok {
			panic(fmt.Errorf("invalid item type in btree: want %T, got %T", found, i))
		}
		if found.entityId != entityId {
			return false
		}
		s.btree.Delete(found)
		return true
	})
}

func (s *attribStore) Remove(entityId Id, v Value) {
	node := &attribNode{
		entityId: entityId,
		value:    v,
	}
	s.btree.Delete(node)
}

func (s *attribStore) Clear() {
	s.btree.Clear(false)
}

var errAttributeNotFound = fmt.Errorf("attribute not found")

func IsAttributeNotFoundError(err error) bool {
	return errors.Is(err, errAttributeNotFound)
}

func overwrite(to, from Value) {
	vTo := reflect.ValueOf(to)
	vFrom := reflect.ValueOf(from)
	if vTo.Kind() != reflect.Ptr {
		panic(fmt.Errorf("destination should be a pointer, not %T", to))
	}
	vTo = vTo.Elem()
	if vFrom.Kind() == reflect.Ptr {
		vFrom = vFrom.Elem()
	}
	if vTo.Type() != vFrom.Type() {
		panic(fmt.Errorf("expected identical types, not %T and %T", to, from))
	}
	vTo.Set(vFrom)
}

func (s *attribStore) FirstAttrib(entityId Id, v Value) bool {
	start := &attribNode{entityId: entityId, value: nil}
	success := false
	s.btree.AscendGreaterOrEqual(start, func(i btree.Item) bool {
		found, ok := i.(*attribNode)
		if !ok {
			panic(fmt.Errorf("invalid item type in btree: want %T, got %T", found, i))
		}
		if found.entityId != entityId {
			return false
		}
		success = true
		overwrite(v, found.value)
		return false
	})
	if !success {
		return false
	}
	return true
}

func (s *attribStore) NextAttrib(entityId Id, v Value) bool {
	start := &attribNode{entityId: entityId, value: v}
	first := true
	success := false
	s.btree.AscendGreaterOrEqual(start, func(i btree.Item) bool {
		found, ok := i.(*attribNode)
		if !ok {
			panic(fmt.Errorf("invalid item type in btree: want %T, got %T", found, i))
		}
		if found.entityId != entityId {
			return false
		}
		if first {
			first = false
			// Check if the first value exists so we skip to the next.
			if !v.Less(found.value) && !found.value.Less(v) {
				// return true, request the next element
				return true
			}
		}
		success = true
		overwrite(v, found.value)
		return false
	})
	if !success {
		return false
	}
	return true
}

func (s *attribStore) FirstEntity() (Id, bool) {
	searchId := Id{}
	if s.FindEntity(&searchId) == false {
		return Id{}, false
	}
	return searchId, true
}

func (s *attribStore) FindEntity(entityId *Id) bool {
	id := *entityId
	pivot := &attribNode{
		entityId: id,
	}
	success := false
	s.btree.AscendGreaterOrEqual(pivot, func(i btree.Item) bool {
		found, ok := i.(*attribNode)
		if !ok {
			panic(fmt.Errorf("invalid item type in btree: want %T, got %T", found, i))
		}
		success = true
		*entityId = found.entityId
		return false
	})
	if !success {
		return false
	}
	return true
}
