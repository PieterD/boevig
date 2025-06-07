package rang

import (
	"iter"
	"testing"

	"github.com/google/btree"
	"github.com/stretchr/testify/require"
)

func TestOrdered_SeekIterator(t *testing.T) {
	less := func(a, b int) bool {
		return a < b
	}
	t.Run("simple", func(t *testing.T) {
		searcher := NewTestSearcher(-2, -1, 1, 2, 3, 4, 5, 6)
		seq := NewOrdered(less).SeekIterator(searcher.Search)
		require.Equal(t, []int{-2, -1, 1, 2, 3, 4, 5, 6}, ToSlice(UnSeek(seq)))
	})
	t.Run("break", func(t *testing.T) {
		searcher := NewTestSearcher(-2, -1, 1, 2, 3, 4, 5, 6)
		seq := NewOrdered(less).SeekIterator(searcher.Search)
		var got []int
		for v := range seq {
			if v.Value() >= 4 {
				break
			}
			got = append(got, v.Value())
		}
		require.Equal(t, []int{-2, -1, 1, 2, 3}, got)
	})
	t.Run("seek", func(t *testing.T) {
		searcher := NewTestSearcher(-2, -1, 1, 2, 3, 4, 5, 6)
		seq := NewOrdered(less).SeekIterator(searcher.Search)
		var got []int
		seekOnce := false
		for v := range seq {
			if v.Value() < 4 {
				require.False(t, seekOnce)
				seekOnce = true
				v.Seek(4)
				continue
			}
			got = append(got, v.Value())
		}
		require.True(t, seekOnce)
		require.Equal(t, []int{4, 5, 6}, got)
	})
	t.Run("seek back fails", func(t *testing.T) {
		searcher := NewTestSearcher(-2, -1, 1, 2, 3, 4, 5, 6)
		seq := NewOrdered(less).SeekIterator(searcher.Search)
		var got []int
		seekOnce := false
		for v := range seq {
			if v.Value() == 4 {
				require.False(t, seekOnce)
				seekOnce = true
				v.Seek(2)
				continue
			}
			got = append(got, v.Value())
		}
		require.True(t, seekOnce)
		require.Equal(t, []int{-2, -1, 1, 2, 3, 5, 6}, got)
	})
}

func TestOrdered_Collections(t *testing.T) {
	less := func(a, b int) bool {
		return a < b
	}
	tests := []struct {
		desc        string
		f           func(seqs ...iter.Seq[Seekable[int]]) iter.Seq[Seekable[int]]
		left        []int
		right       []int
		seekAt      int
		seekTo      int
		exIntersect []int
		exUnion     []int
	}{
		{
			desc:        "no overlap",
			left:        []int{1, 2, 3},
			right:       []int{4, 5, 6},
			exIntersect: []int(nil),
			exUnion:     []int{1, 2, 3, 4, 5, 6},
		},
		{
			desc:        "left empty",
			left:        []int(nil),
			right:       []int{4, 5, 6},
			exIntersect: []int(nil),
			exUnion:     []int{4, 5, 6},
		},
		{
			desc:        "right empty",
			left:        []int{4, 5, 6},
			right:       []int(nil),
			exIntersect: []int(nil),
			exUnion:     []int{4, 5, 6},
		},
		{
			desc:        "both empty",
			left:        []int(nil),
			right:       []int(nil),
			exIntersect: []int(nil),
			exUnion:     []int(nil),
		},
		{
			desc:        "overlap all",
			left:        []int{1, 2, 3},
			right:       []int{1, 2, 3},
			exIntersect: []int{1, 2, 3},
			exUnion:     []int{1, 2, 3},
		},
		{
			desc:        "overlap middle",
			left:        []int{1, 2, 3, 4},
			right:       []int{3, 4, 5, 6},
			exIntersect: []int{3, 4},
			exUnion:     []int{1, 2, 3, 4, 5, 6},
		},
		{
			desc:        "overlap big left",
			left:        []int{1, 2, 3, 4, 5, 6, 7, 8},
			right:       []int{6},
			exIntersect: []int{6},
			exUnion:     []int{1, 2, 3, 4, 5, 6},
		},
		{
			desc:        "overlap big right",
			left:        []int{6},
			right:       []int{1, 2, 3, 4, 5, 6, 7, 8},
			exIntersect: []int{6},
			exUnion:     []int{1, 2, 3, 4, 5, 6, 7, 8},
		},
	}
	for _, test := range tests {
		test := test
		o := NewOrdered(less)
		type TC struct {
			desc string
			f    func(seqs ...iter.Seq[Seekable[int]]) iter.Seq[Seekable[int]]
			ex   []int
		}
		tcs := []TC{
			{
				desc: "intersect",
				f:    o.Intersect,
				ex:   test.exIntersect,
			},
			{
				desc: "union",
				f:    o.Union,
				ex:   test.exUnion,
			},
		}
		for _, tc := range tcs {
			t.Run(tc.desc, func(t *testing.T) {

			})
		}
		t.Run(test.desc, func(t *testing.T) {
			for _, tc := range tcs {
				tc := tc
				t.Run(tc.desc, func(t *testing.T) {
					leftSeq := o.SeekIterator(NewTestSearcher(test.left...).Search)
					rightSeq := o.SeekIterator(NewTestSearcher(test.right...).Search)
					interSeq := o.Intersect(leftSeq, rightSeq)
					var vs []int
					for v := range interSeq {
						if test.seekAt != 0 && v.Value() == test.seekAt {
							v.Seek(test.seekTo)
							continue
						}
						vs = append(vs, v.Value())
					}
					require.Equal(t, test.exIntersect, vs)
				})
			}
		})
	}
}

type TestSearcher struct {
	tree *btree.BTreeG[int]
}

func NewTestSearcher(values ...int) *TestSearcher {
	tree := btree.NewG[int](5, func(a, b int) bool {
		return a < b
	})
	for _, v := range values {
		tree.ReplaceOrInsert(v)
	}
	return &TestSearcher{
		tree: tree,
	}
}

func (it *TestSearcher) Search(first *int) iter.Seq[int] {
	return func(yield func(int) bool) {
		visitor := func(item int) bool {
			if !yield(item) {
				return false
			}
			return true
		}
		if first == nil {
			it.tree.Ascend(visitor)
			return
		}
		it.tree.AscendGreaterOrEqual(*first, visitor)
	}
}
