package ansi

import (
	"bytes"
	"fmt"
	"io"
)

type Generator struct {
	buf *bytes.Buffer
	esc byte
}

const defaultEscape = 27

func NewGenerator() *Generator {
	return &Generator{
		buf: bytes.NewBuffer(nil),
		esc: defaultEscape,
	}
}

func (g *Generator) Write(buf []byte) (n int, err error) {
	n, _ = g.buf.Write(buf)
	return n, nil
}

func (g *Generator) Flush(w io.Writer) error {
	_, err := io.Copy(w, g.buf)
	if err != nil {
		return fmt.Errorf("copying buffer to writer: %w", err)
	}
	return nil
}

func (g *Generator) write(bytes ...byte) {
	_, _ = g.buf.Write(bytes)
}

func (g *Generator) CursorPos(row, col uint) {
	_, _ = fmt.Fprintf(g, "%c[%d;%df", g.esc, row, col)
}

func (g *Generator) ClearScreen() {
	_, _ = fmt.Fprintf(g, "%c[2J", g.esc)
}

func (g *Generator) ForegroundColor(color Color, bright bool) {
	if bright {
		color += brightAdditive
	}
	_, _ = fmt.Fprintf(g, "%c[%dm", g.esc, color)
}

func (g *Generator) BackgroundColor(color Color) {
	_, _ = fmt.Fprintf(g, "%c[%dm", g.esc, color+backgroundAdditive)
}

func (g *Generator) Reset() {
	_, _ = fmt.Fprintf(g, "%c[0m", g.esc)
}

func (g *Generator) Printf(f string, args ...interface{}) {
	_, _ = fmt.Fprintf(g, f, args...)
}
