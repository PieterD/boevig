package ecs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TestComponentString struct {
	ComponentHeader[TestComponentString, *TestComponentString]
	String string
}

type TestComponentNum struct {
	ComponentHeader[TestComponentNum, *TestComponentNum]
	Int int
}

type TestComponentFloat struct {
	ComponentHeader[TestComponentFloat, *TestComponentFloat]
	Float float64
}

type TestComponentBool struct {
	ComponentHeader[TestComponentBool, *TestComponentBool]
	Bool bool
}

func TestCStoreBook_Get(t *testing.T) {
	t.Run("get existing component one at a time", func(t *testing.T) {
		book, tcs, _, _ := cStoreBookDefaults()
		ok := book.Get(1, &tcs)
		require.True(t, ok)
		require.Equal(t, "c1", tcs.String)
	})
	t.Run("get existing component two at a time", func(t *testing.T) {
		book, tcs, tci, _ := cStoreBookDefaults()
		ok := book.Get(3, &tcs, &tci)
		require.True(t, ok)
		require.Equal(t, "c3", tcs.String)
		require.Equal(t, 3, tci.Int)
	})
	t.Run("get existing component three at a time", func(t *testing.T) {
		book, tcs, tci, tcf := cStoreBookDefaults()
		ok := book.Get(7, &tcs, &tci, &tcf)
		require.True(t, ok)
		require.Equal(t, "c7", tcs.String)
		require.Equal(t, 7, tci.Int)
		require.Equal(t, 7.0, tcf.Float)
	})
	t.Run("get single missing component", func(t *testing.T) {
		book, tcs, _, _ := cStoreBookDefaults()
		ok := book.Get(2, &tcs)
		require.False(t, ok)
	})
	t.Run("get one missing component out of two", func(t *testing.T) {
		book, tcs, tci, _ := cStoreBookDefaults()
		ok := book.Get(2, &tcs, &tci)
		require.False(t, ok)
	})
	t.Run("get missing entity", func(t *testing.T) {
		book, tcs, tci, _ := cStoreBookDefaults()
		ok := book.Get(4, &tcs, &tci)
		require.False(t, ok)
	})
	t.Run("get missing entity", func(t *testing.T) {
		book, _, _, _ := cStoreBookDefaults()
		var tcb TestComponentBool
		ok := book.Get(4, &tcb)
		require.False(t, ok)
	})
}

func TestCStoreBook_Remove(t *testing.T) {
	t.Run("remove entity with single existing component", func(t *testing.T) {
		book, tcs, _, _ := cStoreBookDefaults()
		book.Remove(1)
		ok := book.Get(1, &tcs)
		require.False(t, ok)
	})
	t.Run("remove entity with two components", func(t *testing.T) {
		book, tcs, tci, _ := cStoreBookDefaults()
		book.Remove(3)
		require.False(t, book.Get(3, &tcs))
		require.False(t, book.Get(3, &tci))
	})
	t.Run("remove entity with three components", func(t *testing.T) {
		book, tcs, tci, tcf := cStoreBookDefaults()
		book.Remove(7)
		require.False(t, book.Get(7, &tcs))
		require.False(t, book.Get(7, &tci))
		require.False(t, book.Get(7, &tcf))
	})
	t.Run("remove entity with no components", func(t *testing.T) {
		book, _, _, _ := cStoreBookDefaults()
		book.Remove(7)
	})
}

func TestCStoreBook_All(t *testing.T) {
	t.Run("get all s elements", func(t *testing.T) {
		book, tcs, _, _ := cStoreBookDefaults()
		var ids []EntityID
		var values []string
		for id := range book.All(&tcs) {
			ids = append(ids, id.Value())
			values = append(values, tcs.String)
		}
		require.Equal(t, []EntityID{1, 3, 5, 7}, ids)
		require.Equal(t, []string{"c1", "c3", "c5", "c7"}, values)
	})
	t.Run("get all i elements", func(t *testing.T) {
		book, _, tci, _ := cStoreBookDefaults()
		var ids []EntityID
		var values []int
		for id := range book.All(&tci) {
			ids = append(ids, id.Value())
			values = append(values, tci.Int)
		}
		require.Equal(t, []EntityID{2, 3, 6, 7}, ids)
		require.Equal(t, []int{2, 3, 6, 7}, values)
	})
	t.Run("get all i and f elements", func(t *testing.T) {
		book, _, tci, tcf := cStoreBookDefaults()
		var ids []EntityID
		var iValues []int
		var fValues []float64
		for id := range book.All(&tci, &tcf) {
			ids = append(ids, id.Value())
			iValues = append(iValues, tci.Int)
			fValues = append(fValues, tcf.Float)
		}
		require.Equal(t, []EntityID{6, 7}, ids)
		require.Equal(t, []int{6, 7}, iValues)
		require.Equal(t, []float64{6, 7}, fValues)
	})
	t.Run("get nonexistent element", func(t *testing.T) {
		book, _, _, _ := cStoreBookDefaults()
		var tcb TestComponentBool
		var ids []EntityID
		for id := range book.All(&tcb) {
			ids = append(ids, id.Value())
		}
		require.Equal(t, []EntityID(nil), ids)
	})
}

func cStoreBookDefaults() (*cStoreBook, TestComponentString, TestComponentNum, TestComponentFloat) {
	book := newCStoreBook()
	book.Add(1, TestComponentString{String: "c1"})
	book.Add(2, TestComponentNum{Int: 2})
	book.Add(3, TestComponentString{String: "c3"}, TestComponentNum{Int: 3})
	book.Add(4, TestComponentFloat{Float: 4})
	book.Add(5, TestComponentFloat{Float: 5}, TestComponentString{String: "c5"})
	book.Add(6, TestComponentFloat{Float: 6}, TestComponentNum{Int: 6})
	book.Add(7, TestComponentFloat{Float: 7}, TestComponentString{String: "c7"}, TestComponentNum{Int: 7})
	return book, TestComponentString{}, TestComponentNum{}, TestComponentFloat{}
}
