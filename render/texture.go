package render

import (
	"fmt"
	"image"
	"image/draw"
	// register png decoder
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// Texture represents a 2D texture object.
type Texture struct {
	id     uint32
	width  uint32
	height uint32
}

// LoadRGBATexture loads an image file into an RGBA texture.
func LoadRGBATexture(filename string, wrap int32) (*Texture, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("texture file `%s` not found on disk: %v", filename, err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	texture := &Texture{
		width:  uint32(rgba.Rect.Size().X),
		height: uint32(rgba.Rect.Size().Y),
	}
	gl.GenTextures(1, &texture.id)
	gl.BindTexture(gl.TEXTURE_2D, texture.id)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, wrap)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, wrap)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(texture.width),
		int32(texture.height),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return texture, nil
}

// Width returns the width of the texture.
func (t *Texture) Width() uint32 {
	return t.width
}

// Height returns the height of the texture.
func (t *Texture) Height() uint32 {
	return t.height
}

// Bind activates the provided texture unit and binds the texture.
func (t *Texture) Bind(location uint32) {
	gl.ActiveTexture(location)
	gl.BindTexture(gl.TEXTURE_2D, t.id)
}

// Unbind will unbind the texture.
func (t *Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// Destroy deallocates the texture buffer.
func (t *Texture) Destroy() {
	gl.DeleteTextures(1, &t.id)
	t.id = 0
}
