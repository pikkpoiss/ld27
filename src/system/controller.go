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
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"log"
)

// Handles window close events.
type CloseHandler func() int

// Handles key press events.
type KeyHandler func(key int, state int)

// Main system controller.
type Controller struct {
	Win *Window
}

// Creates an initialized controller.
func NewController() (c *Controller, err error) {
	if err = glfw.Init(); err != nil {
		return
	}
	c = &Controller{}
	return
}

// Opens a new window.
func (c *Controller) Open(win *Window) (err error) {
	c.Win = win
	mode := glfw.Windowed
	if win.Fullscreen {
		mode = glfw.Fullscreen
	}
	if win.Resize == false {
		glfw.OpenWindowHint(glfw.WindowNoResize, 1)
	}
	if err = glfw.OpenWindow(win.Width, win.Height, 0, 0, 0, 0, 0, 0, mode); err != nil {
		return
	}
	gl.Init()
	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	v1, v2, v3 := glfw.GLVersion()
	log.Printf("OpenGL version: %v %v %v\n", v1, v2, v3)
	fb_supported := glfw.ExtensionSupported("GL_EXT_framebuffer_object")
	log.Printf("Framebuffer supported: %v\n", fb_supported)
	c.SetClearColor(0, 0, 0, 0)
	if win.VSync == true {
		glfw.SetSwapInterval(1) // Limit to refresh
	}
	glfw.SetWindowTitle(win.Title)
	glfw.SetWindowSizeCallback(func(w, h int) {
		log.Printf("Resizing window to %v, %v\n", w, h)
		c.resize()
	})
	err = c.resize()
	return
}

// Call to clean up after you're done.
func (c *Controller) Terminate() {
	glfw.Terminate()
}

// Handles window resize.
func (c *Controller) resize() (err error) {
	c.Win.Width, c.Win.Height = glfw.WindowSize()
	return
}

// Clamps a value to a max.
func (c *Controller) clamp(i int, max int) gl.GLclampf {
	return gl.GLclampf(float64(i) / float64(max))
}

// Sets the background clear color.
func (c *Controller) SetClearColor(r int, g int, b int, a int) {
	gl.ClearColor(c.clamp(r, 255), c.clamp(g, 255), c.clamp(b, 255), c.clamp(a, 255))
	gl.ClearDepth(1.0)
}

// Specify a function to call when the window is closed.
func (c *Controller) SetCloseCallback(handler CloseHandler) {
	glfw.SetWindowCloseCallback(glfw.WindowCloseHandler(handler))
}

// Specify a function to call when a key is pressed.
func (c *Controller) SetKeyCallback(handler KeyHandler) {
	glfw.SetKeyCallback(glfw.KeyHandler(handler))
}

// Check whether a key is pressed.
func (c *Controller) Key(key int) int {
	return glfw.Key(key)
}
