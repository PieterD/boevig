package ecs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TestComponentIndex struct {
	ComponentHeader[TestComponentIndex, *TestComponentIndex]
	String string
	Num    int
	Bool   bool
}

func (c TestComponentIndex) Index() []Indexer {
	return []Indexer{
		EQ("test_str", c.String),
		EQ("test_num", c.Num),
		EQ("test_bool", c.Bool),
	}
}

func TestDB_Get(t *testing.T) {
	t.Run("get single component", func(t *testing.T) {
		db, ids := dbDefaults()
		str := TestComponentString{}
		require.True(t, db.Get(ids[1], &str))
		require.Equal(t, "string_1", str.String)
	})
	t.Run("get double component", func(t *testing.T) {
		db, ids := dbDefaults()
		num := TestComponentNum{}
		idx := TestComponentIndex{}
		require.True(t, db.Get(ids[3], &idx, &num))
		require.Equal(t, "indexed_string_3", idx.String)
		require.Equal(t, 3, num.Int)
	})
	t.Run("get one missing component", func(t *testing.T) {
		db, ids := dbDefaults()
		num := TestComponentNum{}
		require.False(t, db.Get(ids[2], &num))
	})
	t.Run("get one missing component, one existing", func(t *testing.T) {
		db, ids := dbDefaults()
		num := TestComponentNum{}
		str := TestComponentString{}
		require.False(t, db.Get(ids[2], &num, &str))
		require.False(t, db.Get(ids[3], &num, &str))
	})
}

func TestDB_Search(t *testing.T) {
	type ptrs struct {
		str TestComponentString
		num TestComponentNum
		idx TestComponentIndex
		boo TestComponentBool
	}
	tests := []struct {
		desc            string
		f               func(b *SearchBuilder, ptrs *ptrs) *SearchBuilder
		expectedIds     []EntityID
		expectedStrs    []TestComponentString
		expectedNums    []TestComponentNum
		expectedIndexes []TestComponentIndex
		expectedBools   []TestComponentBool
	}{
		{
			desc: "components - one component one result",
			f: func(b *SearchBuilder, ptrs *ptrs) *SearchBuilder {
				return b.Components(&ptrs.boo)
			},
			expectedIds: []EntityID{4},
			expectedBools: []TestComponentBool{
				{Bool: true},
			},
		},
		{
			desc: "components - one component multiple results",
			f: func(b *SearchBuilder, ptrs *ptrs) *SearchBuilder {
				return b.Components(&ptrs.str)
			},
			expectedIds: []EntityID{1, 2, 4},
			expectedStrs: []TestComponentString{
				{String: "string_1"},
				{String: "string_2"},
				{String: "string_4"},
			},
		},
		{
			desc: "components - two components one result",
			f: func(b *SearchBuilder, ptrs *ptrs) *SearchBuilder {
				return b.Components(&ptrs.str, &ptrs.idx)
			},
			expectedIds: []EntityID{1},
			expectedStrs: []TestComponentString{
				{String: "string_1"},
			},
			expectedIndexes: []TestComponentIndex{
				{String: "indexed_string_1", Num: 1, Bool: true},
			},
		},
		{
			desc: "components - two components multiple results",
			f: func(b *SearchBuilder, ptrs *ptrs) *SearchBuilder {
				return b.Components(&ptrs.num, &ptrs.idx)
			},
			expectedIds: []EntityID{1, 3},
			expectedNums: []TestComponentNum{
				{Int: 1},
				{Int: 3},
			},
			expectedIndexes: []TestComponentIndex{
				{String: "indexed_string_1", Num: 1, Bool: true},
				{String: "indexed_string_3", Num: 3, Bool: false},
			},
		},
		{
			desc: "index - one index one result",
			f: func(b *SearchBuilder, ptrs *ptrs) *SearchBuilder {
				return b.Index((EQ("test_bool", false)))
			},
			expectedIds: []EntityID{3},
		},
		{
			desc: "index - one index multiple results",
			f: func(b *SearchBuilder, ptrs *ptrs) *SearchBuilder {
				return b.Index((EQ("test_bool", true)))
			},
			expectedIds: []EntityID{1, 5},
		},
		{
			desc: "index - limited by component",
			f: func(b *SearchBuilder, ptrs *ptrs) *SearchBuilder {
				return b.
					Index((EQ("test_bool", true))).
					Components(&ptrs.num)
			},
			expectedIds: []EntityID{1},
			expectedNums: []TestComponentNum{
				{Int: 1},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			db, _ := dbDefaults()
			ptrs := &ptrs{}
			seq := test.f(db.Search(), ptrs).Done()
			var ids []EntityID
			var strs []TestComponentString
			var idxs []TestComponentIndex
			for id := range seq {
				ids = append(ids, id)
				if test.expectedStrs != nil {
					strs = append(strs, ptrs.str)
				}
				if test.expectedIndexes != nil {
					idxs = append(idxs, ptrs.idx)
				}
			}
			require.Equal(t, test.expectedIds, ids)
			require.Equal(t, test.expectedStrs, strs)
			require.Equal(t, test.expectedIndexes, idxs)
		})
	}
}

func dbDefaults() (*DB, []EntityID) {
	db := New()
	ids := []EntityID{0}
	ids = append(ids, db.NewEntity(
		TestComponentIndex{String: "indexed_string_1", Num: 1, Bool: true},
		TestComponentString{String: "string_1"},
		TestComponentNum{Int: 1}))
	ids = append(ids, db.NewEntity(
		TestComponentString{String: "string_2"}))
	ids = append(ids, db.NewEntity(
		TestComponentIndex{String: "indexed_string_3", Num: 3, Bool: false},
		TestComponentNum{Int: 3}))
	ids = append(ids, db.NewEntity(
		TestComponentString{String: "string_4"},
		TestComponentBool{Bool: true}))
	ids = append(ids, db.NewEntity(
		TestComponentIndex{String: "indexed_string_5", Num: 5, Bool: true}))
	return db, ids
}
