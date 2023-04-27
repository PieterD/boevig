package main

import (
	_ "image/png"

	"github.com/PieterD/boevig/game"
	"github.com/PieterD/boevig/grid"
)

func main() {
	grid.Run("resources/rogue_yun_16x16.png", 16, 16, game.New())
}
