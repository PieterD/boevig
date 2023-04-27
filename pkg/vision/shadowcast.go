package vision

import (
	"github.com/PieterD/boevig/pkg/core"
	"github.com/PieterD/boevig/pkg/frac"
)

type LocalMap interface {
	IsTileOpaque(entityID core.EntityID, where core.LocalCoordinate) bool
	SetTileVisibility(entityID core.EntityID, where core.LocalCoordinate)
}

type ShadowCaster struct {
	LocalMap LocalMap
}

func NewShadowCaster(lm LocalMap) *ShadowCaster {
	return &ShadowCaster{
		LocalMap: lm,
	}
}

func (sc ShadowCaster) Cast(entityID core.EntityID, source core.LocalCoordinate, radius int) {
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				oc := octantCaster{
					entityID:       entityID,
					shadeRepo:      newShadeRepo(),
					offset:         source,
					radius:         radius,
					flipDiagonal:   k == 1,
					flipHorizontal: j == 1,
					flipVertical:   i == 1,
				}
				oc.Cast()
			}
		}
	}
}

type octantCaster struct {
	localMap       LocalMap
	shadeRepo      *shadeRepo
	entityID       core.EntityID
	offset         core.LocalCoordinate
	radius         int
	flipDiagonal   bool
	flipHorizontal bool
	flipVertical   bool
}

func (oc octantCaster) SetTileVisibility(where core.LocalCoordinate) {

}

func (oc octantCaster) IsTileOpaque(where core.LocalCoordinate) bool {

}

func (oc octantCaster) Cast() {
	if oc.radius <= 0 {
		return
	}
	if oc.radius == 1 {
		oc.SetTileVisibility(core.LocalCoordinate{X: 0, Y: 0})
		return
	}
	if oc.IsTileOpaque(core.LocalCoordinate{X: 0, Y: 0}) {
		oc.SetTileVisibility(core.LocalCoordinate{X: 0, Y: 0})
		return
	}
	for y := 1; y < oc.radius; y++ {
		x := 0
		sweep := oc.LeftAngle(x, y)
		shades := oc.shadeRepo.Shades()
		for _, s := range shades {
		}
	}
}

func (oc octantCaster) LeftAngle(x, y int) frac.Fraction {
	return frac.New(int64(x)*2-1, int64(y)*2-1)
}
