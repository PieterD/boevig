package journal

import (
	"fmt"

	"github.com/PieterD/rlpoc/old/internal/internal/spacetime"
	"github.com/PieterD/rlpoc/old/internal/internal/spacetime/command"
)

type Journal struct {
	commands            []spacetime.Command
	currentCommandIndex int
	loom                spacetime.Loom
}

func NewJournal(loom spacetime.Loom) *Journal {
	return &Journal{
		commands: []spacetime.Command{
			command.Init{},
		},
		currentCommandIndex: 0,
		loom:                loom,
	}
}

func (j *Journal) Apply(c spacetime.Command) error {
	j.commands = append(j.commands, c)
	switch cmd := c.(type) {
	case *command.Init:
		if err := j.processInit(cmd); err != nil {
			return fmt.Errorf("processing init command in journal: %w", err)
		}
	case *command.Undo:
		return j.processUndo(cmd)
	case *command.Redo:
		return j.processRedo(cmd)
	}
	if err := j.loom.Apply(c); err != nil {
		return fmt.Errorf("applying event to loom: %w", err)
	}
	return nil
}

func (j *Journal) processInit(cmd *command.Init) error {
	return nil
}

func (j *Journal) processUndo(cmd *command.Undo) error {
	panic("not implemented")
}

func (j *Journal) processRedo(cmd *command.Redo) error {
	panic("not implemented")
}

func (j *Journal) RedoOptions() ([]spacetime.Command, error) {
	panic("implement me")
}

var _ spacetime.Journal = &Journal{}
