package sim

import (
	"fmt"
	"strings"

	"github.com/PieterD/rlpoc/old/internal/lib/cardinal"
)

type Grid struct {
	bounds cardinal.Box
	cells  map[cardinal.Coord]Cell
}

type Cell struct {
	Wall bool
}

func NewGridFromSource(source string) (*Grid, error) {
	g := &Grid{
		cells: make(map[cardinal.Coord]Cell),
	}
	if err := g.parseGrid("HC_FIRST_ROOM", strings.NewReader(source)); err != nil {
		return nil, fmt.Errorf("parsing grid: %w", err)
	}
	return g, nil
}

func (g *Grid) Contains(at cardinal.Coord) bool {
	_, ok := g.cells[at]
	if !ok {
		return false
	}
	return true
}

func (g *Grid) IsWall(at cardinal.Coord) bool {
	cell, ok := g.cells[at]
	if !ok {
		return false
	}
	if !cell.Wall {
		return false
	}
	return true
}

func (g *Grid) Set(at cardinal.Coord, wall bool) {
	g.Include(at)
	g.cells[at] = Cell{
		Wall: wall,
	}
}

func (g *Grid) Clear(at cardinal.Coord) {
	g.Include(at)
	delete(g.cells, at)
}

func (g *Grid) Include(at cardinal.Coord) {
	g.bounds = g.bounds.Include(at)
}

func (g *Grid) Source() string {
	builder := &strings.Builder{}
	_ = g.bounds.Visit(func(coord cardinal.Coord, start, newRow bool) error {
		if newRow && !start {
			builder.WriteByte('\n')
		}
		if newRow {
			builder.WriteString("grid: ")
		}
		if !g.Contains(coord) {
			builder.WriteByte(' ')
			return nil
		}
		if g.IsWall(coord) {
			builder.WriteByte('#')
			return nil
		} else {
			builder.WriteByte('.')
			return nil
		}
	})
	builder.WriteByte('\n')
	return builder.String()
}
