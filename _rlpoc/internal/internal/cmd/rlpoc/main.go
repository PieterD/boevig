package main

import (
	"fmt"
	"os"

	"github.com/PieterD/rlpoc/old/internal/internal/spacetime/command"
	"github.com/PieterD/rlpoc/old/internal/internal/spacetime/journal"
	sim2 "github.com/PieterD/rlpoc/old/internal/internal/spacetime/sim"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	f := sim2.NewFabric()
	l := sim2.NewLoom(f)
	j := journal.NewJournal(l)
	if err := j.Apply(command.Init{}); err != nil {
		return fmt.Errorf("applying Init command: %w", err)
	}
	return nil
}
