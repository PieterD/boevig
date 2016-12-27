package wallify

import (
	"image"
)

type Wallifier interface {
	IsWallable(p image.Point) bool
	IsSeen(p image.Point) bool
}

/*
Wallify determins the type of wall rune to use for a given cell.
A bitmask of its neighbors is made:
           1
          8#2
           4
The bit is set if the corresponding neighbor IsWallable.
For the given bitmask, the corresponding rune index is found in the wallRune slice.
This rune index is then found in the provided runes slice, which can be,
say, SingleWall or DoubleWall, and these contain the rune code to be returned.

Normally, if the entire level is discovered (IsSeen), this would be enough.
However, if this is all we do, T-shapes and Cross-shapes cause problems
if one of the blocks adjacent to a T or Cross has not yet been discovered. Consider:
 │
 ├─
 │
If the lower straight wall has not yet been discovered, then the fact that the
corner is a T-shape rather than a corner-shape gives away the fact that there is a
wall on the other side.

To fix this, if the first step results in a cross or a T, we create a visibility mask
of whether the cell's neighbors are visible, and AND it with the bitmask.

This works for most situations, except when not a single neighbor has been seen;
this would result in a pillar, which would turn into a wall when neighbors are
later discovered. This looks terrible.

So, in order to fix this, if the retuls of step 2 is a pillar, we try to fit one
of the invisiWall elements to the mask, which contains only straight walls and corners.

Only if we can't fit any of these, a pillar is probably right.
*/

//var singleWall = []int{79, 179, 196, 192, 218, 191, 217, 195, 194, 180, 193, 197}
//var singleWall = []int{9, 179, 196, 192, 218, 191, 217, 195, 194, 180, 193, 197}
var SingleWall = []int{233, 179, 196, 192, 218, 191, 217, 195, 194, 180, 193, 197}
var DoubleWall = []int{233, 186, 205, 200, 201, 187, 188, 204, 203, 185, 202, 206}
var wallRune = []int{0, 1, 2, 3, 1, 1, 4, 7, 2, 6, 2, 10, 5, 9, 8, 11}
var invisiWall = []int{5, 10, 3, 6, 12, 9, 1, 2, 4, 8}

func Wallify(w Wallifier, p image.Point, runes []int) int {
	x := image.Point{X: 1}
	y := image.Point{Y: 1}
	bits := 0
	if w.IsWallable(p.Sub(y)) {
		bits |= 1
	}
	if w.IsWallable(p.Add(x)) {
		bits |= 2
	}
	if w.IsWallable(p.Add(y)) {
		bits |= 4
	}
	if w.IsWallable(p.Sub(x)) {
		bits |= 8
	}
	switch bits {
	case 7, 14, 13, 11, 15:
		// Cross or T-shape; check visibility
		visibits := 0
		if w.IsWallable(p.Sub(y)) && w.IsSeen(p.Sub(y)) {
			visibits |= 1
		}
		if w.IsWallable(p.Add(x)) && w.IsSeen(p.Add(x)) {
			visibits |= 2
		}
		if w.IsWallable(p.Add(y)) && w.IsSeen(p.Add(y)) {
			visibits |= 4
		}
		if w.IsWallable(p.Sub(x)) && w.IsSeen(p.Sub(x)) {
			visibits |= 8
		}
		if visibits&bits == 0 {
			// No neighbors are visible, which would result in a pillar.
			// Instead, try to fit one of the non-T, non-cross shapes.
			for _, iw := range invisiWall {
				if bits&visibits&iw == iw {
					bits = iw
					break
				}
			}
		} else {
			bits = bits & visibits
		}
	}
	return runes[wallRune[bits]]
}
