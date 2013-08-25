// Copyright 2013 Arne Roomann-Kurrik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package system

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"image"
	"image/draw"
	"image/png"
	"os"
)

type Texture struct {
	texture gl.Texture
	Width   int
	Height  int
	Frames  [][]int
}

func LoadTexture(path string, smoothing int, framewidth int, frameheight int) (texture *Texture, err error) {
	var (
		img    image.Image
		bounds image.Rectangle
		obounds image.Rectangle
		gltex  gl.Texture
	)
	if img, err = loadPNG(path); err != nil {
		return
	}
	obounds = img.Bounds()
	img = getPow2Image(img)
	bounds = img.Bounds()
	if gltex, err = getGLTexture(img, smoothing); err != nil {
		return
	}
	texture = &Texture{
		texture: gltex,
		Width:   bounds.Dx(),
		Height:  bounds.Dy(),
		Frames:  make([][]int, 0),
	}
	frames := obounds.Dx() / framewidth * obounds.Dy() / frameheight
	for i := 0; i < frames; i++ {
		var (
			minx = (i * framewidth) % obounds.Dx()
			maxx = minx + framewidth
			miny = ((i * framewidth) / obounds.Dx()) * frameheight
			maxy = miny + frameheight
		)
		texture.Frames = append(texture.Frames, []int{
			minx,
			maxx,
			miny,
			maxy,
		})
	}
	return
}

func (t *Texture) MinX(i int) float64 {
	return float64(t.Frames[i][0]) / float64(t.Width)
}

func (t *Texture) MaxX(i int) float64 {
	return float64(t.Frames[i][1]) / float64(t.Width)
}

func (t *Texture) MinY(i int) float64 {
	return 1 - float64(t.Frames[i][2])/float64(t.Height)
}

func (t *Texture) MaxY(i int) float64 {
	return 1 - float64(t.Frames[i][3])/float64(t.Height)
}

func (t *Texture) Bind() {
	t.texture.Bind(gl.TEXTURE_2D)
}

func (t *Texture) Unbind() {
	t.texture.Unbind(gl.TEXTURE_2D)
}

func (t *Texture) Dispose() {
	t.texture.Delete()
}

func getPow2(i int) int {
	p2 := 1
	for p2 < i {
		p2 = p2 << 1
	}
	return p2
}

func getPow2Image(img image.Image) image.Image {
	var (
		b   = img.Bounds()
		p2w = getPow2(b.Max.X)
		p2h = getPow2(b.Max.Y)
	)
	if p2w == b.Max.X && p2h == b.Max.Y {
		return img
	}
	out := image.NewRGBA(image.Rect(0, 0, p2w, p2h))
	draw.Draw(out, b, img, image.ZP, draw.Src)
	return out
}

func loadPNG(path string) (img image.Image, err error) {
	var file *os.File
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()
	img, err = png.Decode(file)
	return
}

func getGLTexture(img image.Image, smoothing int) (gltexture gl.Texture, err error) {
	var data *bytes.Buffer
	if data, err = encodeTGA("texture", img); err != nil {
		return
	}
	gltexture = gl.GenTexture()
	gltexture.Bind(gl.TEXTURE_2D)
	if !glfw.LoadMemoryTexture2D(data.Bytes(), 0) {
		err = fmt.Errorf("Failed to load texture")
		return
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, smoothing)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, smoothing)
	return
}

func encodeTGA(name string, img image.Image) (buf *bytes.Buffer, err error) {
	var (
		bounds image.Rectangle = img.Bounds()
		ident  []byte          = []byte(name)
		width  []byte          = make([]byte, 2)
		height []byte          = make([]byte, 2)
		nrgba  *image.NRGBA
		data   []byte
	)
	binary.LittleEndian.PutUint16(width, uint16(bounds.Dx()))
	binary.LittleEndian.PutUint16(height, uint16(bounds.Dy()))

	// See http://paulbourke.net/dataformats/tga/
	buf = &bytes.Buffer{}
	buf.WriteByte(byte(len(ident)))
	buf.WriteByte(0)
	buf.WriteByte(2) // uRGBI
	buf.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0})
	buf.Write([]byte(width))
	buf.Write([]byte(height))
	buf.WriteByte(32) // Bits per pixel
	buf.WriteByte(8)
	if buf.Len() != 18 {
		err = fmt.Errorf("TGA header is not 18 bytes: %v", buf.Len())
		return
	}

	nrgba = image.NewNRGBA(bounds)
	draw.Draw(nrgba, bounds, img, bounds.Min, draw.Src)
	buf.Write(ident)
	data = make([]byte, bounds.Dx()*bounds.Dy()*4)
	var (
		lineLength int = bounds.Dx() * 4
		destOffset int = len(data) - lineLength
	)
	for srcOffset := 0; srcOffset < len(nrgba.Pix); {
		var (
			dest   = data[destOffset : destOffset+lineLength]
			source = nrgba.Pix[srcOffset : srcOffset+nrgba.Stride]
		)
		copy(dest, source)
		destOffset -= lineLength
		srcOffset += nrgba.Stride
	}
	for x := 0; x < len(data); {
		buf.WriteByte(data[x+2])
		buf.WriteByte(data[x+1])
		buf.WriteByte(data[x+0])
		buf.WriteByte(data[x+3])
		x += 4
	}
	return
}
