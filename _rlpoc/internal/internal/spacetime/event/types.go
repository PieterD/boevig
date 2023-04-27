package event

import (
	"github.com/PieterD/rlpoc/old/internal/internal/pkg/cardinal"
	"github.com/PieterD/rlpoc/old/internal/internal/spacetime"
)

type Location struct {
	Room  spacetime.EntityId
	Coord cardinal.Coord
}
