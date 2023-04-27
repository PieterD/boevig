package vision

import (
	"github.com/PieterD/boevig/pkg/frac"
	"sort"
)

type shadeRepo struct {
	shades []shade
}

func newShadeRepo() *shadeRepo {
	return &shadeRepo{}
}

type shade struct {
	From frac.Fraction
	To   frac.Fraction
}

func (s shade) Range(y int) (fromX, toX int) {
	panic("not implemented")
}

func (s shade) Merge(s2 shade) (s3 shade, ok bool) {
	if s2.From.Less(s.From) {
		s, s2 = s2, s
	}
	if s.To.Less(s2.From) {
		return shade{}, false
	}
	return shade{
		From: s.From,
		To:   s2.To,
	}, true
}

func (s shade) Contains(what frac.Fraction) bool {
	if what.Less(s.From) {
		return false
	}
	if s.To.Less(what) {
		return false
	}
	return true
}

func (r *shadeRepo) FindNextShade(start frac.Fraction) (s shade, ok bool) {
	for _, search := range r.shades {
		if search.Contains(start) {
			return s, true
		}
	}
	return shade{}, false
}

func (r *shadeRepo) InsertShade(s shade) {
	r.shades = append(r.shades, s)
	if len(r.shades) == 1 {
		return
	}
	sort.Slice(r.shades, func(i, j int) bool {
		return r.shades[i].From.Less(r.shades[j].From)
	})
	newShades := []shade{r.shades[0]}
	j := 0
	for i := 1; i < len(r.shades); i++ {
		m, ok := newShades[j].Merge(r.shades[i])
		if ok {
			newShades[j] = m
			continue
		}
		j++
		newShades = append(newShades, r.shades[i])
	}
	sort.Slice(r.shades, func(i, j int) bool {
		return r.shades[i].From.Less(r.shades[j].From)
	})
	r.shades = newShades
}

func (r *shadeRepo) Shades() []shade {
	return r.shades
}
