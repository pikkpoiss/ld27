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
	"github.com/go-gl/gl"
)

type Camera struct {
	w float64
	h float64
	x float64
	y float64
}

func NewCamera(x float64, y float64, w float64, h float64) (c *Camera) {
	c = &Camera{
		w: w,
		h: h,
		x: x,
		y: y,
	}
	return
}

func (c *Camera) SetProjection() {
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(c.x, c.x+c.w, c.y+c.h, c.y, 1, -1)
}
