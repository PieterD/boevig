package game

type Player struct{}

type GridLocation struct {
	X int
	Y int
}

type Terrain struct {
	Passable bool
}

type Direction int

const (
	N  Direction = 1
	E  Direction = 2
	S  Direction = 4
	W  Direction = 8
	NE           = N | E
	SE           = S | E
	SW           = S | W
	NW           = N | W
)
