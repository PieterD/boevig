package sim

import (
	"fmt"

	"github.com/PieterD/rlpoc/old/internal/internal/pkg/cardinal"
	"github.com/PieterD/rlpoc/old/internal/internal/spacetime"
	"github.com/PieterD/rlpoc/old/internal/internal/spacetime/command"
	event2 "github.com/PieterD/rlpoc/old/internal/internal/spacetime/event"
)

type Loom struct {
	fabric       spacetime.Fabric
	playerId     spacetime.EntityId
	nextEntityId spacetime.EntityId
}

var _ spacetime.Loom = &Loom{}

func NewLoom(f spacetime.Fabric) *Loom {
	return &Loom{
		fabric:       f,
		playerId:     1,
		nextEntityId: 3,
	}
}

func (l *Loom) Apply(e spacetime.Command) error {
	switch evt := e.(type) {
	case *command.Init:
		return l.processInit(evt)
	case *command.Move:
		return l.processMove(evt)
	default:
		return fmt.Errorf("event type %T unsupported by loom", e)
	}
}

const firstRoom = `
#######
#.....#
#.....#
#.....#
#######
`

func (l *Loom) processInit(initCommand *command.Init) error {
	initEvent := &event2.Init{}
	var playerId spacetime.EntityId = 1
	var roomId spacetime.EntityId = 2
	var nextEntityId spacetime.EntityId = 3
	createPlayerEvent := &event2.CreatePlayer{
		Id:   l.playerId,
		Name: "Player",
	}
	roomGeometryEvent := &event2.RoomGeometry{
		Id:    roomId,
		Shape: firstRoom,
	}
	spawnEvent := &event2.Spawn{
		Id: l.playerId,
		Location: event2.Location{
			Room: roomId,
			Coord: cardinal.Coord{
				X: 0,
				Y: 0,
			},
		},
	}
	if err := l.fabric.Apply(initEvent); err != nil {
		return fmt.Errorf("applying init event: %w", err)
	}
	if err := l.fabric.Apply(createPlayerEvent); err != nil {
		return fmt.Errorf("applying player creation event: %w", err)
	}
	if err := l.fabric.Apply(roomGeometryEvent); err != nil {
		return fmt.Errorf("applying room geometry event: %w", err)
	}
	if err := l.fabric.Apply(spawnEvent); err != nil {
		return fmt.Errorf("applying spawn event: %w", err)
	}
	l.playerId = playerId
	l.nextEntityId = nextEntityId
	return nil
}

func (l *Loom) processMove(moveCommand *command.Move) error {
	moveEvent := &event2.Move{
		Id:        l.playerId,
		Direction: moveCommand.Direction,
	}
	if err := l.fabric.Apply(moveEvent); err != nil {
		return fmt.Errorf("applying move event to fabric: %w", err)
	}
	return nil
}

func (l *Loom) Undo() error {
	panic("implement me")
}
