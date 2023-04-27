package ansi

type Color byte

const (
	Black Color = 30 + iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

const (
	backgroundAdditive Color = 10
	brightAdditive     Color = 60
)
