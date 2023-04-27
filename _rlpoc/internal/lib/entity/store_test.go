package entity_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/PieterD/rlpoc/lib/entity"
)

func testStore(t *testing.T) (*entity.Store, []entity.Id) {
	t.Helper()
	s := entity.NewStore()
	alice, bob, charlie := s.NewEntityId(), s.NewEntityId(), s.NewEntityId()
	s.Give(alice, &Name{"Alice"}, &Country{"Netherlands"})
	s.Give(bob, &Name{"Bob"}, &Country{"Netherlands"}, &Tag{"female"}, &Tag{"married"})
	s.Give(charlie, &Name{"Charlie"}, &Tag{"male"})
	return s, []entity.Id{alice, bob, charlie}
}

func TestStore_Entities(t *testing.T) {
	tests := []struct {
		desc  string
		attrs []entity.Value
		want  []int
	}{
		{
			desc:  "name only",
			attrs: []entity.Value{&Name{}},
			want:  []int{0, 1, 2},
		},
		{
			desc:  "country only",
			attrs: []entity.Value{&Country{}},
			want:  []int{0, 1},
		},
		{
			desc:  "tag only",
			attrs: []entity.Value{&Tag{}},
			want:  []int{1, 2},
		},
		{
			desc:  "country and tag",
			attrs: []entity.Value{&Country{}, &Tag{}},
			want:  []int{1},
		},
		{
			desc: "name and tag",
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			s, entities := testStore(t)
			cursor := s.Entities(test.attrs...)
			got := unpackCursor(t, entities, cursor)
			if !reflect.DeepEqual(test.want, got) {
				t.Logf("want: %v", test.want)
				t.Logf("got : %v", got)
				t.Errorf("invalid result set")
			}
		})
	}
}

func unpackCursor(t *testing.T, entities []entity.Id, cursor entity.Cursor) []int {
	t.Helper()
	entityMap := make(map[entity.Id]int)
	for i, e := range entities {
		entityMap[e] = i
	}
	var result []int
	for {
		entityId, ok := cursor.Fix()
		if !ok {
			break
		}
		entityIndex, ok := entityMap[entityId]
		if !ok {
			t.Fatalf("unknown entity id: %v", entityId)
		}
		result = append(result, entityIndex)
		cursor.Advance()
	}
	return result
}

type Name struct {
	content string
}

func (a Name) Less(value entity.Value) bool {
	return false
}

var _ entity.Value = &Name{}

type Country struct {
	content string
}

func (a Country) Less(value entity.Value) bool {
	return false
}

var _ entity.Value = &Country{}

type Tag struct {
	content string
}

func (a Tag) Less(value entity.Value) bool {
	b, ok := value.(*Tag)
	if !ok {
		panic(fmt.Errorf("invalid value: want %T, got %T", b, value))
	}
	return a.content < b.content
}

var _ entity.Value = &Tag{}
