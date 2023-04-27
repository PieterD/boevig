package event

import (
	"github.com/PieterD/rlpoc/old/internal/internal/pkg/cardinal"
	"github.com/PieterD/rlpoc/old/internal/internal/spacetime"
)

type Init struct{}

type CreatePlayer struct {
	Id   spacetime.EntityId
	Name string
}

type RoomGeometry struct {
	Id    spacetime.EntityId
	Shape string
}

type Spawn struct {
	Id       spacetime.EntityId
	Location Location
}

type Move struct {
	Id        spacetime.EntityId
	Direction cardinal.Direction
}
