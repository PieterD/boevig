package cardinal

type Box struct {
	TL Coord
	BR Coord
}

// Contains returns true if the Coord lies within Box.
func (b Box) Contains(at Coord) bool {
	if b.TL.X > at.X {
		return false
	}
	if b.TL.Y > at.Y {
		return false
	}
	if b.BR.X < at.X {
		return false
	}
	if b.BR.Y < at.Y {
		return false
	}
	return true
}

// Include returns a new Box, which has (maybe) been enlarged to include a new Coord.
func (b Box) Include(at Coord) Box {
	if b.TL.X > at.X {
		b.TL.X = at.X
	}
	if b.TL.Y > at.Y {
		b.TL.Y = at.Y
	}
	if b.BR.X < at.X {
		b.BR.X = at.X
	}
	if b.BR.Y < at.Y {
		b.BR.Y = at.Y
	}
	return b
}

// Visit will visit every cell, row by row, starting at the top left.
// If f returns an error, Visit will abort and return it immediately.
// newRow will be true only for the leftmost column of Coords.
// start will be true only for the top-leftmost Coord.
func (b Box) Visit(f func(coord Coord, start bool, newRow bool) error) error {
	start := true
	for y := b.TL.Y; y <= b.BR.Y; y++ {
		newRow := true
		for x := b.TL.X; x <= b.BR.X; x++ {
			if err := f(Coord{x, y}, start, newRow); err != nil {
				return err
			}
			newRow = false
			start = false
		}
	}
	return nil
}
