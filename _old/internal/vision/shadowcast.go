package vision

import (
	"image"
	"runtime"
)

type ShadowCastMap interface {
	SetVisible(p image.Point)
	IsTransparent(p image.Point) bool
}

type callShadowCastMap struct {
	m  ShadowCastMap
	r  Radius
	vt visionTransformer
}

var sendchan chan callShadowCastMap
var recvchan chan bool

func init() {
	routines := runtime.GOMAXPROCS(-1)
	sendchan = make(chan callShadowCastMap, 4)
	recvchan = make(chan bool, 4)
	for i := 0; i < routines; i++ {
		go func() {
			for call := range sendchan {
				shadowCastOctant(1, one(), zero(), call.m, call.r, call.vt)
				recvchan <- true
			}
		}()
	}
}

func ShadowCastPar(m ShadowCastMap, r Radius, source image.Point) {
	vt := visionTransformer{source: source}
	for i := range quads {
		sendchan <- callShadowCastMap{
			m:  m,
			r:  r,
			vt: vt.quad(i),
		}
	}
	for range quads {
		<-recvchan
	}
	for i := range quads {
		sendchan <- callShadowCastMap{
			m:  m,
			r:  r,
			vt: vt.quad(i).swap(),
		}
	}
	for range quads {
		<-recvchan
	}
}

func ShadowCast(m ShadowCastMap, r Radius, source image.Point) {
	vt := visionTransformer{source: source}
	for i := range quads {
		shadowCastOctant(1, one(), zero(), m, r, vt.quad(i))
		shadowCastOctant(1, one(), zero(), m, r, vt.quad(i).swap())
	}
}

func shadowCastOctant(col int, startSlope, endSlope fraction, m ShadowCastMap, r Radius, vt visionTransformer) {
	wall := false
	for x := col; endSlope.less(startSlope) && !wall && r.In(x, 0); x++ {
		approxStart := startSlope.mul(x).whole() + 1
		for y := approxStart; y >= 0; y-- {
			if !r.In(x, y) {
				continue
			}
			hiSlope := newFraction(y*2-1, x*2-1)
			if startSlope.less(hiSlope) {
				continue
			}
			loSlope := newFraction(y*2+1, x*2-1)
			if loSlope.less(endSlope) {
				break
			}
			pos := vt.trans(image.Point{X: x, Y: y})
			m.SetVisible(pos)
			if !m.IsTransparent(pos) {
				if !wall {
					wall = true
					shadowCastOctant(x+1, startSlope, loSlope, m, r, vt)
				}
			} else {
				if wall {
					wall = false
					startSlope = newFraction(y*2+1, x*2+1)
				}
			}
		}
	}
}
