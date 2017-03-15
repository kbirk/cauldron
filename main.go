package main

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/unchartedsoftware/plog"

	"github.com/kbirk/cauldron/render"
	"github.com/kbirk/cauldron/shape"
)

const (
	windowWidth  = 1200
	windowHeight = 800
)

var (
	camera             *render.Transform
	explosionTechnique *render.Technique
	smokeTechnique     *render.Technique
	shockwaveTechnique *render.Technique
	effects            []*Effect
	projection         mgl32.Mat4
	view               mgl32.Mat4
)

// Effect represents an animated effect.
type Effect struct {
	Explosion *render.Renderable
	Smoke     *render.Renderable
	Shockwave *render.Renderable
	Time      time.Time
	Position  mgl32.Vec2
}

// Draw renders the effect based on the provided time value.
func (e *Effect) Draw(now time.Time) {
	// time relative to start of effect
	t := float32(now.Sub(e.Time).Seconds())
	// model matrix
	model := mgl32.Translate3D(e.Position[0], e.Position[1], 0.0)
	shockwaveTechnique.Draw(
		drawShockwave(
			e.Shockwave,
			mgl32.Vec4{1.0, 0.98, 0.96, 0.4},
			projection,
			view,
			model,
			t))
	smokeTechnique.Draw(
		drawSmoke(
			e.Smoke,
			mgl32.Vec4{0.41, 0.4, 0.39, 0.2},
			projection,
			view,
			model,
			t))
	explosionTechnique.Draw(
		drawExplosion(
			e.Explosion,
			mgl32.Vec4{0.8, 0.4, 0.2, 0.8},
			projection,
			view,
			model,
			t))
}

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func newFlatTechnique(viewport *render.Viewport) (*render.Technique, error) {
	// create shader
	shader, err := render.NewVertFragShader(
		"resources/shaders/flat.vert",
		"resources/shaders/flat.frag")
	if err != nil {
		return nil, err
	}
	// create technique
	technique := &render.Technique{}
	technique.Enable(gl.BLEND)
	technique.Disable(gl.DEPTH_TEST)
	technique.Shader(shader)
	technique.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	technique.Viewport(viewport)
	// return technique
	return technique, nil
}

func newExplosionTechnique(viewport *render.Viewport) (*render.Technique, error) {
	// create shader
	shader, err := render.NewVertFragShader(
		"resources/shaders/particle.vert",
		"resources/shaders/particle.frag")
	if err != nil {
		return nil, err
	}
	// create technique
	technique := &render.Technique{}
	technique.Enable(gl.BLEND)
	technique.Disable(gl.DEPTH_TEST)
	technique.Shader(shader)
	technique.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	technique.Viewport(viewport)
	// return technique
	return technique, nil
}

func newSmokeTechnique(viewport *render.Viewport) (*render.Technique, error) {
	// create shader
	shader, err := render.NewVertFragShader(
		"resources/shaders/smoke.vert",
		"resources/shaders/smoke.frag")
	if err != nil {
		return nil, err
	}
	// create technique
	technique := &render.Technique{}
	technique.Enable(gl.BLEND)
	technique.Disable(gl.DEPTH_TEST)
	technique.Shader(shader)
	technique.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	technique.Viewport(viewport)
	// return technique
	return technique, nil
}

func newShockwaveTechnique(viewport *render.Viewport) (*render.Technique, error) {
	// create shader
	shader, err := render.NewVertFragShader(
		"resources/shaders/shockwave.vert",
		"resources/shaders/shockwave.frag")
	if err != nil {
		return nil, err
	}
	// create technique
	technique := &render.Technique{}
	technique.Enable(gl.BLEND)
	technique.Disable(gl.DEPTH_TEST)
	technique.Shader(shader)
	technique.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	technique.Viewport(viewport)
	// return technique
	return technique, nil
}

func randVec2() mgl32.Vec2 {
	return mgl32.Vec2{
		(rand.Float32()*2 - 1),
		(rand.Float32()*2 - 1),
	}.Normalize()
}

