package vision

import "image"

type visionTransform func(image.Point) image.Point

type visionTransformer struct {
	source image.Point
	mul    image.Point
	swapxy bool
}

var quads = []image.Point{
	{X: 1, Y: 1},
	{X: 1, Y: -1},
	{X: -1, Y: 1},
	{X: -1, Y: -1},
}

func (vt visionTransformer) quad(i int) visionTransformer {
	vt.mul = quads[i]
	return vt
}

func (vt visionTransformer) swap() visionTransformer {
	vt.swapxy = true
	return vt
}

func (vt visionTransformer) trans(p image.Point) image.Point {
	x := p.X * vt.mul.X
	y := p.Y * vt.mul.Y
	if vt.swapxy {
		x, y = y, x
	}
	return image.Point{X: x + vt.source.X, Y: y + vt.source.Y}
}

type fraction struct {
	num, den int64
}

func newFraction(num, den int) fraction {
	return fraction{
		num: int64(num),
		den: int64(den),
	}
}

func one() fraction {
	return newFraction(1, 1)
}

func zero() fraction {
	return newFraction(0, 1)
}

func (a fraction) less(b fraction) bool {
	return a.num*b.den < b.num*a.den
}

func (a fraction) mul(s int) fraction {
	a.num *= int64(s)
	return a
}

func (a fraction) whole() int {
	return int(a.num / a.den)
}
