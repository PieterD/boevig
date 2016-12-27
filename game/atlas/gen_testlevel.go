package atlas

import "github.com/PieterD/boevig/game/atlas/aspect"

func GenTestlevel(atlas *Atlas) {
	w := atlas.bounds.Max.X
	h := atlas.bounds.Max.Y
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
	for x := 0; x < w; x += 10 {
		for y := 0; y < h; y += 10 {
			atlas.setFeature(x, y, aspect.Wall)
		}
	}

	wallTest(atlas)

	max := 5
	for i := 1; i <= max; i++ {
		atlas.setFeature(max, i, aspect.Wall)
		atlas.setFeature(i, max, aspect.Wall)
	}
	atlas.setFeature(max, 2, aspect.ClosedDoor)

	for x := 0; x <= 10; x++ {
		atlas.setFeature(45+x, 5, aspect.Wall)
		atlas.setFeature(45+x, 15, aspect.Wall)
		atlas.setFeature(45, 5+x, aspect.Wall)
		atlas.setFeature(55, 5+x, aspect.Wall)
	}
	atlas.setFeature(50, 5, aspect.ClosedDoor)
	atlas.setFeature(50, 10, aspect.Floor)
}

func wallTest(atlas *Atlas) {
	// specials
	atlas.setFeature(10, 10, aspect.Wall)

	atlas.setFeature(10, 20, aspect.Wall)
	atlas.setFeature(10, 19, aspect.Wall)
	atlas.setFeature(10, 21, aspect.Wall)

	atlas.setFeature(10, 30, aspect.Wall)
	atlas.setFeature(9, 30, aspect.Wall)
	atlas.setFeature(11, 30, aspect.Wall)

	atlas.setFeature(10, 40, aspect.Wall)
	atlas.setFeature(9, 40, aspect.Wall)
	atlas.setFeature(11, 40, aspect.Wall)
	atlas.setFeature(10, 39, aspect.Wall)
	atlas.setFeature(10, 41, aspect.Wall)

	// doubles
	atlas.setFeature(20, 10, aspect.Wall)
	atlas.setFeature(20, 9, aspect.Wall)

	atlas.setFeature(20, 20, aspect.Wall)
	atlas.setFeature(21, 20, aspect.Wall)

	atlas.setFeature(20, 30, aspect.Wall)
	atlas.setFeature(20, 31, aspect.Wall)

	atlas.setFeature(20, 40, aspect.Wall)
	atlas.setFeature(19, 40, aspect.Wall)

	// triples
	atlas.setFeature(30, 10, aspect.Wall)
	atlas.setFeature(30, 9, aspect.Wall)
	atlas.setFeature(31, 10, aspect.Wall)

	atlas.setFeature(30, 20, aspect.Wall)
	atlas.setFeature(31, 20, aspect.Wall)
	atlas.setFeature(30, 21, aspect.Wall)

	atlas.setFeature(30, 30, aspect.Wall)
	atlas.setFeature(29, 30, aspect.Wall)
	atlas.setFeature(30, 31, aspect.Wall)

	atlas.setFeature(30, 40, aspect.Wall)
	atlas.setFeature(30, 39, aspect.Wall)
	atlas.setFeature(29, 40, aspect.Wall)

	// quads
	atlas.setFeature(40, 10, aspect.Wall)
	atlas.setFeature(40, 9, aspect.Wall)
	atlas.setFeature(40, 11, aspect.Wall)
	atlas.setFeature(41, 10, aspect.Wall)

	atlas.setFeature(40, 20, aspect.Wall)
	atlas.setFeature(41, 20, aspect.Wall)
	atlas.setFeature(40, 21, aspect.Wall)
	atlas.setFeature(39, 20, aspect.Wall)

	atlas.setFeature(40, 30, aspect.Wall)
	atlas.setFeature(40, 29, aspect.Wall)
	atlas.setFeature(39, 30, aspect.Wall)
	atlas.setFeature(40, 31, aspect.Wall)

	atlas.setFeature(40, 40, aspect.Wall)
	atlas.setFeature(40, 39, aspect.Wall)
	atlas.setFeature(41, 40, aspect.Wall)
	atlas.setFeature(39, 40, aspect.Wall)
}