func createParticleRenderable(positions []float32, indices []uint16, offsets []float32, velocities []float32, sizes []float32) *render.Renderable {
	// create vertexbuffer
	vb := &render.VertexBuffer{}
	vb.AllocateBuffer((len(positions) + len(offsets) + len(velocities) + len(sizes)) * 4)
	vb.BufferSubFloat32(positions, 0)
	vb.BufferSubFloat32(offsets, len(positions)*4)
	vb.BufferSubFloat32(velocities, (len(positions)+len(offsets))*4)
	vb.BufferSubFloat32(sizes, (len(positions)+len(offsets)+len(velocities))*4)
	// create indexbuffer
	ib := &render.IndexBuffer{}
	ib.BufferUint16(indices)
	// create renderable
	renderable := &render.Renderable{}
	renderable.SetVertexBuffer(vb)
	renderable.SetIndexBuffer(ib)
	renderable.SetPointer(0, &render.AttributePointer{
		Type:       gl.FLOAT,
		Size:       3,
		ByteOffset: 0,
	})
	renderable.SetPointer(1, &render.AttributePointer{
		Type:       gl.FLOAT,
		Size:       2,
		ByteOffset: len(positions) * 4,
	})
	renderable.SetPointer(2, &render.AttributePointer{
		Type:       gl.FLOAT,
		Size:       2,
		ByteOffset: (len(positions) + len(offsets)) * 4,
	})
	renderable.SetPointer(3, &render.AttributePointer{
		Type:       gl.FLOAT,
		Size:       1,
		ByteOffset: (len(positions) + len(offsets) + len(velocities)) * 4,
	})
	renderable.SetInstancedAttributes([]uint32{
		1, 2, 3,
	})
	renderable.SetDrawElementsInstanced(
		gl.TRIANGLES,
		int32(len(indices)),
		gl.UNSIGNED_SHORT,
		0,
		int32(len(offsets)/2))
	renderable.Upload()
	return renderable
}

func createExplosion(num int, radius float32, force float32, size float32) *render.Renderable {
	positions, indices := shape.Quad(size, true, false)
	offsets := make([]float32, 2*num)
	for i := 0; i < num; i++ {
		offset := randVec2().Mul(rand.Float32() * radius)
		offsets[i*2] = offset[0]
		offsets[i*2+1] = offset[1]
	}
	velocities := make([]float32, 2*num)
	direction := mgl32.Vec2{0, force * 2}
	for i := 0; i < num; i++ {
		velocity := randVec2().Mul(rand.Float32() * force).Add(direction)
		velocities[i*2] = velocity[0]
		velocities[i*2+1] = velocity[1]
	}
	sizes := make([]float32, num)
	for i := 0; i < num; i++ {
		// size
		sizes[i] = rand.Float32() * size
	}
	return createParticleRenderable(
		positions,
		indices,
		offsets,
		velocities,
		sizes)
}

func createSmoke(num int, radius float32, force float32, size float32) *render.Renderable {
	positions, indices := shape.Circle(size, 64, true, false)
	offsets := make([]float32, 2*num)
	for i := 0; i < num; i++ {
		offset := randVec2().Mul(rand.Float32() * radius)
		offsets[i*2] = offset[0]
		offsets[i*2+1] = offset[1]
	}
	velocities := make([]float32, 2*num)
	direction := mgl32.Vec2{0, force * 1.5}
	for i := 0; i < num; i++ {
		velocity := randVec2().Mul(rand.Float32() * force).Add(direction)
		velocities[i*2] = velocity[0]
		velocities[i*2+1] = velocity[1]
	}
	sizes := make([]float32, num)
	for i := 0; i < num; i++ {
		// size
		sizes[i] = ((rand.Float32() * 0.5) + 0.5) * size
	}
	return createParticleRenderable(
		positions,
		indices,
		offsets,
		velocities,
		sizes)
}

func createQuad(size float32) *render.Renderable {
	vertices, indices := shape.Quad(size, true, true)
	// create vertexbuffer
	vb := &render.VertexBuffer{}
	vb.BufferFloat32(vertices)
	// create indexbuffer
	ib := &render.IndexBuffer{}
	ib.BufferUint16(indices)
	// create renderable
	quad := &render.Renderable{}
	quad.SetVertexBuffer(vb)
	quad.SetIndexBuffer(ib)
	quad.SetPointer(0, &render.AttributePointer{
		Type:       gl.FLOAT,
		Size:       3,
		ByteStride: (3 + 2) * 4,
		ByteOffset: 0,
	})
	quad.SetPointer(1, &render.AttributePointer{
		Type:       gl.FLOAT,
		Size:       2,
		ByteStride: (3 + 2) * 4,
		ByteOffset: 3 * 4,
	})
	quad.SetDrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_SHORT, 0)
	quad.Upload()
	return quad
}

func createCircle(radius float32, segments int) *render.Renderable {
	vertices, indices := shape.Circle(radius, segments, true, true)
	// create vertexbuffer
	vb := &render.VertexBuffer{}
	vb.BufferFloat32(vertices)
	// create indexbuffer
	ib := &render.IndexBuffer{}
	ib.BufferUint16(indices)
	// create renderable
	circle := &render.Renderable{}
	circle.SetVertexBuffer(vb)
	circle.SetIndexBuffer(ib)
	circle.SetPointer(0, &render.AttributePointer{
		Type:       gl.FLOAT,
		Size:       3,
		ByteStride: (3 + 2) * 4,
		ByteOffset: 0,
	})
	circle.SetPointer(1, &render.AttributePointer{
		Type:       gl.FLOAT,
		Size:       2,
		ByteStride: (3 + 2) * 4,
		ByteOffset: 3 * 4,
	})
	circle.SetDrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_SHORT, 0)
	circle.Upload()
	return circle
}

