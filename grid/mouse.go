package grid

import (
	"image"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type mouseTranslator struct {
	grid   *Grid
	eh     EventHandler
	last   image.Point
	down   bool
	drag   bool
	button MouseButton
	start  image.Point
}

func newMouseTranslator(grid *Grid, eh EventHandler) *mouseTranslator {
	return &mouseTranslator{
		grid: grid,
		eh:   eh,
		last: image.Point{X: -1, Y: -1},
	}
}

func (trans *mouseTranslator) Pos(posx, posy float64) {
	x := (int(posx) - trans.grid.padx)
	y := (int(posy) - trans.grid.pady)

	if x < 0 || y < 0 {
		return
	}

	x /= trans.grid.runewidth
	y /= trans.grid.runeheight

	if trans.last.X == x && trans.last.Y == y {
		return
	}
	if x >= trans.grid.cols || y >= trans.grid.rows {
		return
	}

	trans.last.X = x
	trans.last.Y = y

	if trans.down && !trans.drag {
		trans.drag = true
		trans.eh.MouseDrag(MouseDragEvent{
			To:     trans.last,
			From:   trans.start,
			State:  StartDrag,
			Button: trans.button,
		})
		return
	}

	if trans.down && trans.drag {
		trans.eh.MouseDrag(MouseDragEvent{
			To:     trans.last,
			From:   trans.start,
			State:  ContinueDrag,
			Button: trans.button,
		})
		return
	}

	trans.eh.MouseMove(MouseMoveEvent{
		Pos: trans.last,
	})
}

type MouseButton int

const (
	MouseButtonLeft   MouseButton = MouseButton(glfw.MouseButtonLeft)
	MouseButtonRight  MouseButton = MouseButton(glfw.MouseButtonRight)
	MouseButtonMiddle MouseButton = MouseButton(glfw.MouseButtonMiddle)
)

func (trans *mouseTranslator) Button(gbutton glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if gbutton != glfw.MouseButtonLeft && gbutton != glfw.MouseButtonRight && gbutton != glfw.MouseButtonMiddle {
		return
	}
	button := MouseButton(gbutton)
	if action == glfw.Press {
		if trans.down {
			return
		}
		trans.down = true
		trans.drag = false
		trans.button = button
		trans.start = trans.last
		return
	}
	if action == glfw.Release {
		if !trans.down || trans.button != button {
			return
		}
		if trans.drag {
			trans.eh.MouseDrag(MouseDragEvent{
				To:     trans.last,
				From:   trans.start,
				State:  EndDrag,
				Button: trans.button,
			})
		} else {
			trans.eh.MouseClick(MouseClickEvent{
				Pos:    trans.last,
				Button: trans.button,
			})
		}
		trans.down = false
		trans.drag = false
		return
	}
}

type MouseMoveEvent struct {
	Pos image.Point
}

type MouseClickEvent struct {
	Pos    image.Point
	Button MouseButton
}

type MouseDragEvent struct {
	From   image.Point
	To     image.Point
	State  DragState
	Button MouseButton
}

type DragState int

const (
	StartDrag DragState = iota
	ContinueDrag
	EndDrag
)
