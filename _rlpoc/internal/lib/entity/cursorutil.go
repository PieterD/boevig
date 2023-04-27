package entity

func CursorFromAttribute(s *Store, v Value) Cursor {
	return newCursorFromAttribute(s, v)
}

func CursorFromIds(ids ...Id) Cursor {
	return newCursorFromIds(ids...)
}

func CursorJoin(cursors ...Cursor) Cursor {
	return newCursorFromCursors(cursors...)
}

func CursorVisit(cur Cursor, f func(Id) bool) {
	for {
		entityId, ok := cur.Fix()
		if !ok {
			return
		}
		if f(entityId) == false {
			return
		}
		cur.Advance()
	}
}
