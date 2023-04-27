package sim

import (
	"fmt"
	"io"

	"github.com/PieterD/rlpoc/old/internal/lib/cardinal"
	"github.com/PieterD/rlpoc/old/internal/lib/lexer"
)

type (
	tokGridRowStart struct{}
	tokGridRowEnd   struct{}
	tokWall         struct{}
	tokFloor        struct{}
	tokSpace        struct{}
)

func (g *Grid) parseGrid(inputName string, input io.Reader) error {
	coord := cardinal.Coord{}
	err := lexer.Run(inputName, input, stInit, func(token interface{}) error {
		switch token.(type) {
		case tokGridRowStart:
		case tokGridRowEnd:
			coord.Y++
			coord.X = 0
		case tokWall:
			g.Set(coord, true)
			coord.X++
		case tokFloor:
			g.Set(coord, false)
			coord.X++
		case tokSpace:
			g.Clear(coord)
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
	default:
		return ctx.Error(fmt.Errorf("invalid character in a grid row: %X'%c'", in, in))
	}
}
