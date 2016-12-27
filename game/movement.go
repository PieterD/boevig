package game

import (
	"fmt"

	"github.com/PieterD/boevig/game/atlas/aspect"
)

type turnDirection int

const (
	turnLeft   turnDirection = 3
	turnRight  turnDirection = 1
	turnAround turnDirection = 2
)

func (g *Game) CmdMove() {
	cpos := g.pos
	npos := cpos
	switch g.dir {
	case 0:
		npos.Y--
	case 1:
		npos.X++
	case 2:
		npos.Y++
	case 3:
		npos.X--
	}
	if g.atlas.IsPassable(npos) {
		g.pos = npos
		g.atlas.Vision(g.pos)
	} else if g.atlas.GetFeature(npos) == aspect.ClosedDoor {
		fmt.Printf("you open the door.\n")
		g.atlas.SetFeature(npos, aspect.OpenDoor)
		g.atlas.Vision(g.pos)
	} else {
		fmt.Printf("thump!\n")
	}
}

func (g *Game) CmdTurn(td turnDirection) {
	g.dir = (g.dir + int(td)) % 4
}
