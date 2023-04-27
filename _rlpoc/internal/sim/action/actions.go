package action

import (
	"fmt"
	"reflect"

	"github.com/PieterD/rlpoc/lib/entity"

	"github.com/PieterD/rlpoc/internal/sim/attrib"

	"github.com/PieterD/rlpoc/internal/sim"
	"github.com/PieterD/rlpoc/lib/cardinal"
)

var _ sim.ActionConstructor = NewAction

func NewAction(name string, settable bool) (sim.Action, error) {
	constructor, ok := actionConstructors[name]
	if !ok {
		return nil, fmt.Errorf("unknown action name: %s", name)
	}
	action := constructor()
	if !settable {
		v := reflect.ValueOf(action)
		action = v.Elem().Interface().(sim.Action)
	}
	return constructor(), nil
}

var actionConstructors = map[string]func() sim.Action{
	SpawnRoom{}.Name():   func() sim.Action { return &SpawnRoom{} },
	SpawnPlayer{}.Name(): func() sim.Action { return &SpawnPlayer{} },
	Walk{}.Name():        func() sim.Action { return &Walk{} },
}

type SpawnRoom struct {
	EntityId sim.EntityId
	Shape    string
}

func (s SpawnRoom) Name() string {
	return "SpawnRoom"
}

func (s SpawnRoom) Less(value entity.Value) bool {
	return false
}

func (s SpawnRoom) Apply(state *sim.State) error {
	roomId := s.EntityId
	state.EntityStore.Give(roomId, &attrib.Room{Map: s.Shape})
	return nil
}

func (s SpawnRoom) Revert(state *sim.State) error {
	state.EntityStore.DestroyEntity(s.EntityId)
	return nil
}

type SpawnPlayer struct {
	EntityId sim.EntityId
	Room     sim.EntityId
	Position cardinal.Coord
}

func (s SpawnPlayer) Name() string {
	return "SpawnPlayer"
}

func (s SpawnPlayer) Less(value entity.Value) bool {
	return false
}

func (s SpawnPlayer) Apply(state *sim.State) error {
	player := &attrib.Player{}
	location := &attrib.Location{
		Room:     s.Room,
		Position: s.Position,
	}
	state.EntityStore.Clear(player)
	state.EntityStore.Give(s.EntityId, player, location)
	return nil
}

func (s SpawnPlayer) Revert(state *sim.State) error {
	state.EntityStore.DestroyEntity(s.EntityId)
	return nil
}

type Walk struct {
	EntityId  sim.EntityId
	Direction cardinal.Direction
}

func (a Walk) Name() string {
	return "Walk"
}

func (a Walk) Less(value entity.Value) bool {
	return false
}

func (a Walk) Apply(state *sim.State) error {
	player := &attrib.Player{}
	location := &attrib.Location{}
	playerId, ok := state.EntityStore.Entities(player).Fix()
	if !ok {
		return fmt.Errorf("no player found")
	}
	if state.EntityStore.First(playerId, location) == false {
		return fmt.Errorf("player has no location")
	}
	a.Direction.Move(&location.Position)
	state.EntityStore.Give(playerId, location)
	return nil
}

func (a Walk) Revert(state *sim.State) error {
	location := &attrib.Location{}
	playerId, ok := state.Player(location)
	if !ok {
		return fmt.Errorf("no player found")
	}
	a.Direction.Flip().Move(&location.Position)
	state.EntityStore.Give(playerId, location)
	return nil
}
