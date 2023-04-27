package grid

var colorData = []float32{
	// Black to White
	0.0, 0.0, 0.0,
	0.25, 0.25, 0.25,
	0.5, 0.5, 0.5,
	229.0 / 256.0, 229.0 / 256.0, 229.0 / 256.0,
	1.0, 1.0, 1.0,

	// Red
	0.5, 0.0, 0.0,
	205.0 / 256.0, 0.0, 0.0,
	1.0, 0.0, 0.0,

	// Green
	0.0, 0.5, 0.0,
	0.0, 205.0 / 256.0, 0.0,
	0.0, 1.0, 0.0,

	// Blue
	0.0, 0.0, 0.5,
	0.0, 0.0, 205.0 / 256.0,
	0.0, 0.0, 1.0,

	// Yellow
	0.5, 0.5, 0.0,
	205.0 / 256.0, 205.0 / 256.0, 0.0,
	1.0, 1.0, 0.0,

	// Magenta
	0.5, 0.0, 0.5,
	205.0 / 256.0, 0.0, 205.0 / 256.0,
	1.0, 0.0, 1.0,

	// Cyan
	0.0, 0.5, 0.5,
	0.0, 205.0 / 256.0, 205.0 / 256.0,
	0.0, 1.0, 1.0,
}

type Color int

const (
	Black Color = iota
	VeryDarkGray
	DarkGray
	Gray
	White

	DarkRed
	Red
	BrightRed

	DarkGreen
	Green
	BrightGreen

	DarkBlue
	Blue
	BrightBlue

	DarkYellow
	Yellow
	BrightYellow

	DarkMagenta
	Magenta
	BrightMagenta

	DarkCyan
	Cyan
	BrightCyan
)
