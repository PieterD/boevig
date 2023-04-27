package atlas

import (
	"image"
	"math/rand"

	"github.com/PieterD/boevig/game/atlas/aspect"
)

type caveSlice struct {
	w     int
	h     int
	cells []byte
	swp   []byte
}

func newCaveSlice(atlas *Atlas) *caveSlice {
	w := atlas.bounds.Max.X
	h := atlas.bounds.Max.Y
	cells := make([]byte, w*h*2)
	swp := cells[w*h:]
	cells = cells[:w*h]
	cs := caveSlice{
		w:     w,
		h:     h,
		cells: cells,
		swp:   swp,
	}
	return &cs
}

func (cs *caveSlice) set(i int, x, y int) {
	cs.cells[x+y*cs.w] = byte(i)
}

func (cs *caveSlice) get(x, y int) int {
	return int(cs.cells[x+y*cs.w])
}

func (cs *caveSlice) setSwp(i int, x, y int) {
	cs.swp[x+y*cs.w] = byte(i)
}
func (cs *caveSlice) getSwp(x, y int) int {
	return int(cs.swp[x+y*cs.w])
}

func (cs *caveSlice) fillBorders(size int) {
	for i := 0; i < size; i++ {
		for x := 0; x < cs.w; x++ {
			cs.set(1, x, i)
			cs.set(1, x, cs.h-1-i)
		}
		for y := 0; y < cs.h; y++ {
			cs.set(1, i, y)
			cs.set(1, cs.w-1-i, y)
		}
	}
}

func (cs *caveSlice) fillRandomCells(perc int) {
	for y := 0; y < cs.h; y++ {
		for x := 0; x < cs.w; x++ {
			if rand.Intn(100) < perc {
				cs.set(1, x, y)
			}
		}
	}
}

func (cs *caveSlice) automaton() {
	cs.swp, cs.cells = cs.cells, cs.swp
	for y := 1; y < cs.h-1; y++ {
		for x := 1; x < cs.w-1; x++ {
			neighbors := 0
			for yy := -1; yy <= 1; yy++ {
				for xx := -1; xx <= 1; xx++ {
					if cs.getSwp(x+xx, y+yy) == 1 {
						neighbors++
					}
				}
			}
			if neighbors >= 5 {
				cs.set(1, x, y)
			} else {
				cs.set(0, x, y)
			}
		}
	}
}

func (cs *caveSlice) closeDiagonalGaps() {
	for y := 1; y < cs.h-1; y++ {
		for x := 1; x < cs.w-1; x++ {
			b00 := cs.get(x-1, y-1) == 1
			b10 := cs.get(x, y-1) == 1
			b01 := cs.get(x-1, y) == 1
			b11 := cs.get(x, y) == 1
			if b00 && !b10 && !b01 && b11 {
				cs.set(1, x, y-1)
				cs.set(1, x-1, y)
			}
			if !b00 && b10 && b01 && !b11 {
				cs.set(1, x-1, y-1)
				cs.set(1, x, y)
			}
		}
	}
}

func (cs *caveSlice) clear() {
	for i := range cs.cells {
		cs.cells[i] = 0
	}
}

func (cs *caveSlice) clearSwp() {
	for i := range cs.swp {
		cs.swp[i] = 0
	}
}

func (cs *caveSlice) randomFloor() (x, y int) {
	for {
		x = rand.Intn(cs.w)
		y = rand.Intn(cs.h)
		if cs.get(x, y) == 0 {
			break
		}
	}
	return
}

func (cs *caveSlice) singleCave() bool {
	if !cs.fillLargeCave() {
		return false
	}
	for y := 0; y < cs.h; y++ {
		for x := 0; x < cs.w; x++ {
			if cs.getSwp(x, y) == 0 && cs.get(x, y) == 0 {
				cs.set(1, x, y)
			}
		}
	}
	return true
}

func (cs *caveSlice) fillLargeCave() bool {
	try := 0
	for try = 0; try < 50; try++ {
		cs.clearSwp()
		sx, sy := cs.randomFloor()
		num := cs.fill(sx, sy, func(x, y int) bool {
			if x < 0 || y < 0 || x >= cs.w || y >= cs.h {
				// Out of bounds
				return false
			}
			if cs.getSwp(x, y) == 1 {
				// Cell already visited
				return false
			}
			if cs.get(x, y) == 1 {
				// Cell is a wall
				return false
			}
			// Set cell as visited
			cs.setSwp(1, x, y)
			return true
		})
		perc := (num * 100) / (cs.w * cs.h)
		if perc >= 55 {
			return true
		}
	}
	return false
}

func (cs *caveSlice) fill(xStart, y int, f func(x, y int) bool) int {
	if !f(xStart, y) {
		return 0
	}
	num := 1
	xMin := xStart
	xMax := xStart
	for x := xStart + 1; f(x, y); x++ {
		xMax = x
		num++
	}
	for x := xStart - 1; f(x, y); x-- {
		xMin = x
		num++
	}
	for x := xMin; x <= xMax; x++ {
		num += cs.fill(x, y-1, f)
		num += cs.fill(x, y+1, f)
	}
	return num
}

func (cs *caveSlice) setAtlas(atlas *Atlas) {
	for y := 1; y < cs.h-1; y++ {
		for x := 1; x < cs.w-1; x++ {
			cell := atlas.cell(image.Point{X: x, Y: y})
			if cs.get(x, y) == 1 {
				cell.feature = aspect.Wall
			}
		}
	}
}

func GenCave(atlas *Atlas) {
	cs := newCaveSlice(atlas)
	for try := 0; try < 50; try++ {
		cs.fillBorders(3)
		cs.fillRandomCells(45)

		// Run the automaton 5 times
		for cycle := 0; cycle < 7; cycle++ {
			cs.automaton()
		}
		cs.closeDiagonalGaps()
		if cs.singleCave() {
			cs.setAtlas(atlas)
			return
		} else {
			cs.clear()
			cs.clearSwp()
		}
	}
	panic("Could not generate a suitable cave after 50 tries!")
}
