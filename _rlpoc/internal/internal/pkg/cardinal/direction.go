package cardinal

import "fmt"

type (
	FaceCoord struct {
		Face  Direction
		Coord Coord
	}
	Coord struct {
		X, Y int32
	}
	Direction byte
	Rotation  int
)

const (
	ClockWise   Rotation = 1
	CounterWise Rotation = -1
)

const (
	North Direction = 1 << iota
	East
	South
	West
)

func NewFaceCoord(face Direction, coord Coord) FaceCoord {
	return FaceCoord{
		Face:  face,
		Coord: coord,
	}
}

func (c Coord) Less(than Coord) bool {
	switch {
	case c.Y < than.Y:
		return true
	case c.Y == than.Y:
		return c.X < than.X
	}
	return false
}

func (r Rotation) Normalise() Rotation {
	return r % 4
}

func (direction Direction) String() string {
	switch direction {
	case North:
		return "north"
	case North | East:
		return "northeast"
	case East:
		return "east"
	case South | East:
		return "southeast"
	case South:
		return "south"
	case South | West:
		return "southwest"
	default:
		return fmt.Sprintf("unknown_direction_%d", direction)
	}
}

func (direction Direction) Valid() bool {
	switch direction {
	case North, North | East, East, South | East, South, South | West, West, North | West:
		return true
	default:
		return false
	}
}

func (direction Direction) Split() []Direction {
	var split []Direction
	for _, dir := range []Direction{North, East, South, West} {
		if direction.Has(dir) {
			split = append(split, dir)
		}
	}
	return split
}

func (direction Direction) Has(has Direction) bool {
	return direction&has == has
}

func (direction Direction) Flip() Direction {
	var nd Direction
	if direction.Has(North) {
		nd |= South
	}
	if direction.Has(East) {
		nd |= West
	}
	if direction.Has(South) {
		nd |= North
	}
	if direction.Has(West) {
		nd |= East
	}
	return nd
}

func (direction Direction) Move(coord *Coord) {
	for _, dir := range direction.Split() {
		switch dir {
		case North:
			coord.Y--
		case East:
			coord.X++
		case South:
			coord.Y++
		case West:
			coord.X--
		}
	}
}

func (direction Direction) Rotate(r Rotation) Direction {
	if r < 0 {
		for ; r < 0; r++ {
			direction = direction.RotateRight()
		}
	} else if r > 0 {
		for ; r > 0; r-- {
			direction = direction.RotateLeft()
		}
	}
	return direction
}

func (direction Direction) RotateLeft() Direction {
	switch direction {
	case North:
		return West
	case North | East:
		return North | West
	case East:
		return North
	case South | East:
		return North | East
	case South:
		return East
	case South | West:
		return South | East
	case West:
		return South
	case North | West:
		return South | West
	default:
		panic(fmt.Errorf("unknown direction %d", direction))
	}
}

func (direction Direction) RotateRight() Direction {
	switch direction {
	case North:
		return East
	case North | East:
		return South | East
	case East:
		return South
	case South | East:
		return South | West
	case South:
		return West
	case South | West:
		return North | West
	case West:
		return North
	case North | West:
		return North | East
	default:
		panic(fmt.Errorf("unknown direction %d", direction))
	}
}
