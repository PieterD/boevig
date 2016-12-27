package grid

import (
	"fmt"
	"image"
)

type DrawableGrid interface {
	GridSize() image.Rectangle
	Set(x, y, r int, fore, back Color)
}

type Grid struct {
	screenwidth  int
	screenheight int
	runewidth    int
	runeheight   int
	texwidth     int
	texheight    int

	cols int
	rows int
	padx int
	pady int

	texcols  int
	texrows  int
	texcodes int

	coord []float32
	index []uint32
	data  []uint8
}

func (grid *Grid) GridSize() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: grid.cols, Y: grid.rows},
	}
}

func (grid *Grid) RuneSize() image.Point {
	return image.Point{X: grid.runewidth, Y: grid.runeheight}
}

func NewGrid(runewidth, runeheight, texwidth, texheight int) (*Grid, error) {
	if texwidth%runewidth != 0 {
		return nil, fmt.Errorf("Texture width must be a multiple of the rune width")
	}
	if texheight%runeheight != 0 {
		return nil, fmt.Errorf("Texture height must be a multiple of the rune height")
	}
	grid := &Grid{
		runewidth:  runewidth,
		runeheight: runeheight,
		texwidth:   texwidth,
		texheight:  texheight,
	}

	grid.texcols = grid.texwidth / grid.runewidth
	grid.texrows = grid.texheight / grid.runeheight
	grid.texcodes = grid.texcols * grid.texrows
	return grid, nil
}

func (grid *Grid) Resize(width, height int) {
	grid.screenwidth = width
	grid.screenheight = height
	grid.cols = grid.screenwidth / grid.runewidth
	grid.rows = grid.screenheight / grid.runeheight
	grid.padx = (grid.screenwidth - grid.cols*grid.runewidth) / 2
	grid.pady = (grid.screenheight - grid.rows*grid.runeheight) / 2
	grid.updateSlices()
	grid.updateCoordinates()
	grid.clearData()
}

func (grid *Grid) Buffers() ([]float32, []uint32, []uint8) {
	return grid.coord, grid.index, grid.data
}

func (grid *Grid) updateSlices() {
	coordSize := grid.cols * grid.rows * 2 * 4
	indexSize := grid.cols * grid.rows * 6
	dataSize := grid.cols * grid.rows * 4 * 4
	if cap(grid.coord) < coordSize {
		grid.coord = make([]float32, coordSize)
	} else {
		grid.coord = grid.coord[:coordSize]
	}

	if cap(grid.index) < indexSize {
		grid.index = make([]uint32, indexSize)
	} else {
		grid.index = grid.index[:indexSize]
	}

	if cap(grid.data) < dataSize {
		grid.data = make([]uint8, dataSize)
	} else {
		grid.data = grid.data[:dataSize]
	}
}

func (grid *Grid) updateCoordinates() {
	vpos := uint32(0)
	ipos := 0
	for y := 0; y < grid.rows; y++ {
		for x := 0; x < grid.cols; x++ {
			vstart := (y*grid.cols + x) * 2 * 4
			vcur := grid.coord[vstart : vstart+2*4]
			vcur[0], vcur[1] = grid.fcoord(x, y+1)
			vcur[2], vcur[3] = grid.fcoord(x+1, y+1)
			vcur[4], vcur[5] = grid.fcoord(x+1, y)
			vcur[6], vcur[7] = grid.fcoord(x, y)
			icur := grid.index[ipos : ipos+6]
			icur[0] = vpos
			icur[1] = vpos + 1
			icur[2] = vpos + 2
			icur[3] = vpos
			icur[4] = vpos + 2
			icur[5] = vpos + 3
			vpos += 4
			ipos += 6
		}
	}
}

func (grid *Grid) clearData() {
	for i := range grid.data {
		grid.data[i] = 0
	}
}

func (grid *Grid) fcoord(x, y int) (fx float32, fy float32) {
	fwidth := float32(grid.screenwidth)
	fheight := float32(grid.screenheight)
	itlx := grid.padx + x*grid.runewidth
	itly := grid.pady + y*grid.runeheight
	fx = float32(itlx)/fwidth*2.0 - 1.0
	fy = float32(itly)/fheight*2.0 - 1.0
	return fx, fy
}

func (grid *Grid) Set(x, y, r int, fore, back Color) {
	y = grid.rows - 1 - y
	if x >= grid.cols || y >= grid.rows {
		return
	}
	idx := x + y*grid.cols
	rx := r % grid.texcols
	ry := r / grid.texcols
	data := grid.data[idx*4*4 : idx*4*4+4*4]
	data1 := data[0:4]
	data2 := data[4:8]
	data3 := data[8:12]
	data4 := data[12:16]
	data1[0] = uint8(rx)
	data1[1] = uint8(ry)
	data1[2] = uint8(fore)
	data1[3] = uint8(back)

	data2[0] = uint8(rx + 1)
	data2[1] = uint8(ry)
	data2[2] = uint8(fore)
	data2[3] = uint8(back)

	data3[0] = uint8(rx + 1)
	data3[1] = uint8(ry + 1)
	data3[2] = uint8(fore)
	data3[3] = uint8(back)

	data4[0] = uint8(rx)
	data4[1] = uint8(ry + 1)
	data4[2] = uint8(fore)
	data4[3] = uint8(back)
}

func (grid *Grid) Vertices() int32 {
	return int32(len(grid.index))
}
