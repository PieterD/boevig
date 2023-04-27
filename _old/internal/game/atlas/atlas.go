package atlas

import (
	"image"

	"math/rand"
	"time"

	"github.com/PieterD/boevig/game/atlas/aspect"
	"github.com/PieterD/boevig/grid"
	"github.com/PieterD/boevig/vision"
	"github.com/PieterD/boevig/wallify"
)

type Glyph struct {
	Code int
	Fore grid.Color
	Back grid.Color
}

func Translate(screen image.Rectangle, center image.Point, atlas image.Rectangle) image.Point {
	tl := center.Sub(screen.Max.Div(2))
	if screen.Max.X > atlas.Max.X {
		tl.X = -(screen.Max.X - atlas.Max.X) / 2
	} else {
		if tl.X < 0 {
			tl.X = 0
		}
		if tl.X >= atlas.Max.X-screen.Max.X {
			tl.X = atlas.Max.X - screen.Max.X
		}
	}
	if screen.Max.Y > atlas.Max.Y {
		tl.Y = -(screen.Max.Y - atlas.Max.Y) / 2
	} else {
		if tl.Y < 0 {
			tl.Y = 0
		}
		if tl.Y >= atlas.Max.Y-screen.Max.Y {
			tl.Y = atlas.Max.Y - screen.Max.Y
		}
	}
	return tl
}

type Atlas struct {
	cells   []Cell
	bounds  image.Rectangle
	visible uint64
}

func (atlas *Atlas) cell(p image.Point) *Cell {
	if !p.In(atlas.bounds) {
		return &Cell{}
	}
	return &atlas.cells[p.X+p.Y*atlas.bounds.Max.X]
}

func New() *Atlas {
	rand.Seed(time.Now().UnixNano())
	w := 100
	h := 60
	atlas := &Atlas{
		cells: make([]Cell, w*h),
		bounds: image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: w, Y: h},
		},
		visible: 1,
	}
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			atlas.setFeature(x, y, aspect.Floor)
		}
	}
	for x := 0; x < w; x++ {
		atlas.setFeature(x, 0, aspect.Wall)
		atlas.setFeature(x, h-1, aspect.Wall)
	}
	for y := 0; y < h; y++ {
		atlas.setFeature(0, y, aspect.Wall)
		atlas.setFeature(w-1, y, aspect.Wall)
	}
	//GenTestlevel(atlas)
	GenCave(atlas)
	return atlas
}

func (atlas *Atlas) ExploreAll() {
	for y := 0; y < atlas.bounds.Max.Y; y++ {
		for x := 0; x < atlas.bounds.Max.X; x++ {
			atlas.cell(image.Point{X: x, Y: y}).seen = true
		}
	}
}

func (atlas *Atlas) Bounds() image.Rectangle {
	return atlas.bounds
}

func (atlas *Atlas) GetFeature(pos image.Point) aspect.Feature {
	return atlas.cell(pos).feature
}

func (atlas *Atlas) SetFeature(pos image.Point, feature aspect.Feature) {
	atlas.cell(pos).feature = feature
}

func (atlas *Atlas) setFeature(x, y int, ft aspect.Feature) {
	atlas.SetFeature(image.Point{X: x, Y: y}, ft)
}

func (atlas *Atlas) Glyph(p image.Point) Glyph {
	cell := atlas.cell(p)
	glyph := Glyph{
		Code: 32,
		Fore: grid.Black,
		Back: grid.Black,
	}
	if atlas.cell(p).seen {
		switch cell.feature {
		case aspect.Wall:
			glyph = Glyph{
				Code: wallify.Wallify(atlas, p, wallify.SingleWall),
				Fore: grid.Gray,
				Back: grid.Black,
			}
		case aspect.Floor:
			glyph = Glyph{
				Code: atlas.floorrune(p, floorRune),
				Fore: grid.DarkGray,
				Back: grid.Black,
			}
		case aspect.ClosedDoor:
			glyph = Glyph{
				Code: 43,
				Fore: grid.DarkRed,
				Back: grid.Black,
			}
		case aspect.OpenDoor:
			glyph = Glyph{
				Code: 47,
				Fore: grid.DarkRed,
				Back: grid.Black,
			}
		}
		if !atlas.IsVisible(p) {
			if glyph.Fore != grid.Black {
				glyph.Fore = grid.VeryDarkGray
			}
			glyph.Back = grid.Black
		}
	}
	return glyph
}

func (atlas *Atlas) IsPassable(p image.Point) bool {
	return atlas.cell(p).feature.Passable
}

func (atlas *Atlas) IsTransparent(p image.Point) bool {
	return atlas.cell(p).feature.Transparent
}

func (atlas *Atlas) SetVisible(p image.Point) {
	atlas.cell(p).visible = atlas.visible
	atlas.cell(p).seen = true
}

func (atlas *Atlas) IsVisible(p image.Point) bool {
	return atlas.cell(p).visible == atlas.visible
}

func (atlas *Atlas) IsWallable(p image.Point) bool {
	return atlas.cell(p).feature.Wallable
}

func (atlas *Atlas) IsSeen(p image.Point) bool {
	return atlas.cell(p).seen
}

func (atlas *Atlas) RandomFloor() image.Point {
	for try := 0; try < 1000; try++ {
		x := rand.Intn(atlas.bounds.Max.X)
		y := rand.Intn(atlas.bounds.Max.Y)
		p := image.Point{X: x, Y: y}
		if atlas.cell(p).feature == aspect.Floor {
			return p
		}
	}
	panic("Couldn't find random floor after 1000 tries!")
}

func (atlas *Atlas) Vision(source image.Point) {
	atlas.visible++
	vision.ShadowCastPar(atlas, vision.EndlessRadius(), source)
}

//var floorRune = []int{44, 46, 96, 249, 250}
//var floorRune = []int{44, 46, 96, 249, 39}
//var floorRune = []int{44, 46, 96, 249, 39, 250, 250}
var floorRune = []int{250, 44, 250, 46, 250, 96, 250, 249, 250, 39, 250}

//var floorRune = []int{250}

func (atlas *Atlas) floorrune(p image.Point, runes []int) int {
	x := uint64(p.X)
	y := uint64(p.Y)
	ui := ((x<<32)|y)*(x^y) + y - x
	return runes[ui%uint64(len(runes))]
}
