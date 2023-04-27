package gridutil

import (
	"fmt"
	"image"
	"unicode/utf8"

	"github.com/PieterD/boevig/grid"
)

func Text(d grid.DrawableGrid, start image.Point, text string, fore, back grid.Color) {
	x := start.X
	y := start.Y
	for {
		r, l := utf8.DecodeRuneInString(text)
		if l == 0 {
			return
		}
		if l == 1 && r == utf8.RuneError {
			panic(fmt.Errorf("Invalid utf-8"))
		}
		if l > 1 {
			panic(fmt.Errorf("Non-ascii character: '%c (%d)'", r, r))
		}
		d.Set(x, y, int(r), fore, back)
		x++
		text = text[l:]
	}
}
