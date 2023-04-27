package grid

import (
	"fmt"
	"io"
	"strings"

	"github.com/PieterD/rlpoc/old/internal/internal/pkg/cardinal"
	"github.com/PieterD/rlpoc/old/internal/internal/pkg/lexer"
	"github.com/PieterD/rlpoc/old/internal/internal/spacetime"
)

type portalId byte

type Room struct {
	bounds  cardinal.Box
	cells   map[cardinal.Coord]Cell
	portals []Portal
}

type Cell struct {
	Wall bool
}

type Portal struct {
	enabled        bool
	Location       cardinal.Coord
	LocalPortalId  portalId
	RemoteRoomId   spacetime.EntityId
	RemotePortalId portalId
}

func NewRoomFromSource(name string, source string) (*Room, error) {
	r := &Room{
		cells:   make(map[cardinal.Coord]Cell),
		portals: make([]Portal, 10),
	}
	if err := r.parseRoom(name, strings.NewReader(source)); err != nil {
		return nil, fmt.Errorf("parsing room: %w", err)
	}
	return r, nil
}

func (r *Room) Contains(at cardinal.Coord) bool {
	_, ok := r.cells[at]
	if !ok {
		return false
	}
	return true
}

func (r *Room) IsWall(at cardinal.Coord) bool {
	cell, ok := r.cells[at]
	if !ok {
		return false
	}
	if !cell.Wall {
		return false
	}
	return true
}

func (r *Room) Set(at cardinal.Coord, wall bool) {
	r.Include(at)
	r.cells[at] = Cell{
		Wall: wall,
	}
}

func (r *Room) Clear(at cardinal.Coord) {
	r.Include(at)
	delete(r.cells, at)
}

func (r *Room) Include(at cardinal.Coord) {
	r.bounds = r.bounds.Include(at)
}

func (r *Room) Source() string {
	builder := &strings.Builder{}
	_ = r.bounds.Visit(func(coord cardinal.Coord, start, newRow bool) error {
		if newRow && !start {
			builder.WriteByte('\n')
		}
		if newRow {
			builder.WriteString("grid: ")
		}
		if !r.Contains(coord) {
			builder.WriteByte(' ')
			return nil
		}
		if r.IsWall(coord) {
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

type (
	tokGridRowStart struct{}
	tokGridRowEnd   struct{}
	tokWall         struct{}
	tokFloor        struct{}
	tokSpace        struct{}
	tokPortal       struct {
		PortalId byte
	}
)

func (r *Room) parseRoom(inputName string, input io.Reader) error {
	coord := cardinal.Coord{}
	err := lexer.Run(inputName, input, stInit, func(token interface{}) error {
		switch token.(type) {
		case tokGridRowStart:
		case tokGridRowEnd:
			coord.Y++
			coord.X = 0
		case tokWall:
			r.Set(coord, true)
			coord.X++
		case tokFloor:
			r.Set(coord, false)
			coord.X++
		case tokSpace:
			r.Clear(coord)
			coord.X++
		default:
			return fmt.Errorf("unexpected token: %#v", token)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("running lexer: %w", err)
	}
	return nil
}

func stInit(ctx lexer.Context, in rune) lexer.StateFunc {
	switch in {
	case '/':
		return stComment
	case '\n':
		// ignore empty lines
		ctx.Clear()
		return stInit
	case 'g':
		// grid definition
		return lexer.MatchString("grid: ", tokGridRowStart{}, stGridRow)(ctx, in)
	default:
		return ctx.Error(fmt.Errorf("invalid character: %X'%c' cannot start a line", in, in))
	}
}

func stComment(ctx lexer.Context, in rune) lexer.StateFunc {
	switch in {
	case '/':
		return stCommentContents
	default:
		return ctx.Error(fmt.Errorf("no valid characters after '/' other than another '/' to begin a comment"))
	}
}

func stCommentContents(ctx lexer.Context, in rune) lexer.StateFunc {
	switch in {
	case '\n':
		ctx.Clear()
		return stInit
	default:
		return stCommentContents
	}
}

func stGridRow(ctx lexer.Context, in rune) lexer.StateFunc {
	switch in {
	case '#':
		ctx.Emit(tokWall{})
		ctx.Clear()
		return stGridRow
	case '.':
		ctx.Emit(tokFloor{})
		ctx.Clear()
		return stGridRow
	case ' ':
		ctx.Emit(tokSpace{})
		ctx.Clear()
		return stGridRow
	case '\n':
		ctx.Emit(tokGridRowEnd{})
		ctx.Clear()
		return stInit
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		ctx.Emit(tokPortal{PortalId: byte(in - '0')})
		ctx.Clear()
		return stGridRow
	default:
		return ctx.Error(fmt.Errorf("invalid character in a grid row: %X'%c'", in, in))
	}
}
