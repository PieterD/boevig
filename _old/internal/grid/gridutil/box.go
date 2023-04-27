package gridutil

import (
	"image"

	"github.com/PieterD/boevig/grid"
)

func SingleBox(d grid.DrawableGrid, bounds image.Rectangle, fore, back grid.Color) {
	drawBox(d, bounds, fore, back, []int{196, 179, 218, 191, 217, 192})
}

func DoubleBox(d grid.DrawableGrid, bounds image.Rectangle, fore, back grid.Color) {
	drawBox(d, bounds, fore, back, []int{205, 186, 201, 187, 188, 200})
}

func drawBox(d grid.DrawableGrid, bounds image.Rectangle, fore, back grid.Color, parts []int) {
	for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
		d.Set(x, bounds.Min.Y, parts[0], fore, back)
		d.Set(x, bounds.Max.Y-1, parts[0], fore, back)
	}
	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		d.Set(0, y, parts[1], fore, back)
		d.Set(bounds.Max.X-1, y, parts[1], fore, back)
	}
	d.Set(bounds.Min.X, bounds.Min.Y, parts[2], fore, back)
	d.Set(bounds.Max.X-1, bounds.Min.Y, parts[3], fore, back)
	d.Set(bounds.Max.X-1, bounds.Max.Y-1, parts[4], fore, back)
	d.Set(bounds.Min.X, bounds.Max.Y-1, parts[5], fore, back)
}
