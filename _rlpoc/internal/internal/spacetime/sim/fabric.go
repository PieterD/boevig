package sim

import (
	"github.com/PieterD/rlpoc/old/internal/internal/pkg/cardinal"
	"github.com/PieterD/rlpoc/old/internal/internal/spacetime"
	"github.com/PieterD/rlpoc/old/internal/internal/spacetime/event"
)

type Fabric struct {
	playerId   spacetime.EntityId
	playerName string
	//rooms      map[spacetime.EntityId]*fabricRoom
}

func (f *Fabric) Apply(e spacetime.Event) error {
	switch evt := e.(type) {
	case event.Init:
	case event.CreatePlayer:
		f.playerId = evt.Id
		f.playerName = evt.Name
	case event.RoomGeometry:
		panic("not implemented")
	case event.Spawn:
		panic("not implemented")
	case event.Move:
		panic("not implemented")
	}
	panic("not implemented")
}

func (f *Fabric) Undo(event spacetime.Event) error {
	panic("implement me")
}

func (f *Fabric) Draw(dst spacetime.Canvas, r cardinal.Box, sp cardinal.Coord) error {
	panic("implement me")
}

func NewFabric() *Fabric {
	return &Fabric{}
}

var _ spacetime.Fabric = &Fabric{}
