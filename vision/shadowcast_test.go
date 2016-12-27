package vision

import (
	"image"
	"testing"
)

type testShadowCastMap struct {
	visible uint64
	cells   []testCell
	bounds  image.Rectangle
}

type testCell struct {
	transparent bool
	visible     uint64
}

func (m *testShadowCastMap) cell(p image.Point) *testCell {
	if !p.In(m.bounds) {
		return &testCell{}
	}
	return &m.cells[p.X+p.Y*m.bounds.Max.X]
}

func newTestShadowCastMap(w, h int) *testShadowCastMap {
	m := &testShadowCastMap{
		visible: 1,
		cells:   make([]testCell, w*h),
		bounds:  image.Rectangle{Max: image.Point{X: w, Y: h}},
	}
	for x := 1; x < w-1; x++ {
		for y := 1; y < h-1; y++ {
			m.cell(image.Point{X: x, Y: y}).transparent = true
		}
	}
	return m
}

func (m *testShadowCastMap) SetVisible(p image.Point) {
	m.cell(p).visible = m.visible
}

func (m *testShadowCastMap) IsTransparent(p image.Point) bool {
	return m.cell(p).transparent
}

type shadowCastTest struct {
	name   string
	m      *testShadowCastMap
	r      Radius
	source image.Point
	f      func(m ShadowCastMap, r Radius, source image.Point)
}

func BenchmarkTinyRoom(b *testing.B) {
	m := newTestShadowCastMap(100, 100)
	r := EndlessRadius()
	for i := 49; i <= 51; i++ {
		m.cell(image.Point{X: i, Y: 49}).transparent = false
		m.cell(image.Point{X: i, Y: 51}).transparent = false
		m.cell(image.Point{X: 49, Y: i}).transparent = false
		m.cell(image.Point{X: 51, Y: i}).transparent = false
	}
	source := image.Point{X: 50, Y: 50}
	tests := []shadowCastTest{
		{"ShadowCastFloat", m, r, source, ShadowCastFloat},
		{"ShadowCast", m, r, source, ShadowCast},
		{"ShadowCastPar", m, r, source, ShadowCastPar},
	}
	run(b, tests)
}

func run(bb *testing.B, tests []shadowCastTest) {
	for _, test := range tests {
		bb.Run(test.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				test.m.visible++
				test.f(test.m, test.r, test.source)
			}
		})
	}
}

func BenchmarkLargeEmptyCenter(b *testing.B) {
	m := newTestShadowCastMap(100, 100)
	r := EndlessRadius()
	source := image.Point{X: 50, Y: 50}
	tests := []shadowCastTest{
		{"ShadowCastFloat", m, r, source, ShadowCastFloat},
		{"ShadowCast", m, r, source, ShadowCast},
		{"ShadowCastPar", m, r, source, ShadowCastPar},
	}
	run(b, tests)
}

func BenchmarkLargeEmptyLeft(b *testing.B) {
	m := newTestShadowCastMap(100, 100)
	r := EndlessRadius()
	source := image.Point{X: 1, Y: 50}
	tests := []shadowCastTest{
		{"ShadowCastFloat", m, r, source, ShadowCastFloat},
		{"ShadowCast", m, r, source, ShadowCast},
		{"ShadowCastPar", m, r, source, ShadowCastPar},
	}
	run(b, tests)
}

func BenchmarkLargeScattered10(b *testing.B) {
	m := newTestShadowCastMap(100, 100)
	r := EndlessRadius()
	for x := 0; x < m.bounds.Max.X; x += 10 {
		for y := 0; y < m.bounds.Max.Y; y += 10 {
			m.cell(image.Point{X: x, Y: y}).transparent = false
		}
	}
	source := image.Point{X: 50, Y: 50}
	tests := []shadowCastTest{
		{"ShadowCastFloat", m, r, source, ShadowCastFloat},
		{"ShadowCast", m, r, source, ShadowCast},
		{"ShadowCastPar", m, r, source, ShadowCastPar},
	}
	run(b, tests)
}

func BenchmarkLargeScattered5(b *testing.B) {
	m := newTestShadowCastMap(100, 100)
	r := EndlessRadius()
	for x := 0; x < m.bounds.Max.X; x += 5 {
		for y := 0; y < m.bounds.Max.Y; y += 5 {
			m.cell(image.Point{X: x, Y: y}).transparent = false
		}
	}
	source := image.Point{X: 50, Y: 50}
	tests := []shadowCastTest{
		{"ShadowCastFloat", m, r, source, ShadowCastFloat},
		{"ShadowCast", m, r, source, ShadowCast},
		{"ShadowCastPar", m, r, source, ShadowCastPar},
	}
	run(b, tests)
}
