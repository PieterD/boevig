package grid

import (
	"image"
	"runtime"

	. "github.com/PieterD/boevig/pan"
	"github.com/PieterD/glimmer/gli"
	"github.com/PieterD/glimmer/win"
)

type EventHandler interface {
	Draw(g DrawableGrid)
	Char(r rune)
	Key(k KeyEvent)
	MouseMove(k MouseMoveEvent)
	MouseClick(k MouseClickEvent)
	MouseDrag(k MouseDragEvent)
	Fin(last bool) bool
}

func init() {
	runtime.LockOSThread()
}

var vertexShaderText = `
#version 110
attribute vec2 position;
attribute float foreColor;
attribute float backColor;
attribute vec2 texCoord;
uniform vec3 colorData[23];
uniform vec2 runeSize;
varying vec3 theForeColor;
varying vec3 theBackColor;
varying vec2 theTexCoord;
void main() {
	gl_Position = vec4(position, 0.0, 1.0);
	theForeColor = colorData[int(foreColor)];
	theBackColor = colorData[int(backColor)];
	theTexCoord = vec2(texCoord.x / runeSize.x, texCoord.y / runeSize.y);
}
`

var fragmentShaderText = `
#version 110
varying vec3 theForeColor;
varying vec3 theBackColor;
varying vec2 theTexCoord;
uniform sampler2D tex;
void main() {
	vec4 texColor = texture2D(tex, theTexCoord);
	gl_FragColor = vec4(mix(theBackColor, theForeColor, texColor.a), 1.0);
}
`

func Run(charset string, charwidth, charheight int, eh EventHandler) {
	defer eh.Fin(true)
	width := 800
	height := 600
	Panic(win.Start(
		win.Size(width, height),
		win.Resizable(),
		win.Func(func (window *win.Window){
			// Create shaders and program
			program, err := gli.NewProgram(vertexShaderText, fragmentShaderText)
			Panic(err)
			defer program.Delete()

			// Load and initialize texture
			img, err := gli.LoadImage(charset)
			Panic(err)
			texture, err := gli.NewTexture(img,
				gli.TextureFilter(gli.LINEAR, gli.LINEAR),
				gli.TextureWrap(gli.CLAMP_TO_EDGE, gli.CLAMP_TO_EDGE))
			Panic(err)
			defer texture.Delete()

			// Create Vertex ArrayObject
			vao, err := gli.NewVAO()
			Panic(err)
			defer vao.Delete()

			// Create grid
			grid, err := NewGrid(charwidth, charheight, texture.Size().X, texture.Size().Y)
			Panic(err)
			grid.Resize(width, height)
			vCoords, vIndex, vData := grid.Buffers()

			// Create grid buffers
			posvbo, err := gli.NewBuffer(vCoords)
			Panic(err)
			defer posvbo.Delete()
			idxvbo, err := gli.NewBuffer(vIndex, gli.BufferElementArray())
			Panic(err)
			defer idxvbo.Delete()
			vbo, err := gli.NewBuffer(vData, gli.BufferAccessFrequency(gli.DYNAMIC))
			Panic(err)
			defer vbo.Delete()

			mousetrans := newMouseTranslator(grid, eh)
			keytrans := newKeyTranslator()

			// Set up VAO
			vao.Enable(2, posvbo, program.Attrib("position"))
			vao.Enable(2, vbo, program.Attrib("texCoord"),
				gli.VAOStride(4))
			vao.Enable(1, vbo, program.Attrib("foreColor"),
				gli.VAOStride(4), gli.VAOOffset(2))
			vao.Enable(1, vbo, program.Attrib("backColor"),
				gli.VAOStride(4), gli.VAOOffset(3))

			// Set uniforms
			program.Uniform("tex").SetSampler(1)
			program.Uniform("colorData[0]").SetFloat(colorData...)
			program.Uniform("runeSize").SetFloat(float32(grid.RuneSize().X), float32(grid.RuneSize().Y))

			draw, err := gli.NewDraw(gli.TRIANGLES, program, vao,
				gli.DrawIndex(idxvbo),
				gli.DrawTexture(texture, 1))
			Panic(err)

			clear, err := gli.NewClear(gli.ClearColor(0, 0, 0, 1))
			Panic(err)

			for {
				ie := window.Poll()
				if ie == nil {
					if eh.Fin(false) {
						break
					}
					// Render scene
					grid.clearData()
					eh.Draw(grid)
					_, _, vData = grid.Buffers()
					vbo.Update(0, vData)

					// Draw scene
					clear.Clear()
					//texture.Use(1)
					draw.Draw(0, grid.Vertices())

					window.Swap()
					continue
				}

				switch e := ie.(type) {
				case win.EventMouseButton:
					mousetrans.Button(e.Button, e.Action, e.Mod)
				case win.EventMousePos:
					mousetrans.Pos(e.X, e.Y)
				case win.EventChar:
					eh.Char(e.Char)
				case win.EventKey:
					ev, ok := keytrans.Key(e.Key, e.Action, e.Mod)
					if ok {
						eh.Key(ev)
					}
				case win.EventResize:
					width = e.Width
					height = e.Height
					gli.Viewport(image.Rectangle{Max: image.Point{X: width, Y: height}})
					grid.Resize(width, height)
					vCoords, vIndex, vData := grid.Buffers()
					posvbo.Upload(vCoords)
					idxvbo.Upload(vIndex)
					vbo.Upload(vData)

				}
			}
		})))
}
