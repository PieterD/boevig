package game

import (
	"fmt"
	"image"

	"github.com/PieterD/boevig/game/atlas"
	"github.com/PieterD/boevig/grid"
)

type Game struct {
	die   bool
	pos   image.Point
	dir   int
	atlas *atlas.Atlas
}

func New() *Game {
	atlas := atlas.New()
	pos := atlas.RandomFloor()
	g := &Game{
		pos:   pos,
		dir:   2,
		atlas: atlas,
	}

	g.atlas.Vision(g.pos)
	return g
}

func (g *Game) Draw(d grid.DrawableGrid) {
	screenBounds := d.GridSize()
	atlasBounds := g.atlas.Bounds()
	translator := atlas.Translate(screenBounds, g.pos, atlasBounds)
	for x := 0; x < screenBounds.Max.X; x++ {
		for y := 0; y < screenBounds.Max.Y; y++ {
			screenCoord := image.Point{X: x, Y: y}
			atlasCoord := screenCoord.Add(translator)
			if atlasCoord.In(atlasBounds) {
				glyph := g.atlas.Glyph(atlasCoord)
				d.Set(screenCoord.X, screenCoord.Y, glyph.Code, glyph.Fore, glyph.Back)
			}
		}
	}
	//gridutil.SingleBox(d, d.GridSize(), grid.Black, grid.White)
	//gridutil.Text(d, image.Point{X: 2, Y: 2}, "Hello moo", grid.Red, grid.Black)
	var arrow int
	switch g.dir {
	case 0:
		arrow = 24
	case 1:
		arrow = 26
	case 2:
		arrow = 25
	case 3:
		arrow = 27
	}
	//d.Set(g.x, g.y, arrow, grid.Green, grid.Black)
	trpos := g.pos.Sub(translator)
	d.Set(trpos.X, trpos.Y, arrow, grid.BrightGreen, grid.Black)
}

func (g *Game) Char(r rune) {
	fmt.Printf("char %d(%c)\n", r, r)
}

func (g *Game) Key(e grid.KeyEvent) {
	fmt.Printf("key %#v\n", e)
	switch e.Key {
	case grid.KeyEscape:
		g.die = true
	case grid.KeyUp:
		g.CmdMove()
	case grid.KeyLeft:
		g.CmdTurn(turnLeft)
	case grid.KeyRight:
		g.CmdTurn(turnRight)
	case grid.KeyDown:
		g.CmdTurn(turnAround)
	case grid.KeySpace:
		g.atlas.ExploreAll()
	}
}

func (g *Game) MouseMove(e grid.MouseMoveEvent) {
	fmt.Printf("mousemove  %#v\n", e)
}

func (g *Game) MouseClick(e grid.MouseClickEvent) {
	fmt.Printf("mouseclick %#v\n", e)
}

func (g *Game) MouseDrag(e grid.MouseDragEvent) {
	fmt.Printf("mousedrag  %#v\n", e)
}

func (g *Game) Fin(last bool) bool {
	//fmt.Printf("finish %t\n", last)
	return g.die
}
