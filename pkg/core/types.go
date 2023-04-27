package core

type EntityID int64

type LocalCoordinate struct {
	X int
	Y int
}

func (lc LocalCoordinate) Negative() LocalCoordinate {
	return LocalCoordinate{
		X: -lc.X,
		Y: -lc.Y,
	}
}

type Direction int8

const (
	North Direction = iota
	East
	South
	West
)