func drawFlat(renderable *render.Renderable, color mgl32.Vec4, projection, view, model mgl32.Mat4) []*render.Command {
	command := &render.Command{}
	command.Uniform("uProjection", &projection[0])
	command.Uniform("uView", &view[0])
	command.Uniform("uModel", &model[0])
	command.Uniform("uColor", &color[0])
	command.Renderable(renderable)
	return []*render.Command{
		command,
	}
}

func drawExplosion(renderable *render.Renderable, color mgl32.Vec4, projection, view, model mgl32.Mat4, time float32) []*render.Command {
	command := &render.Command{}
	command.Uniform("uProjection", &projection[0])
	command.Uniform("uView", &view[0])
	command.Uniform("uModel", &model[0])
	command.Uniform("uColor", &color[0])
	gravity := mgl32.Vec2{0, -200}
	command.Uniform("uGravity", &gravity[0])
	command.Uniform("uTime", time)
	command.Renderable(renderable)
	return []*render.Command{
		command,
	}
}

func drawSmoke(renderable *render.Renderable, color mgl32.Vec4, projection, view, model mgl32.Mat4, time float32) []*render.Command {
	command := &render.Command{}
	command.Uniform("uProjection", &projection[0])
	command.Uniform("uView", &view[0])
	command.Uniform("uModel", &model[0])
	command.Uniform("uColor", &color[0])
	rise := mgl32.Vec2{0, 10}
	command.Uniform("uRise", &rise[0])
	command.Uniform("uTime", time)
	command.Renderable(renderable)
	return []*render.Command{
		command,
	}
}

func drawShockwave(renderable *render.Renderable, color mgl32.Vec4, projection, view, model mgl32.Mat4, time float32) []*render.Command {
	command := &render.Command{}
	command.Uniform("uProjection", &projection[0])
	command.Uniform("uView", &view[0])
	command.Uniform("uModel", &model[0])
	command.Uniform("uColor", &color[0])
	command.Uniform("uForce", float32(15.0))
	command.Uniform("uTime", time)
	command.Renderable(renderable)
	return []*render.Command{
		command,
	}
}

func handleKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
}

func handleMouseButton(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action == glfw.Press {
		x, y := w.GetCursorPos()
		_, height := w.GetSize()
		effects = append(effects, &Effect{
			Explosion: createExplosion(500, 20, 200, 4),
			Smoke:     createSmoke(1000, 20, 140, 10),
			Shockwave: createCircle(10.0, 64),
			Position:  mgl32.Vec2{float32(x), float32(float64(height) - y)},
			Time:      time.Now(),
		})
	}
}

func main() {

	// init glfw
	err := glfw.Init()
	if err != nil {
		log.Error("Failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	// set window hints
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// create window
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "cauldron", nil, nil)
	if err != nil {
		log.Error("Failed to create window:", err)
		return
	}
	window.MakeContextCurrent()
	window.SetKeyCallback(handleKey)
	window.SetMouseButtonCallback(handleMouseButton)

	// init glow
	err = gl.Init()
	if err != nil {
		log.Error("Failed to init glow:", err)
		return
	}

	// log opengl version
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Info("OpenGL version", version)

	// create viewport
	viewportWidth, viewportHeight := window.GetFramebufferSize()
	viewport := &render.Viewport{
		X:      0,
		Y:      0,
		Width:  int32(viewportWidth),
		Height: int32(viewportHeight),
	}

	// create camera
	camera = render.NewTransform()

	// create techniques
	explosionTechnique, err = newExplosionTechnique(viewport)
	if err != nil {
		log.Error(err)
		return
	}
	smokeTechnique, err = newSmokeTechnique(viewport)
	if err != nil {
		log.Error(err)
		return
	}
	shockwaveTechnique, err = newShockwaveTechnique(viewport)
	if err != nil {
		log.Error(err)
		return
	}

	// projection matrix
	width, height := window.GetSize()
	projection = mgl32.Ortho(
		0, float32(width),
		0, float32(height),
		-1.0, 1.0)

	// view matrix
	view = camera.ViewMatrix()

	// Configure global settings
	gl.ClearColor(0.1, 0.1, 0.1, 1.0)

	// frame loop
	for !window.ShouldClose() {

		// poll events
		glfw.PollEvents()

		// clear buffer
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// grab current time
		now := time.Now()

		// draw animations
		for _, effect := range effects {
			effect.Draw(now)
		}

		// remove stale effects
		j := 0
		for i := 0; i < len(effects); i++ {
			if now.Sub(effects[i].Time).Seconds() < 3.0 {
				effects[j] = effects[i]
				j++
			}
		}
		effects = effects[:j]

		// swap buffers
		window.SwapBuffers()
	}
}
