package spacetime

import (
	cardinal2 "github.com/PieterD/rlpoc/old/internal/internal/pkg/cardinal"
)

type EntityId uint64

type Command interface{}

type Event interface{}

// Journal manages the full history of all user inputs, including Undo.
// Each Command is stored, and non-journal-only events are converted and passed on to the Loom as regular Events.
type Journal interface {
	Apply(Command) error
	RedoOptions() ([]Command, error)
}

// Loom consumes Commands, and runs the game world until another input is required.
// It will create Events, and pass them on to the Fabric.
type Loom interface {
	Apply(Command) error
	Undo() error
}

// Fabric represents the state of the world following on a particular history of Events.
// Fabric itself stores no Events, and creates none.
type Fabric interface {
	Apply(Event) error
	Undo(Event) error
	// Draw will align r.TL in dst with sp in the Fabric, and then fill the rectangle with content.
	Draw(dst Canvas, r cardinal2.Box, sp cardinal2.Coord) error
}

type Canvas interface {
	Clear() error
	Set(cardinal2.Coord, rune) error
}
