package attrib

import (
	"fmt"

	"github.com/PieterD/rlpoc/lib/cardinal"
	"github.com/PieterD/rlpoc/lib/entity"
)

type Room struct {
	Map string
}

func (a Room) Less(value entity.Value) bool {
	return false
}

var _ entity.Value = Room{}

type Player struct{}

func (a Player) Less(value entity.Value) bool {
	return false
}

var _ entity.Value = Player{}

type Location struct {
	Room     entity.Id
	Position cardinal.Coord
}

func (a Location) Less(value entity.Value) bool {
	return false
}

var _ entity.Value = Location{}

type Event struct {
	Depth  uint64
	Action string
}

func (a Event) Less(value entity.Value) bool {
	return false
}

var _ entity.Value = Event{}

type EventParentLink struct {
	Parent entity.Id
}

func (a EventParentLink) Less(value entity.Value) bool {
	b, ok := value.(*EventParentLink)
	if !ok {
		panic(fmt.Errorf("invalid type: want %T, got %T", b, value))
	}
	return a.Parent.Less(b.Parent)
}

var _ entity.Value = EventParentLink{}

type EventChildLink struct {
	Child entity.Id
}

func (a EventChildLink) Less(value entity.Value) bool {
	b, ok := value.(*EventChildLink)
	if !ok {
		panic(fmt.Errorf("invalid type: want %T, got %T", b, value))
	}
	return a.Child.Less(b.Child)
}

var _ entity.Value = EventChildLink{}

type EventRecentParent struct {
	ParentId entity.Id
}

func (a EventRecentParent) Less(value entity.Value) bool {
	return false
}

var _ entity.Value = EventRecentParent{}

type EventRecentChild struct {
	ChildId entity.Id
}

func (a EventRecentChild) Less(value entity.Value) bool {
	return false
}

var _ entity.Value = EventRecentChild{}
