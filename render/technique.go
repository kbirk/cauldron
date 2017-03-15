package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type blendFunc struct {
	sfactor uint32
	dfactor uint32
}

type cullFace struct {
	mode uint32
}

type depthMask struct {
	flag bool
}

type depthFunc struct {
	xfunc uint32
}

type clearColor struct {
	r float32
	g float32
	b float32
	a float32
}

// Technique represents a render technique.
type Technique struct {
	enables    []uint32
	disables   []uint32
	shader     *Shader
	viewport   *Viewport
	blendFunc  *blendFunc
	cullFace   *cullFace
	depthMask  *depthMask
	depthFunc  *depthFunc
	clearColor *clearColor
}

// Enable enables the rendering states for the technique.
func (t *Technique) Enable(enable uint32) {
	t.enables = append(t.enables, enable)
}

// Disable disables the rendering states for the technique.
func (t *Technique) Disable(disable uint32) {
	t.disables = append(t.disables, disable)
}

// Shader sets the shader for the technique.
func (t *Technique) Shader(shader *Shader) {
	t.shader = shader
}

// Viewport sets the viewport for the technique.
func (t *Technique) Viewport(viewport *Viewport) {
	t.viewport = viewport
}

// BlendFunc sets the blend func for the technique.
func (t *Technique) BlendFunc(sfactor uint32, dfactor uint32) {
	t.blendFunc = &blendFunc{
		sfactor: sfactor,
		dfactor: dfactor,
	}
}

// CullFace sets the cull face mode for the technique.
func (t *Technique) CullFace(mode uint32) {
	t.cullFace = &cullFace{
		mode: mode,
	}
}

// DepthMask sets the depth mask for the technique.
func (t *Technique) DepthMask(flag bool) {
	t.depthMask = &depthMask{
		flag: flag,
	}
}

// DepthFunc sets the depth mask for the technique.
func (t *Technique) DepthFunc(xfunc uint32) {
	t.depthFunc = &depthFunc{
		xfunc: xfunc,
	}
}

// ClearColor sets the clear color for the frame.
func (t *Technique) ClearColor(r, g, b, a float32) {
	t.clearColor = &clearColor{
		r: r,
		g: g,
		b: b,
		a: a,
	}
}

// Draw renders all commands using the technique.
func (t *Technique) Draw(commands []*Command) {
	t.bind()
	for _, command := range commands {
		command.Execute(t.shader)
	}
	t.unbind()
}

func (t *Technique) bind() {

	// use shader
	t.shader.Use()

	// enable state
	for _, state := range t.enables {
		gl.Enable(state)
	}
	// disable state
	for _, state := range t.disables {
		gl.Disable(state)
	}

	// state functions
	if t.blendFunc != nil {
		gl.BlendFunc(t.blendFunc.sfactor, t.blendFunc.dfactor)
	}
	if t.cullFace != nil {
		gl.CullFace(t.cullFace.mode)
	}
	if t.depthMask != nil {
		gl.DepthMask(t.depthMask.flag)
	}
	if t.depthFunc != nil {
		gl.DepthFunc(t.depthFunc.xfunc)
	}

	// update viewport
	if t.viewport != nil {
		gl.Viewport(
			t.viewport.X,
			t.viewport.Y,
			t.viewport.Width,
			t.viewport.Height)
	}
}

func (t *Technique) unbind() {
	// TODO: shouldn't require this method, only change what is needed in Bind

	// disable state
	// for _, state := range t.enables {
	// 	gl.Disable(state)
	// }

	// reset state functions
	if t.blendFunc != nil {
		gl.BlendFunc(gl.ONE, gl.ZERO)
	}
	if t.cullFace != nil {
		gl.CullFace(gl.BACK)
	}
	if t.depthMask != nil {
		gl.DepthMask(true)
	}
	if t.depthFunc != nil {
		gl.DepthFunc(gl.LESS)
	}
}
