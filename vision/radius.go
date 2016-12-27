package vision

type Radius interface {
	In(x, y int) bool
}

type endlessRadius struct{}

func EndlessRadius() Radius {
	return endlessRadius{}
}

func (_ endlessRadius) In(x, y int) bool {
	return true
}

type circularRadius struct {
	radius        int
	radiusSquared int64
}

func CircularRadius(r int) Radius {
	rr := int64(r)
	return circularRadius{
		radius:        r,
		radiusSquared: rr * rr,
	}
}

func (r circularRadius) In(x, y int) bool {
	xx := int64(x)
	yy := int64(y)
	return xx*xx+yy*yy <= r.radiusSquared
}
