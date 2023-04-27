package entity

import (
	"fmt"
	"sort"
)

type Cursor interface {
	Fix() (Id, bool)
	Advance()
	Seek(Id)
}

func newMultiAttributeCursor(s *Store, vs ...Value) Cursor {
	var cursors []Cursor
	for _, v := range vs {
		cursor := CursorFromAttribute(s, v)
		cursors = append(cursors, cursor)
	}
	return CursorJoin(cursors...)
}

type cursorFixedIds struct {
	ids   []Id
	index int
	damp  bool
}

func newCursorFromIds(ids ...Id) *cursorFixedIds {
	return &cursorFixedIds{
		ids: ids,
	}
}

func (c *cursorFixedIds) Fix() (Id, bool) {
	if c.index >= len(c.ids) {
		return Id{}, false
	}
	c.damp = false
	return c.ids[c.index], true
}

func (c *cursorFixedIds) Advance() {
	if c.damp {
		panic(fmt.Errorf("cursorMultiAttr is already damp, Fix it first"))
	}
	c.index++
	c.damp = true
}

func (c *cursorFixedIds) Seek(id Id) {
	newIndex := sort.Search(len(c.ids), func(i int) bool {
		return c.ids[i].Less(id)
	})
	c.index = newIndex
	c.damp = true
}

var _ Cursor = &cursorFixedIds{}

type cursorJoin struct {
	cursors []Cursor
	damp    bool
}

func (c *cursorJoin) Fix() (Id, bool) {
	if len(c.cursors) == 0 {
		return Id{}, false
	}
	c.damp = false
	for {
		maxId := Id{}
		maxCount := 0
		for i, cursor := range c.cursors {
			id, ok := cursor.Fix()
			if !ok {
				return Id{}, false
			}
			if i == 0 || maxId.Less(id) {
				maxId = id
				maxCount = 1
				continue
			}
			if maxId == id {
				maxCount++
				continue
			}
			// maxId > id: we do not change maxCount.
			// one of the cursors does not have this id.
		}
		if maxCount == len(c.cursors) {
			return maxId, true
		}
		for _, cursor := range c.cursors {
			cursor.Seek(maxId)
		}
	}
}

func (c *cursorJoin) Advance() {
	if c.damp {
		panic(fmt.Errorf("cursorMultiAttr is already damp, Fix it first"))
	}
	for _, cursor := range c.cursors {
		cursor.Advance()
	}
	c.damp = true
}

func (c *cursorJoin) Seek(id Id) {
	for _, cursor := range c.cursors {
		cursor.Seek(id)
	}
	c.damp = true
}

func newCursorFromCursors(cursors ...Cursor) *cursorJoin {
	return &cursorJoin{
		cursors: cursors,
	}
}

var _ Cursor = &cursorJoin{}

type cursorAttrib struct {
	v        Value
	store    *attribStore
	entityId Id
	damp     bool
}

func newCursorFromAttribute(s *Store, v Value) Cursor {
	_, attribStore := s.mustGetAttribStore(v)
	firstEntityId, ok := attribStore.FirstEntity()
	if !ok {
		return &cursorEmpty{}
	}
	return &cursorAttrib{
		v:        v,
		store:    attribStore,
		entityId: firstEntityId,
	}
}

func (c *cursorAttrib) Fix() (Id, bool) {
	id := c.entityId
	if c.store.FindEntity(&id) == false {
		return Id{}, false
	}
	c.entityId = id
	c.damp = false
	return id, true
}

func (c *cursorAttrib) Advance() {
	if c.damp {
		panic(fmt.Errorf("cursorAttr is already damp, Fix it first"))
	}
	c.entityId = c.entityId.next()
	c.damp = true
}

func (c *cursorAttrib) Seek(id Id) {
	c.entityId = id
	c.damp = true
}

var _ Cursor = &cursorAttrib{}

type cursorEmpty struct{}

func (c cursorEmpty) Fix() (Id, bool) {
	return Id{}, false
}

func (c cursorEmpty) Advance() {}

func (c cursorEmpty) Seek(id Id) {}

var _ Cursor = cursorEmpty{}
