package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/unchartedsoftware/plog"
)

// NewVertFragShader instantiates a new shader object.
func NewVertFragShader(vert, frag string) (*Shader, error) {

	log.Infof("--- %s, %s ---", vert, frag)

	shader := &Shader{}

	vertex, err := shader.CreateShader(vert, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}

	fragment, err := shader.CreateShader(frag, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}

	shader.AttachShader(vertex)
	shader.AttachShader(fragment)

	err = shader.LinkProgram()
	if err != nil {
		return nil, err
	}
	return shader, nil
}
