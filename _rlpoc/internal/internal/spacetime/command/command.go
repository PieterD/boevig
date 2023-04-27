package command

import (
	"github.com/PieterD/rlpoc/old/internal/internal/pkg/cardinal"
)

type Init struct{}

type Undo struct{}

type Redo struct{}

type Move struct {
	Direction cardinal.Direction
}
