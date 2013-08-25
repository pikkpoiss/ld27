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

package main

import (
	"./system"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
)

func BeginPaint() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func EndPaint() {
	gl.Flush()
	glfw.SwapBuffers()
}

func PaintCast(ctrl *system.Controller, c *Cast) {
	c.Texture.Bind()
	for _, a := range c.Actors {
		var (
			minx  = int(a.X()) - c.OffsetX
			miny  = int(a.Y()) - c.OffsetY
			maxx  = minx + c.Width
			maxy  = miny + c.Height
			frame = a.GetFrame() + c.TextureCols*a.TextureRow()
		)
		if a.FlipX() {
			maxx ^= minx
			minx ^= maxx
			maxx ^= minx
		}
		paintSprite(minx, miny, maxx, maxy, c.Texture, frame)
	}
	c.Texture.Unbind()
}

func PaintMap(ctrl *system.Controller, tm *system.TiledMap) {
	var x int
	var y int
	var tw int
	var th int
	for _, ts := range tm.Tilesets {
		tw = ts.Tilewidth
		th = ts.Tileheight
		ts.Texture.Bind()
		for _, l := range tm.Layers {
			if l.Type != "tilelayer" {
				continue
			}
			for i, gid := range l.Data {
				//log.Printf("PaintMap %v %v %v\n", gid, ts.Firstgid, ts.Lastgid)
				if gid >= ts.Firstgid && gid < ts.Lastgid {
					x = i % l.Width
					y = i / l.Width
					paintTile(x, y, tw, th, ts.Texture, gid-ts.Firstgid)
				}
			}
		}
		ts.Texture.Unbind()
	}
}

func paintSprite(minx int, miny int, maxx int, maxy int, t *system.Texture, index int) {
	gl.MatrixMode(gl.TEXTURE)
	gl.Begin(gl.QUADS)
	gl.TexCoord2d(t.MinX(index), t.MinY(index))
	gl.Vertex2i(minx, miny)
	gl.TexCoord2d(t.MaxX(index), t.MinY(index))
	gl.Vertex2i(maxx, miny)
	gl.TexCoord2d(t.MaxX(index), t.MaxY(index))
	gl.Vertex2i(maxx, maxy)
	gl.TexCoord2d(t.MinX(index), t.MaxY(index))
	gl.Vertex2i(minx, maxy)
	gl.End()
	gl.MatrixMode(gl.MODELVIEW)
}

func paintTile(x int, y int, w int, h int, t *system.Texture, index int) {
	var (
		minx = x * w
		miny = y * h
		maxx = (x + 1) * w
		maxy = (y + 1) * h
	)
	paintSprite(minx, miny, maxx, maxy, t, index)
}
